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
