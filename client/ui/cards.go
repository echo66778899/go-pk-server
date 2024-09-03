// Description: This file contains the code for the cards that are displayed on the client side.
// This code is responsible for displaying the cards in the client.

package ui

import (
	"fmt"
	engine "go-pk-server/core"
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
			fmt.Println() // Print an empty line between cards
		}
	}
}

func PrintBoard(cards []engine.Card) {
	// Print the board cards
	fmt.Println("=====================================================")
	defer fmt.Println("=====================================================")

	firstCard := []string{
		"┌─────┐",
		fmt.Sprintf("│%-2s   │", ranks[cards[0].Value]),
		fmt.Sprintf("│  %s%s%s  │", suitsColor[cards[0].Suit], suitsIcon[cards[0].Suit], reset),
		"└─────┘",
	}

	secondCard := []string{
		"┌─────┐",
		fmt.Sprintf("│%-2s   │", ranks[cards[1].Value]),
		fmt.Sprintf("│  %s%s%s  │", suitsColor[cards[1].Suit], suitsIcon[cards[1].Suit], reset),
		"└─────┘",
	}

	thirdCard := []string{
		"┌─────┐",
		fmt.Sprintf("│%-2s   │", ranks[cards[2].Value]),
		fmt.Sprintf("│  %s%s%s  │", suitsColor[cards[2].Suit], suitsIcon[cards[2].Suit], reset),
		"└─────┘",
	}

	fourthCard := []string{
		"┌─────┐",
		fmt.Sprintf("│%-2s   │", ranks[cards[3].Value]),
		fmt.Sprintf("│  %s%s%s  │", suitsColor[cards[3].Suit], suitsIcon[cards[3].Suit], reset),
		"└─────┘",
	}

	fifthCard := []string{
		"┌─────┐",
		fmt.Sprintf("│%-2s   │", ranks[cards[4].Value]),
		fmt.Sprintf("│  %s%s%s  │", suitsColor[cards[4].Suit], suitsIcon[cards[4].Suit], reset),
		"└─────┘",
	}

	// Print the board cards
	fmt.Println("Board Cards:")
	fmt.Println("-----------------------------------------------------")
	fmt.Println("  ", firstCard[0], "  ", secondCard[0], "  ", thirdCard[0], "  ", fourthCard[0], "  ", fifthCard[0])
	fmt.Println("  ", firstCard[1], "  ", secondCard[1], "  ", thirdCard[1], "  ", fourthCard[1], "  ", fifthCard[1])
	fmt.Println("  ", firstCard[2], "  ", secondCard[2], "  ", thirdCard[2], "  ", fourthCard[2], "  ", fifthCard[2])
	fmt.Println("  ", firstCard[3], "  ", secondCard[3], "  ", thirdCard[3], "  ", fourthCard[3], "  ", fifthCard[3])
	fmt.Println("-----------------------------------------------------")
}
