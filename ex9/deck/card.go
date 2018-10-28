//go:generate stringer -type=Suit,Rank
package deck

import "fmt"

// Suit represents a suit.
type Suit uint8

const (
	Spade Suit = iota
	Diamond
	Club
	Heart
	Joker
)

type Rank uint8

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

// Card represents a playing card.
type Card struct {
	Suit
	Rank
}

// New creates a deck of cards in sorted order.
func New() []Card {
	var ret []Card
	for s := Spade; s < Joker; s++ {
		for r := Ace; r <= King; r++ {
			c := Card{
				Suit: s,
				Rank: r,
			}
			ret = append(ret, c)
		}
	}
	return ret
}

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %ss", c.Rank, c.Suit)
}
