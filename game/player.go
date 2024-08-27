package engine

import "fmt"

// ActionType represents the type of action a player can make.
type ActionType int

const (
	// Fold represents the action of folding, i.e., discarding the hand and giving up on the current round.
	Fold ActionType = iota
	// Check represents the action of checking, i.e., not placing a bet and passing the turn to the next player.
	Check
	// Call represents the action of matching the current bet amount.
	Call
	// Bet represents the action of placing a bet.
	Bet
	// Raise represents the action of increasing the current bet amount.
	Raise
	// AllIn represents the action of betting all of the player's remaining chips.
	AllIn
)

// String returns the string representation of an ActionType.
func (a ActionType) String() string {
	return [...]string{"Fold", "Check", "Call", "Bet", "Raise", "AllIn"}[a]
}

// Action represents a player's action in a Poker game.
type PlayerAction struct {
	PlayerIdx int
	Type      ActionType // The type of action.
	Bet       int        // The amount of chips bet (if applicable).
}

type Player struct {
	PlayerIdx  int
	ID         int
	Name       string
	Chips      int
	Hand       Hand
	HasFolded  bool
	HasActed   bool
	CurrentBet int
}

// NewPlayer creates a new player with the given ID and name.
func NewPlayer(id int, name string) Player {
	return Player{
		ID:    id,
		Name:  name,
		Chips: 2000,
	}
}

// Print player's hand
func (p *Player) PrintHand() {
	fmt.Printf("Player %s's hand: [%s]\n", p.Name, p.Hand.String())
}

func (p *Player) PrintChips() {
	fmt.Printf("Player %s has %d chips\n", p.Name, p.Chips)
}
