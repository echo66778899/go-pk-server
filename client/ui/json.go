package ui

import (
	"encoding/json"
	"fmt"
)

type GameState struct {
	// Define your game state struct here
}

func PrintGameState(jsonData []byte) error {
	var gameState GameState
	err := json.Unmarshal(jsonData, &gameState)
	if err != nil {
		return err
	}

	// Print the game state
	fmt.Printf("Game State: %+v\n", gameState)

	return nil
}
