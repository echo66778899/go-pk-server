// Description: This file contains the code for the cards that are displayed on the client side.
// This code is responsible for displaying the cards in the client.

package ui

import "fmt"

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	black  = "\033[30m"
	white  = "\033[37m"
	green  = "\033[32m"
	yellow = "\033[33m"
)

func TestPrintSuits() {
	// Define colors for suits
	suits := map[string]string{
		"♥": red,
		"♦": yellow,
		"♣": green,
		"♠": white,
	}

	ranks := []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

	for suit, color := range suits {
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
