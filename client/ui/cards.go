// Description: This file contains the code for the commuCards that are displayed on the client side.
// This code is responsible for displaying the commuCards in the client.

package ui

import (
	"fmt"
	engine "go-pk-server/core"
	"go-pk-server/msg"
)

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	black  = "\033[30m"
	white  = "\033[37m"
	green  = "\033[32m"
	yellow = "\033[33m"
)

var suitsColor = map[engine.Suit]string{
	engine.Hearts:   red,
	engine.Diamonds: yellow,
	engine.Clubs:    green,
	engine.Spades:   white,
}

var suitsIcon = map[engine.Suit]string{
	engine.Hearts:   "♥",
	engine.Diamonds: "♦",
	engine.Clubs:    "♣",
	engine.Spades:   "♠",
}

var ranks = map[engine.Value]string{
	engine.Ace:   "A",
	engine.Two:   "2",
	engine.Three: "3",
	engine.Four:  "4",
	engine.Five:  "5",
	engine.Six:   "6",
	engine.Seven: "7",
	engine.Eight: "8",
	engine.Nine:  "9",
	engine.Ten:   "10",
	engine.Jack:  "J",
	engine.Queen: "Q",
	engine.King:  "K",
}

var emptyCard = []string{
	"┌─────┐",
	"│.....│",
	"│.....│",
	"│.....│",
	"└─────┘",
}

func TestPrintSuits() {

	ranks := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

	for suit, color := range suitsColor {
		for _, rank := range ranks {
			card := []string{
				"┌─────┐",
				fmt.Sprintf("│%-2s   │", rank),
				fmt.Sprintf("│  %s%s%s  │", color, suit, reset),
				"└─────┘",
			}
			for _, line := range card {
				fmt.Println(line)
			}
			fmt.Println() // Print an empty line between cardsa
		}
	}
}

func PrintBoardFromGameSyncState(gameStateMsg *msg.CommunicationMessage) {
	if gameStateMsg == nil || gameStateMsg.Payload == nil {
		fmt.Println("Invalid game state message")
		return
	}
	msgPayload, ok := gameStateMsg.Payload.(map[string]interface{})
	if !ok {
		return
	}

	// Get the community cards
	communityCardsMsg, ok := msgPayload["community_cards"].([]interface{})
	if !ok {
		return
	}
	communityCards := []engine.Card{}
	for _, card := range communityCardsMsg {
		cardMap, ok := card.(map[string]interface{})
		if !ok {
			continue
		}
		suit, ok := cardMap["suit"].(int)
		if !ok {
			continue
		}
		value, ok := cardMap["value"].(int)
		if !ok {
			continue
		}
		communityCards = append(communityCards, engine.Card{
			Suit:  engine.Suit(suit),
			Value: engine.Value(value),
		})
	}

	playerHand, ok := msgPayload["player_hand"].([]interface{})
	if !ok {
		return
	}
	playerCards := []engine.Card{}
	for _, card := range playerHand {
		cardMap, ok := card.(map[string]interface{})
		if !ok {
			continue
		}
		suit, ok := cardMap["suit"].(int)
		if !ok {
			continue
		}
		value, ok := cardMap["value"].(int)
		if !ok {
			continue
		}
		playerCards = append(playerCards, engine.Card{
			Suit:  engine.Suit(suit),
			Value: engine.Value(value),
		})
	}

	PrintBoard(communityCards, playerCards)
}

func PrintBoard(commuCards []engine.Card, playerCards []engine.Card) {
	// Print the board commuCards
	fmt.Println("===================================================")
	defer fmt.Println("===================================================")

	// Initialize the commuCards with empty commuCards
	firstCard := emptyCard
	secondCard := emptyCard
	thirdCard := emptyCard
	fourthCard := emptyCard
	fifthCard := emptyCard

	yourFirstCard := emptyCard
	yourSecondCard := emptyCard

	// Base on length of commuCards, print the commuCards. If the length is less than 5, print empty commuCards
	switch len(commuCards) {
	case 0:
		// Print 5 empty commuCards
	case 1:
		firstCard = []string{
			"┌─────┐",
			fmt.Sprintf("│%-2s   │", ranks[commuCards[0].Value]),
			fmt.Sprintf("│  %s%s%s  │", suitsColor[commuCards[0].Suit], suitsIcon[commuCards[0].Suit], reset),
			fmt.Sprintf("│   %2s│", ranks[commuCards[0].Value]),
			"└─────┘",
		}
	case 2:
		firstCard = []string{
			"┌─────┐",
			fmt.Sprintf("│%-2s   │", ranks[commuCards[0].Value]),
			fmt.Sprintf("│  %s%s%s  │", suitsColor[commuCards[0].Suit], suitsIcon[commuCards[0].Suit], reset),
			fmt.Sprintf("│   %2s│", ranks[commuCards[0].Value]),
			"└─────┘",
		}
		secondCard = []string{
			"┌─────┐",
			fmt.Sprintf("│%-2s   │", ranks[commuCards[1].Value]),
			fmt.Sprintf("│  %s%s%s  │", suitsColor[commuCards[1].Suit], suitsIcon[commuCards[1].Suit], reset),
			fmt.Sprintf("│   %2s│", ranks[commuCards[1].Value]),
			"└─────┘",
		}
	case 3:
		firstCard = []string{
			"┌─────┐",
			fmt.Sprintf("│%-2s   │", ranks[commuCards[0].Value]),
			fmt.Sprintf("│  %s%s%s  │", suitsColor[commuCards[0].Suit], suitsIcon[commuCards[0].Suit], reset),
			fmt.Sprintf("│   %2s│", ranks[commuCards[0].Value]),
			"└─────┘",
		}
		secondCard = []string{
			"┌─────┐",
			fmt.Sprintf("│%-2s   │", ranks[commuCards[1].Value]),
			fmt.Sprintf("│  %s%s%s  │", suitsColor[commuCards[1].Suit], suitsIcon[commuCards[1].Suit], reset),
			fmt.Sprintf("│   %2s│", ranks[commuCards[1].Value]),
			"└─────┘",
		}
		thirdCard = []string{
			"┌─────┐",
			fmt.Sprintf("│%-2s   │", ranks[commuCards[2].Value]),
			fmt.Sprintf("│  %s%s%s  │", suitsColor[commuCards[2].Suit], suitsIcon[commuCards[2].Suit], reset),
			fmt.Sprintf("│   %2s│", ranks[commuCards[2].Value]),
			"└─────┘",
		}
	case 4:
		firstCard = []string{
			"┌─────┐",
			fmt.Sprintf("│%-2s   │", ranks[commuCards[0].Value]),
			fmt.Sprintf("│  %s%s%s  │", suitsColor[commuCards[0].Suit], suitsIcon[commuCards[0].Suit], reset),
			fmt.Sprintf("│   %2s│", ranks[commuCards[0].Value]),
			"└─────┘",
		}
		secondCard = []string{
			"┌─────┐",
			fmt.Sprintf("│%-2s   │", ranks[commuCards[1].Value]),
			fmt.Sprintf("│  %s%s%s  │", suitsColor[commuCards[1].Suit], suitsIcon[commuCards[1].Suit], reset),
			fmt.Sprintf("│   %2s│", ranks[commuCards[1].Value]),
			"└─────┘",
		}
		thirdCard = []string{
			"┌─────┐",
			fmt.Sprintf("│%-2s   │", ranks[commuCards[2].Value]),
			fmt.Sprintf("│  %s%s%s  │", suitsColor[commuCards[2].Suit], suitsIcon[commuCards[2].Suit], reset),
			fmt.Sprintf("│   %2s│", ranks[commuCards[2].Value]),
			"└─────┘",
		}
		fourthCard = []string{
			"┌─────┐",
			fmt.Sprintf("│%-2s   │", ranks[commuCards[3].Value]),
			fmt.Sprintf("│  %s%s%s  │", suitsColor[commuCards[3].Suit], suitsIcon[commuCards[3].Suit], reset),
			fmt.Sprintf("│   %2s│", ranks[commuCards[3].Value]),
			"└─────┘",
		}
	case 5:
		firstCard = []string{
			"┌─────┐",
			fmt.Sprintf("│%-2s   │", ranks[commuCards[0].Value]),
			fmt.Sprintf("│  %s%s%s  │", suitsColor[commuCards[0].Suit], suitsIcon[commuCards[0].Suit], reset),
			fmt.Sprintf("│   %2s│", ranks[commuCards[0].Value]),
			"└─────┘",
		}
		secondCard = []string{
			"┌─────┐",
			fmt.Sprintf("│%-2s   │", ranks[commuCards[1].Value]),
			fmt.Sprintf("│  %s%s%s  │", suitsColor[commuCards[1].Suit], suitsIcon[commuCards[1].Suit], reset),
			fmt.Sprintf("│   %2s│", ranks[commuCards[1].Value]),
			"└─────┘",
		}
		thirdCard = []string{
			"┌─────┐",
			fmt.Sprintf("│%-2s   │", ranks[commuCards[2].Value]),
			fmt.Sprintf("│  %s%s%s  │", suitsColor[commuCards[2].Suit], suitsIcon[commuCards[2].Suit], reset),
			fmt.Sprintf("│   %2s│", ranks[commuCards[2].Value]),
			"└─────┘",
		}
		fourthCard = []string{
			"┌─────┐",
			fmt.Sprintf("│%-2s   │", ranks[commuCards[3].Value]),
			fmt.Sprintf("│  %s%s%s  │", suitsColor[commuCards[3].Suit], suitsIcon[commuCards[3].Suit], reset),
			fmt.Sprintf("│   %2s│", ranks[commuCards[3].Value]),
			"└─────┘",
		}
		fifthCard = []string{
			"┌─────┐",
			fmt.Sprintf("│%-2s   │", ranks[commuCards[4].Value]),
			fmt.Sprintf("│  %s%s%s  │", suitsColor[commuCards[4].Suit], suitsIcon[commuCards[4].Suit], reset),
			fmt.Sprintf("│   %2s│", ranks[commuCards[4].Value]),
			"└─────┘",
		}
	}

	// Print the board commuCards
	fmt.Println("Board Cards:")
	fmt.Println("---------------------------------------------------")
	fmt.Println(" ", firstCard[0], secondCard[0], thirdCard[0], " ", fourthCard[0], " ", fifthCard[0])
	fmt.Println(" ", firstCard[1], secondCard[1], thirdCard[1], " ", fourthCard[1], " ", fifthCard[1])
	fmt.Println(" ", firstCard[2], secondCard[2], thirdCard[2], " ", fourthCard[2], " ", fifthCard[2])
	fmt.Println(" ", firstCard[3], secondCard[3], thirdCard[3], " ", fourthCard[3], " ", fifthCard[3])
	fmt.Println(" ", firstCard[4], secondCard[4], thirdCard[4], " ", fourthCard[4], " ", fifthCard[4])
	fmt.Println("---------------------------------------------------")

	if len(playerCards) == 2 {
		yourFirstCard = []string{
			"┌─────┐",
			fmt.Sprintf("│%-2s   │", ranks[playerCards[0].Value]),
			fmt.Sprintf("│  %s%s%s  │", suitsColor[playerCards[0].Suit], suitsIcon[playerCards[0].Suit], reset),
			fmt.Sprintf("│   %2s│", ranks[playerCards[0].Value]),
			"└─────┘",
		}
		yourSecondCard = []string{
			"┌─────┐",
			fmt.Sprintf("│%-2s   │", ranks[playerCards[1].Value]),
			fmt.Sprintf("│  %s%s%s  │", suitsColor[playerCards[1].Suit], suitsIcon[playerCards[1].Suit], reset),
			fmt.Sprintf("│   %2s│", ranks[playerCards[1].Value]),
			"└─────┘",
		}
	}

	fmt.Println("Your Hand:")
	fmt.Println("---------------------------------------------------")
	fmt.Println(" ", yourFirstCard[0], yourSecondCard[0])
	fmt.Println(" ", yourFirstCard[1], yourSecondCard[1])
	fmt.Println(" ", yourFirstCard[2], yourSecondCard[2])
	fmt.Println(" ", yourFirstCard[3], yourSecondCard[3])
	fmt.Println(" ", yourFirstCard[4], yourSecondCard[4])
	fmt.Println("---------------------------------------------------")

}

func test() {
	your := []engine.Card{
		{Value: engine.Ace, Suit: engine.Hearts},
		{Value: engine.Jack, Suit: engine.Diamonds},
	}

	cards := []engine.Card{
		{Value: engine.Ace, Suit: engine.Clubs},
		{Value: engine.Jack, Suit: engine.Hearts},
		{Value: engine.Three, Suit: engine.Clubs},
		{Value: engine.Ten, Suit: engine.Diamonds},
		{Value: engine.Five, Suit: engine.Spades},
	}
	PrintBoard(cards, your)

	your = []engine.Card{
		{Value: engine.Ten, Suit: engine.Hearts},
		{Value: engine.Ten, Suit: engine.Clubs},
	}
	cards_1 := []engine.Card{
		{Value: engine.Ace, Suit: engine.Clubs},
		{Value: engine.Jack, Suit: engine.Hearts},
		{Value: engine.Ten, Suit: engine.Diamonds},
		{Value: engine.Five, Suit: engine.Spades},
	}
	PrintBoard(cards_1, your)

	your = []engine.Card{
		{Value: engine.Ace, Suit: engine.Spades},
		{Value: engine.King, Suit: engine.Spades},
	}
	cards_2 := []engine.Card{
		{Value: engine.Ace, Suit: engine.Clubs},
		{Value: engine.Jack, Suit: engine.Hearts},
		{Value: engine.Ten, Suit: engine.Diamonds},
	}
	PrintBoard(cards_2, your)

	your = []engine.Card{
		{Value: engine.Two, Suit: engine.Spades},
		{Value: engine.Seven, Suit: engine.Diamonds},
	}
	PrintBoard([]engine.Card{}, your)
}
