package blackjack

import (
	"errors"
	"fmt"

	"github.com/jackytck/gophercises/ex9/deck"
)

type state int8

// Options is the game options.
type Options struct {
	Decks           int
	Hands           int
	BlackjackPayout float64
	Verbose         bool
}

const (
	stateBet state = iota
	statePlayerTurn
	stateDealerTurn
	stateHandOver
)

// New creates the game.
func New(opts Options) Game {
	g := Game{
		state:    statePlayerTurn,
		dealerAI: dealerAI{},
		balance:  0,
	}
	if opts.Decks == 0 {
		opts.Decks = 3
	}
	if opts.Hands == 0 {
		opts.Hands = 100
	}
	if opts.BlackjackPayout == 0 {
		opts.BlackjackPayout = 1.5
	}
	g.nDecks = opts.Decks
	g.nHands = opts.Hands
	g.blackjackPayout = opts.BlackjackPayout
	g.verbose = opts.Verbose
	return g
}

// Game represents a game.
type Game struct {
	// unexported fields
	nDecks          int
	nHands          int
	blackjackPayout float64

	state state
	deck  []deck.Card

	dealer   []deck.Card
	dealerAI AI

	player    []hand
	handIdx   int
	playerBet int
	balance   int

	verbose bool
}

// log logs a string if verbose is set.
func (g *Game) log(s string) {
	if g.verbose {
		fmt.Println(s)
	}
}

// currentHand return the address of the current player.
func (g *Game) currentHand() *[]deck.Card {
	switch g.state {
	case stateDealerTurn:
		return &g.dealer
	case statePlayerTurn:
		return &g.player[g.handIdx].cards
	default:
		panic("it isn't currently any player's turn")
	}
}

type hand struct {
	cards []deck.Card
	bet   int
}

func (g *Game) bet(ai AI, shuffled bool) {
	bet := ai.Bet(shuffled)
	if bet < 100 {
		panic("bet must be at least 100")
	}
	g.playerBet = bet
}

func (g *Game) deal() {
	playerHand := make([]deck.Card, 0, 5)
	g.handIdx = 0
	g.dealer = make([]deck.Card, 0, 5)
	// g.player = make([]hand, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, g.deck = draw(g.deck)
		g.dealer = append(g.dealer, card)
		card, g.deck = draw(g.deck)
		playerHand = append(playerHand, card)
	}
	// test splitting
	// playerHand = []deck.Card{
	// 	{Rank: deck.Seven},
	// 	{Rank: deck.Seven},
	// }
	g.player = []hand{
		{
			cards: playerHand,
			bet:   g.playerBet,
		},
	}
	// test blackjack
	// g.player = []deck.Card{
	// 	{Rank: deck.Six},
	// 	{Rank: deck.Five},
	// 	{Rank: deck.Ten},
	// }
	// g.dealer = []deck.Card{
	// 	{Rank: deck.Ace},
	// 	{Rank: deck.Ten},
	// }
	g.state = statePlayerTurn
}

// Play plays the game with ai interface.
func (g *Game) Play(ai AI) int {
	g.deck = nil
	min := 52 * g.nDecks / 3

	for i := 0; i < g.nHands; i++ {
		shuffled := false
		if len(g.deck) < min {
			g.deck = deck.New(deck.Deck(g.nDecks), deck.Shuffle)
			shuffled = true
		}

		g.bet(ai, shuffled)
		g.deal()
		if Blackjack(g.dealer...) {
			endRound(g, ai)
			continue
		}

		for g.state == statePlayerTurn {
			hand := make([]deck.Card, len(*g.currentHand()))
			copy(hand, *g.currentHand())
			move := ai.Play(hand, g.dealer[0])
			err := move(g)
			switch err {
			case errBust:
				MoveStand(g)
			case nil:
				// noop
			default:
				panic(err)
			}
		}

		for g.state == stateDealerTurn {
			hand := make([]deck.Card, len(g.dealer))
			copy(hand, g.dealer)
			move := g.dealerAI.Play(hand, g.dealer[0])
			move(g)
		}

		endRound(g, ai)
	}

	return g.balance
}

func endRound(g *Game, ai AI) {
	dScore := Score(g.dealer...)
	dBlackjack := Blackjack(g.dealer...)
	allHands := make([][]deck.Card, len(g.player))
	for hi, hand := range g.player {
		allHands[hi] = hand.cards
		pScore, pBlackjack := Score(allHands[hi]...), Blackjack(allHands[hi]...)
		winnings := hand.bet
		switch {
		case dBlackjack && pBlackjack:
			g.log("Both got blackjack")
			winnings = 0
		case dBlackjack:
			g.log("Dealer blackjack")
			winnings *= -1
		case pBlackjack:
			g.log("Player blackjack")
			winnings = int(float64(winnings) * g.blackjackPayout)
		case pScore > 21:
			g.log("You busted")
			g.balance *= -1
		case dScore > 21:
			g.log("Dealer busted")
		case pScore > dScore:
			g.log("You win!")
		case dScore > pScore:
			g.log("You lose!")
			g.balance *= -1
		case dScore == pScore:
			// g.log("Draw")
			winnings = 0
		}
		g.balance += winnings
	}
	g.log("")
	ai.Results(allHands, g.dealer)
	g.dealer = nil
	g.player = nil
	g.handIdx = 0
}

var (
	errBust = errors.New("hand score exceeded 21")
)

// Move is a move function.
type Move func(*Game) error

// MoveHit is a hit function.
func MoveHit(g *Game) error {
	hand := g.currentHand()
	var card deck.Card
	card, g.deck = draw(g.deck)
	*hand = append(*hand, card)
	if Score(*hand...) > 21 {
		return errBust
	}
	return nil
}

// MoveSplit splits two cards with the same rank.
func MoveSplit(g *Game) error {
	h := *g.currentHand()
	if len(h) != 2 {
		return errors.New("you can only split with two cards in your hand")
	}
	if h[0].Rank != h[1].Rank {
		return errors.New("both cards must have the same rank to split")
	}
	b := g.player[g.handIdx].bet
	g.player = []hand{
		{cards: []deck.Card{h[0]}, bet: b},
		{cards: []deck.Card{h[1]}, bet: b},
	}
	return nil
}

// MoveDouble doubles down aon a hand with 2 cards.
func MoveDouble(g *Game) error {
	if len(*g.currentHand()) != 2 {
		return errors.New("can only double on a hand with 2 cards")
	}
	g.playerBet *= 2
	if err := MoveHit(g); err != nil {
		return err
	}
	return MoveStand(g)
}

// MoveStand is a stand function.
func MoveStand(g *Game) error {
	if g.state == stateDealerTurn {
		g.state++
		return nil
	}
	if g.state == statePlayerTurn {
		g.handIdx++
		if g.handIdx >= len(g.player) {
			g.state++
		}
		return nil
	}
	return errors.New("invalid state")
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

// minScore computes the minimum score of a hand.
func minScore(hand ...deck.Card) int {
	score := 0
	for _, c := range hand {
		score += min(int(c.Rank), 10)
	}
	return score
}

// Score computes the best blackjack score of a hand.
func Score(hand ...deck.Card) int {
	minScore := minScore(hand...)
	// no card could be reinterpreted as 11, otherwise would bust
	if minScore > 11 {
		return minScore
	}
	for _, c := range hand {
		if c.Rank == deck.Ace {
			// replace an Ace from 1 to 11
			return minScore + 10
		}
	}
	return minScore
}

// Soft returns true if it is a soft socre - that is if an ace is being counted as 11.
func Soft(hand ...deck.Card) bool {
	ms := minScore(hand...)
	score := Score(hand...)
	return ms != score
}

// Blackjack returns true if a hand is a blackjack.
func Blackjack(hand ...deck.Card) bool {
	return len(hand) == 2 && Score(hand...) == 21
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
