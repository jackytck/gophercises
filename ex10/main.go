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

func main() {
	cards := deck.New(deck.Deck(3), deck.Shuffle)
	var card deck.Card
	var dealer, player Hand
	for i := 0; i < 2; i++ {
		for _, hand := range []*Hand{&dealer, &player} {
			card, cards = draw(cards)
			*hand = append(*hand, card)
		}
	}
	var input string
	for input != "s" {
		fmt.Println("Dealer", dealer.DealerString())
		fmt.Println("Player", player)
		fmt.Println("What will you do? (h)it, (s)tand")
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			card, cards = draw(cards)
			player = append(player, card)
		}
	}

	// if dealer score <= 16, then hit
	// if dealer has a soft 17, then hit
	for dealer.Score() <= 16 || (dealer.Score() == 17 && dealer.MinScore() != 17) {
		card, cards = draw(cards)
		dealer = append(dealer, card)
	}

	dScore, pScore := dealer.Score(), player.Score()
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Dealer:", dealer, "\nScore", dScore)
	fmt.Println("Player:", player, "\nScore", pScore)
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