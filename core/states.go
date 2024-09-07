package engine

// Game state is all things that we want to sync to the client for display and game logic

import (
	"encoding/json"
	"fmt"

	msgpb "go-pk-server/gen"
)

type PlayerStatus int

const (
	PlayerStatus_SatOut   PlayerStatus = 0 // Player is sat out and not playing
	PlayerStatus_SatIn    PlayerStatus = 1 // Player is sat in and waiting for the game to start
	PlayerStatus_Playing  PlayerStatus = 2 // Player is playing
	PlayerStatus_Wait4Act PlayerStatus = 3 // Player is waiting for their turn
	PlayerStatus_Check    PlayerStatus = 4 // Player has checked
	PlayerStatus_Call     PlayerStatus = 5 // Player has called
	PlayerStatus_Raise    PlayerStatus = 6 // Player has raised
	PlayerStatus_Fold     PlayerStatus = 7 // Player has folded
	PlayerStatus_AllIn    PlayerStatus = 8 // Player has all in
)

// String of PlayerStatus
func (ps PlayerStatus) String() string {
	return [...]string{"Sat Out", "Sat In", "Playing", "Wait4Act",
		"Check", "Call", "Raise", "Fold", "All In"}[ps]
}

// Game logic state
type GameState struct {
	pot Pot
	cc  CommunityCards

	// Table states
	ButtonPosition   int // slot number of the player who is the dealer
	CurrentRound     msgpb.RoundStateType
	CurrentBet       int
	NumPlayingPlayer int
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
