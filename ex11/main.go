package main

import (
	"fmt"

	"github.com/jackytck/gophercises/ex11/blackjack"
)

func main() {
	game := blackjack.New()
	winnings := game.Play(blackjack.HumanAI())
	fmt.Println(winnings)
}
