package engine

import (
	"fmt"
	"math/rand"
	"time"
)

type Suit int

const (
	Hearts Suit = iota
	Diamonds
	Clubs
	Spades
)

func (s Suit) String() string {
	return [...]string{"Hearts", "Diamonds", "Clubs", "Spades"}[s]
}

type Value int

const (
	Joker Value = iota
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
	Ace
)

// overwrite string method for Value
func (v Value) String() string {
	return [...]string{"Joker", "Two", "Three", "Four",
		"Five", "Six", "Seven", "Eight", "Nine", "Ten",
		"Jack", "Queen", "King", "Ace"}[v]
}

type Card struct {
	Suit  Suit  `json:"suit"`
	Value Value `json:"value"`
}

func (s Card) String() string {
	return fmt.Sprintf("%s of %s", s.Value, s.Suit)
}

type Deck struct {
	Cards  []Card
	Dealed int
}

func NewDeck() *Deck {
	suits := []Suit{Spades, Hearts, Diamonds, Clubs}
	values := []Value{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}

	deck := &Deck{}

	for _, suit := range suits {
		for _, value := range values {
			card := Card{Suit: suit, Value: value}
			deck.Cards = append(deck.Cards, card)
		}
	}

	return deck
}

func (d *Deck) Shuffle() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := len(d.Cards) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}
	d.Dealed = 0
}

// Cut the deck at a random position
func (d *Deck) CutTheCard() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cutIndex := rng.Intn(len(d.Cards))
	d.Cards = append(d.Cards[cutIndex:], d.Cards[:cutIndex]...)
}

func (d *Deck) Draw() Card {
	card := d.Cards[d.Dealed]
	d.Dealed++
	return card
}

// CommunityCards represents the community cards in a Poker game.
type CommunityCards struct {
	Cards []Card
}

func (c CommunityCards) String() string {
	var cardsString string
	for _, card := range c.Cards {
		cardsString += "[" + card.String() + "]"
	}
	return cardsString
}

func (c *CommunityCards) AddCard(card Card) {
	c.Cards = append(c.Cards, card)
}

func (c *CommunityCards) Reset() {
	c.Cards = []Card{}
}
