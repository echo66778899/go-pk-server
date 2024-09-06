package engine

// Game state is all things that we want to sync to the client for display and game logic

import (
	"encoding/json"
	"fmt"

	msgpb "go-pk-server/gen"
)

type PlayerStatus int

const (
	SatOut   PlayerStatus = iota // Player is sat out and not playing
	SatIn                        // Player is sat in and waiting for the game to start
	Playing                      // Player is playing
	Wait4Act                     // Player is waiting for their turn
	Checked                      // Player has checked
	Called                       // Player has called
	Raised                       // Player has raised
	Folded                       // Player has folded
	AlledIn                      // Player has all in
)

// String of PlayerStatus
func (ps PlayerStatus) String() string {
	return [...]string{"SatOut", "SatIn", "Playing", "Wait4Act", "Checked",
		"Called", "Raised", "Folded", "AlledIn"}[ps]
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
