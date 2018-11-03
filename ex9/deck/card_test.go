package deck

import (
	"fmt"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Heart})
	fmt.Println(Card{Rank: Two, Suit: Spade})
	fmt.Println(Card{Rank: Nine, Suit: Diamond})
	fmt.Println(Card{Rank: Jack, Suit: Club})
	fmt.Println(Card{Suit: Joker})

	// Output:
	// Ace of Hearts
	// Two of Spades
	// Nine of Diamonds
	// Jack of Clubs
	// Joker
}

func TestNew(t *testing.T) {
	cards := New()
	if len(cards) != 13*4 {
		t.Error("Wrong number of cards in a new deck.")
	}
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)
	first := Card{Suit: Spade, Rank: Ace}
	if cards[0] != first {
		t.Error("Expected Ace of Spades as first card. Received:", cards[0])
	}
}

func TestSort(t *testing.T) {
	cards := New(Sort(Less))
	first := Card{Suit: Spade, Rank: Ace}
	if cards[0] != first {
		t.Error("Expected Ace of Spades as first card. Received:", cards[0])
	}
}

func TestJokers(t *testing.T) {
	jokerCnt := 3
	cards := New(Jokers(jokerCnt))
	count := 0
	for _, c := range cards {
		if c.Suit == Joker {
			count++
		}
	}
	if count != jokerCnt {
		t.Errorf("Expected %d Jokers, received: %d", jokerCnt, count)
	}
}

func TestFilter(t *testing.T) {
	filter := func(card Card) bool {
		return card.Rank == Two || card.Rank == Three
	}
	cards := New(Filter(filter))
	for _, c := range cards {
		if c.Rank == Two || c.Rank == Three {
			t.Error("Expected all twos and threes to be filtered outi.")
		}
	}
}

func TestDeck(t *testing.T) {
	n := 3
	cards := New(Deck(n))
	if len(cards) != 13*4*n {
		t.Errorf("Expected %d cards, received %d cards.", 13*4*n, len(cards))
	}
}