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
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Dealer:", dealer)
	fmt.Println("Player:", player)
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}
