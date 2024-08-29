package engine

import (
	"testing"
)

func TestHand(t *testing.T) {
	// Create a new hand
	hand := Hand{}

	// Set the cards in the hand
	card1 := Card{Suit: Spades, Value: Ace}
	card2 := Card{Suit: Hearts, Value: King}
	hand.SetCard(card1, 0)
	hand.SetCard(card2, 1)

	// Calculate the best hand with community cards
	communityCards := &CommunityCards{
		Cards: []Card{
			{Suit: Spades, Value: Queen},
			{Suit: Hearts, Value: Jack},
			{Suit: Diamonds, Value: Ten},
			{Suit: Clubs, Value: Nine},
			{Suit: Spades, Value: Eight},
		},
	}
	bestRank := hand.CalcBestHand(communityCards)

	// Ensure that the best rank is correct
	expectedRank := Straight
	if bestRank != expectedRank {
		t.Errorf("Expected best rank to be %s, but got %s", expectedRank, bestRank)
	}

	// Ensure that the string representation of the hand is correct
	expectedString := "Ace of Spades, King of Hearts"
	handString := hand.String()
	if handString != expectedString {
		t.Errorf("Expected hand string to be %s, but got %s", expectedString, handString)
	}

	// Ensure that the string representation of the hand rank is correct
	expectedRankString := "Straight"
	rankString := hand.HandRankingString()
	if rankString != expectedRankString {
		t.Errorf("Expected hand rank string to be %s, but got %s", expectedRankString, rankString)
	}
}

func TestHandRanking(t *testing.T) {
	// Create a new hand rank
	rank := HandRanking(3)

	// Ensure that the string representation of the hand rank is correct
	expectedString := "ThreeOfAKind"
	rankString := rank.String()
	if rankString != expectedString {
		t.Errorf("Expected hand rank string to be %s, but got %s", expectedString, rankString)
	}
}

func TestCard(t *testing.T) {
	// Create a new card
	card := Card{Suit: Spades, Value: Ace}

	// Ensure that the string representation of the card is correct
	expectedString := "Ace of Spades"
	cardString := card.String()
	if cardString != expectedString {
		t.Errorf("Expected card string to be %s, but got %s", expectedString, cardString)
	}
}

func TestDeck(t *testing.T) {
	// Create a new deck
	deck := NewDeck()

	// Ensure that the deck has the correct number of cards
	expectedNumCards := 52
	numCards := len(deck.Cards)
	if numCards != expectedNumCards {
		t.Errorf("Expected number of cards to be %d, but got %d", expectedNumCards, numCards)
	}

	// Ensure that the deck has the correct number of dealed cards
	expectedNumDealed := 0
	numDealed := deck.Dealed
	if numDealed != expectedNumDealed {
		t.Errorf("Expected number of dealed cards to be %d, but got %d", expectedNumDealed, numDealed)
	}
}

func TestDeckShuffle(t *testing.T) {
	// Create a new deck
	deck := NewDeck()

	// Shuffle the deck
	deck.Shuffle()

	// Ensure that the deck has the correct number of cards
	expectedNumCards := 52
	numCards := len(deck.Cards)
	if numCards != expectedNumCards {
		t.Errorf("Expected number of cards to be %d, but got %d", expectedNumCards, numCards)
	}

	// Ensure that the deck has the correct number of dealed cards
	expectedNumDealed := 0
	numDealed := deck.Dealed
	if numDealed != expectedNumDealed {
		t.Errorf("Expected number of dealed cards to be %d, but got %d", expectedNumDealed, numDealed)
	}
}
