package blackjack

import (
	"fmt"

	"github.com/jackytck/gophercises/ex9/deck"
)

type state int8

// Options is the game options.
type Options struct {
	Decks           int
	Hands           int
	BlackjackPayout float64
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
	return g
}

// Game represents a game.
type Game struct {
	// unexported fields
	deck            []deck.Card
	nDecks          int
	nHands          int
	state           state
	dealer          []deck.Card
	player          []deck.Card
	dealerAI        AI
	balance         int
	blackjackPayout float64
}

// CurrentHand return the address of the current player.
func (g *Game) currentHand() *[]deck.Card {
	switch g.state {
	case stateDealerTurn:
		return &g.dealer
	case statePlayerTurn:
		return &g.player
	default:
		panic("it isn't currently any player's turn")
	}
}

func (g *Game) deal() {
	g.dealer = make([]deck.Card, 0, 5)
	g.player = make([]deck.Card, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, g.deck = draw(g.deck)
		g.dealer = append(g.dealer, card)
		card, g.deck = draw(g.deck)
		g.player = append(g.player, card)
	}
	g.state = statePlayerTurn
}

// Play plays the game with ai interface.
func (g *Game) Play(ai AI) int {
	g.deck = nil
	min := 52 * g.nDecks / 3

	for i := 0; i < g.nHands; i++ {
		if len(g.deck) < min {
			g.deck = deck.New(deck.Deck(g.nDecks), deck.Shuffle)
		}

		g.deal()

		for g.state == statePlayerTurn {
			hand := make([]deck.Card, len(g.player))
			copy(hand, g.player)
			move := ai.Play(hand, g.dealer[0])
			move(g)
		}

		for g.state == stateDealerTurn {
			hand := make([]deck.Card, len(g.dealer))
			copy(hand, g.dealer)
			move := g.dealerAI.Play(hand, g.dealer[0])
			move(g)
		}

		endHand(g, ai)
	}

	return g.balance
}

func endHand(g *Game, ai AI) {
	dScore, pScore := Score(g.dealer...), Score(g.player...)
	switch {
	case pScore > 21:
		fmt.Println("You busted")
		g.balance--
	case dScore > 21:
		fmt.Println("Dealer busted")
		g.balance++
	case pScore > dScore:
		fmt.Println("You win!")
		g.balance++
	case dScore > pScore:
		fmt.Println("You lose!")
		g.balance--
	case dScore == pScore:
		fmt.Println("Draw")
	}
	fmt.Println()
	ai.Results([][]deck.Card{g.player}, g.dealer)
	g.dealer = nil
	g.player = nil
}

// Move is a move function.
type Move func(*Game)

// MoveHit is a hit function.
func MoveHit(g *Game) {
	hand := g.currentHand()
	var card deck.Card
	card, g.deck = draw(g.deck)
	*hand = append(*hand, card)
	if Score(*hand...) > 21 {
		MoveStand(g)
	}
}

// MoveStand is a stand function.
func MoveStand(g *Game) {
	g.state++
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
