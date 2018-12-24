package main

import (
	"fmt"
	"strings"

	"github.com/jackytck/gophercises/ex9/deck"
)

// Hand represents a hand in the blackjack game.
type Hand []deck.Card

func (h Hand) String() string {
	strs := make([]string, len(h))
	for i, v := range h {
		strs[i] = v.String()
	}
	return strings.Join(strs, ",")
}

// DealerString displays the first card of the dealer hand only.
func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

// MinScore computes the minimum score of a hand.
func (h Hand) MinScore() int {
	score := 0
	for _, c := range h {
		score += min(int(c.Rank), 10)
	}
	return score
}

// Score computes the score of a hand.
func (h Hand) Score() int {
	minScore := h.MinScore()
	// no card could be reinterpreted as 11, otherwise would bust
	if minScore > 11 {
		return minScore
	}
	for _, c := range h {
		if c.Rank == deck.Ace {
			// replace an Ace from 1 to 11
			return minScore + 10
		}
	}
	return minScore
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Shuffle(gs GameState) GameState {
	ret := clone(gs)
	ret.Deck = deck.New(deck.Deck(3), deck.Shuffle)
	return ret
}

func Deal(gs GameState) GameState {
	ret := clone(gs)
	ret.Dealer = make(Hand, 0, 5)
	ret.Player = make(Hand, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, ret.Deck = draw(ret.Deck)
		ret.Dealer = append(ret.Dealer, card)
		card, ret.Deck = draw(ret.Deck)
		ret.Player = append(ret.Player, card)
	}
	ret.State = StatePlayerTurn
	return ret
}

func Stand(gs GameState) GameState {
	ret := clone(gs)
	ret.State++
	return ret
}

func Hit(gs GameState) GameState {
	ret := clone(gs)
	hand := ret.CurrentPlayer()
	var card deck.Card
	card, ret.Deck = draw(ret.Deck)
	*hand = append(*hand, card)
	if hand.Score() > 21 {
		return Stand(ret)
	}
	return ret
}

func EndHand(gs GameState) GameState {
	ret := clone(gs)
	dScore, pScore := ret.Dealer.Score(), ret.Player.Score()
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Dealer:", ret.Dealer, "\nScore", dScore)
	fmt.Println("Player:", ret.Player, "\nScore", pScore)
	switch {
	case pScore > 21:
		fmt.Println("You busted")
	case dScore > 21:
		fmt.Println("Dealer busted")
	case pScore > dScore:
		fmt.Println("You win!")
	case dScore > pScore:
		fmt.Println("You lose!")
	case dScore == pScore:
		fmt.Println("Draw")
	}
	fmt.Println()
	ret.Dealer = nil
	ret.Player = nil
	return ret
}

func main() {
	var gs GameState
	gs = Shuffle(gs)

	for i := 0; i < 10; i++ {
		gs = Deal(gs)

		var input string
		for gs.State == StatePlayerTurn {
			fmt.Println("Dealer", gs.Dealer.DealerString())
			fmt.Println("Player", gs.Player)
			fmt.Println("What will you do? (h)it, (s)tand")
			fmt.Scanf("%s\n", &input)
			switch input {
			case "h":
				gs = Hit(gs)
			case "s":
				gs = Stand(gs)
			default:
				fmt.Println("Invalid option:", input)
			}
		}

		for gs.State == StateDealerTurn {
			// if dealer score <= 16, then hit
			// if dealer has a soft 17, then hit
			if gs.Dealer.Score() <= 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
				gs = Hit(gs)
			} else {
				gs = Stand(gs)
			}
		}

		gs = EndHand(gs)
	}
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

type State int8

const (
	StatePlayerTurn = iota
	StateDealerTurn
	StateHandOver
)

// GameState represents the game state of the Blackjack game.
type GameState struct {
	Deck   []deck.Card
	State  State
	Dealer Hand
	Player Hand
}

// CurrentPlayer return the address of the current player.
func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.State {
	case StateDealerTurn:
		return &gs.Dealer
	case StatePlayerTurn:
		return &gs.Player
	default:
		panic("it isn't currently any player's turn")
	}
}

func clone(gs GameState) GameState {
	ret := GameState{
		Deck:   make([]deck.Card, len(gs.Deck)),
		State:  gs.State,
		Dealer: make(Hand, len(gs.Dealer)),
		Player: make(Hand, len(gs.Player)),
	}
	copy(ret.Deck, gs.Deck)
	copy(ret.Dealer, gs.Dealer)
	copy(ret.Player, gs.Player)
	return ret
}
