package engine

// Game state is all things that we want to sync to the client for display and game logic

import (
	"encoding/json"
	"fmt"

	msgpb "go-pk-server/gen"
)

// Game logic state
type GameState struct {
	pot Pot
	cc  CommunityCards

	// Table states
	ButtonPosition   int // slot number of the player who is the dealer
	CurrentRound     msgpb.RoundStateType
	CurrentBet       int
	NumPlayingPlayer int

	// Result states
	FinalResult *msgpb.Result
}

// Log game state for development
func (gs *GameState) LogState() {
	jsonBytes, _ := gs.SerializeToJson()
	fmt.Printf("Game state in JSON:\n %s\n", jsonBytes)
}

// Serialize game state to json string
func (gs *GameState) SerializeToJson() (string, error) {
	jsonBytes, err := json.Marshal(gs)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
