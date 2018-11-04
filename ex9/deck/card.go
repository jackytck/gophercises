//go:generate stringer -type=Suit,Rank
package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// Suit represents a suit.
type Suit uint8

const (
	Spade Suit = iota
	Diamond
	Club
	Heart
	Joker
)

var suits = [...]Suit{Spade, Diamond, Club, Heart}

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

const (
	minRank = Ace
	maxRank = King
)

// Card represents a playing card.
type Card struct {
	Suit
	Rank
}

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %ss", c.Rank, c.Suit)
}

// New creates a deck of cards in sorted order.
func New(opts ...func([]Card) []Card) []Card {
	var cards []Card
	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			c := Card{
				Suit: suit,
				Rank: rank,
			}
			cards = append(cards, c)
		}
	}
	for _, f := range opts {
		cards = f(cards)
	}
	return cards
}

// DefaultSort is the functional option for sorting the deck in default order.
func DefaultSort(cards []Card) []Card {
	sort.Slice(cards, Less(cards))
	return cards
}

// Sort returns the functional option for sorting by giving a less function.
func Sort(less func(cards []Card) func(i, j int) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		sort.Slice(cards, less(cards))
		return cards
	}
}

// Less gives the less function for doing default sort.
func Less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absRank(cards[i]) < absRank(cards[j])
	}
}

func absRank(c Card) int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}

var shuffleRand = rand.New(rand.NewSource(time.Now().Unix()))

// Shuffle shuffles the cards.
func Shuffle(cards []Card) []Card {
	ret := make([]Card, len(cards))
	perm := shuffleRand.Perm(len(cards))
	for i, v := range perm {
		ret[i] = cards[v]
	}
	return ret
}

// Jokers put n cards of joker into the deck.
func Jokers(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		for i := 0; i < n; i++ {
			cards = append(cards, Card{
				Suit: Joker,
				Rank: Rank(i + 1),
			})
		}
		return cards
	}
}

// Filter filters the deck of cards by the given filter function.
func Filter(f func(Card) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card
		for _, c := range cards {
			if !f(c) {
				ret = append(ret, c)
			}
		}
		return ret
	}
}

// Deck multiplies the given deck by n times.
func Deck(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card
		for i := 0; i < n; i++ {
			ret = append(ret, cards...)
		}
		return ret
	}
}
