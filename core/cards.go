package engine

import (
	msgpb "go-pk-server/gen"
	"math/rand"
	"time"
)

type Deck struct {
	Cards  []*msgpb.Card
	Dealed int
}

func NewDeck() *Deck {
	suits := []msgpb.SuitType{
		msgpb.SuitType_SPADES,
		msgpb.SuitType_HEARTS,
		msgpb.SuitType_DIAMONDS,
		msgpb.SuitType_CLUBS,
	}
	values := []msgpb.RankType{
		msgpb.RankType_DEUCE,
		msgpb.RankType_THREE,
		msgpb.RankType_FOUR,
		msgpb.RankType_FIVE,
		msgpb.RankType_SIX,
		msgpb.RankType_SEVEN,
		msgpb.RankType_EIGHT,
		msgpb.RankType_NINE,
		msgpb.RankType_TEN,
		msgpb.RankType_JACK,
		msgpb.RankType_QUEEN,
		msgpb.RankType_KING,
		msgpb.RankType_ACE,
	}

	deck := &Deck{}

	for _, suit := range suits {
		for _, value := range values {
			deck.Cards = append(deck.Cards, &msgpb.Card{Suit: suit, Rank: value})
		}
	}

	return deck
}

func (d *Deck) Shuffle() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	// Fisher-Yates shuffle
	for i := len(d.Cards) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	}

	oneFifthLen := len(d.Cards) / 5
	// Do Simple cut-and-riffle shuffle 2 times
	for i := 0; i < 2; i++ {
		rng = rand.New(rand.NewSource(time.Now().UnixNano()))
		// split the deck in nearly half, merge them randomly
		half := rng.Intn(oneFifthLen) + oneFifthLen/5*2
		left := d.Cards[:half]
		right := d.Cards[half:]

		// Merge the two halves by picking random elements from each half
		mergedDeck := make([]*msgpb.Card, 0, len(d.Cards))
		for len(left) > 0 && len(right) > 0 {
			if rand.Intn(2) == 0 {
				mergedDeck = append(mergedDeck, left[0])
				left = left[1:]
			} else {
				mergedDeck = append(mergedDeck, right[0])
				right = right[1:]
			}
		}
		// Append remaining cards if any half has leftover cards
		mergedDeck = append(mergedDeck, left...)
		mergedDeck = append(mergedDeck, right...)

		// Copy merged deck back into the original deck
		copy(d.Cards, mergedDeck)
	}

	// Handover shuffer,cCut the deck into three at two random positions
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	cutPoint1 := rng.Intn(oneFifthLen/5+1) + oneFifthLen
	cutPoint2 := rng.Intn(oneFifthLen/5+1) + oneFifthLen*3
	first := d.Cards[:cutPoint1]
	second := d.Cards[cutPoint1:cutPoint2]
	third := d.Cards[cutPoint2:]
	d.Cards = append(third, second...)
	d.Cards = append(d.Cards, first...)

	// Do last cut-and-riffle shuffle
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	// split the deck in nearly half, merge them randomly
	half := rng.Intn(oneFifthLen) + oneFifthLen/5*2
	left := d.Cards[:half]
	right := d.Cards[half:]

	// Merge the two halves by picking random elements from each half
	mergedDeck := make([]*msgpb.Card, 0, len(d.Cards))
	for len(left) > 0 && len(right) > 0 {
		if rand.Intn(2) == 0 {
			mergedDeck = append(mergedDeck, left[0])
			left = left[1:]
		} else {
			mergedDeck = append(mergedDeck, right[0])
			right = right[1:]
		}
	}
	// Append remaining cards if any half has leftover cards
	mergedDeck = append(mergedDeck, left...)
	mergedDeck = append(mergedDeck, right...)

	// Copy merged deck back into the original deck
	copy(d.Cards, mergedDeck)

	// Cut the deck at a random position in the middle
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	// random cut index from 1/3 to 2/3 of the deck
	cutIndex := rng.Intn(len(d.Cards)/3) + len(d.Cards)/3
	d.Cards = append(d.Cards[cutIndex:], d.Cards[:cutIndex]...)

	if len(d.Cards) != 52 {
		panic("Deck is not 52 cards")
	}
	// Start dealing from the beginning
	d.Dealed = 0
}

func (d *Deck) Draw() *msgpb.Card {
	card := d.Cards[d.Dealed]
	d.Dealed++
	return card
}

// CommunityCards represents the community cards in a Poker game.
type CommunityCards struct {
	Cards []*msgpb.Card
}

func (c CommunityCards) String() string {
	var cardsString string
	for _, card := range c.Cards {
		cardsString += "[" + card.String() + "]"
	}
	return cardsString
}

func (c *CommunityCards) GetCards() []*msgpb.Card {
	return c.Cards
}

func (c *CommunityCards) AddCard(card *msgpb.Card) {
	c.Cards = append(c.Cards, card)
}

func (c *CommunityCards) Reset() {
	c.Cards = []*msgpb.Card{}
}
