package main

import (
	"fmt"

	"github.com/jackytck/gophercises/ex11/blackjack"
)

func main() {
	opts := blackjack.Options{
		Decks:           3,
		Hands:           5,
		BlackjackPayout: 1.5,
	}
	game := blackjack.New(opts)
	winnings := game.Play(blackjack.HumanAI())
	fmt.Println(winnings)
}
