package engine

import "fmt"

// Action represents a player's action in a Poker game.
type ActionIf interface {
	// For game engine
	Execute(gs *GameState, pm *PlayerManager)
	What() string
}

// ========================================
// Action represents a player's action in a Poker game.
type FoldAction struct {
	Position   int // The position of the player in the table making the action.
	Desciption string
}

// NewPlayerAction creates a new PlayerAction.
func NewFoldAction(position int) ActionIf {
	return &FoldAction{
		Position:   position,
		Desciption: "FOLD",
	}
}

// Execute executes the player's action.
func (pa *FoldAction) Execute(gs *GameState, pm *PlayerManager) {
	// Log executing fold action from player
	fmt.Println("Executing FOLD action from player", pm.GetPlayer(pa.Position).Name())
	player := pm.GetPlayer(pa.Position)
	if player.Status() != WaitForAct {
		// Log warning the player is not allowed to fold, the action is invalid
		fmt.Println("Player", player.Name(), "is not allowed to ", pa.Desciption, ", the action is invalid")
		return
	}
	player.UpdateStatus(Folded)
	// Change state
	gs.NumPlayingPlayer--
	if np := pm.NextPlayer(pa.Position, Active); np != nil {
		np.UpdateStatus(WaitForAct)
	}
}

// What returns the name of the action.
func (pa *FoldAction) What() string {
	return "FOLD"
}

// ========================================
// Action represents a player's action in a Poker game.
type CheckAction struct {
	Position   int // The position of the player in the table making the action.
	Desciption string
}

func NewCheckAction(position int) ActionIf {
	return &CheckAction{
		Position:   position,
		Desciption: "CHECK",
	}
}

func (pa *CheckAction) Execute(gs *GameState, pm *PlayerManager) {
	// Log executing check action from player
	fmt.Println("Executing CHECK action from player", pm.GetPlayer(pa.Position).Name())

	player := pm.GetPlayer(pa.Position)
	if player.Status() != WaitForAct {
		// Log warning the player is not allowed to fold, the action is invalid
		fmt.Println("Player", player.Name(), "is not allowed to ", pa.Desciption, ", the action is invalid")
		return
	}
	// Check if the player is allowed to check
	if player.CurrentBet() == gs.CurrentBet {
		player.UpdateStatus(Checked)
		if np := pm.NextPlayer(pa.Position, Active); np != nil {
			np.UpdateStatus(WaitForAct)
		}
	} else {
		// Log info the player name is not allowed to check, the action is invalid
		fmt.Println("Player", player.Name(), "is not allowed to check, the action is invalid")
	}
}

func (pa *CheckAction) What() string {
	return "CHECK"
}

// ========================================
// Action represents a player's action in a Poker game.
type CallAction struct {
	Position   int // The position of the player in the table making the action.
	Desciption string
}

func NewCallAction(position int) ActionIf {
	return &CallAction{
		Position:   position,
		Desciption: "CALL",
	}
}

func (pa *CallAction) Execute(gs *GameState, pm *PlayerManager) {
	player := pm.GetPlayer(pa.Position)
	// Log executing call action from player
	fmt.Println("Executing CALL action from player ", player.Name(), " with status ", player.Status())

	if player.Status() != WaitForAct {
		// Log warning the player is not allowed to fold, the action is invalid
		fmt.Println("Player", player.Name(), "is not allowed to ", pa.Desciption, ", the action is invalid")
		return
	}
	// If the player chip is less than the current bet, the player is all-in
	callChip := gs.CurrentBet - player.CurrentBet()
	player.UpdateCurrentBet(callChip)
	player.TakeChips(callChip)
	player.UpdateStatus(Called)
	gs.pot.AddToPot(player.Position(), callChip)
	gs.CurrentBet = player.CurrentBet()
	if np := pm.NextPlayer(pa.Position, Active); np != nil {
		np.UpdateStatus(WaitForAct)
	}
}

func (pa *CallAction) What() string {
	return "CALL"
}

// ========================================
// Action represents a player's action in a Poker game.
type RaiseAction struct {
	Position   int // The position of the player in the table making the action.
	Bet        int // The amount of chips bet (if applicable).
	Desciption string
}

func NewRaiseAction(position, bet int) ActionIf {
	return &RaiseAction{
		Position:   position,
		Bet:        bet,
		Desciption: "RAISE",
	}
}

func (pa *RaiseAction) Execute(gs *GameState, pm *PlayerManager) {
	// Log
	fmt.Println("Executing RAISE action from player", pm.GetPlayer(pa.Position).Name())

	player := pm.GetPlayer(pa.Position)
	if player.Status() != WaitForAct {
		// Log warning the player is not allowed to fold, the action is invalid
		fmt.Println("Player", player.Name(), "is not allowed to ", pa.Desciption, ", the action is invalid")
		return
	}
	raiseAmount := pa.Bet
	// If the player chip is less than the current bet, the player is all-in
	if player.Chips() > gs.CurrentBet/2 {
		raiseAmount += gs.CurrentBet - player.CurrentBet()
		player.UpdateCurrentBet(raiseAmount)
		player.TakeChips(raiseAmount)
		player.UpdateStatus(Raised)
		gs.pot.AddToPot(player.Position(), raiseAmount)
		gs.CurrentBet = player.CurrentBet()

		// Update all player status to Active and NextPlayer to WaitForAct
		for _, p := range pm.GetListOfOtherPlayers(pa.Position, Called, Raised, Checked) {
			p.UpdateStatus(Active)
		}
		if np := pm.NextPlayer(pa.Position, Active); np != nil {
			np.UpdateStatus(WaitForAct)
		}
	} else {
		// Log warning the player should go all-in
		fmt.Println("Player", player.Name(), "should go all-in rather than Raise")
	}
}

func (pa *RaiseAction) What() string {
	return "RAISE"
}

// ========================================
// Action represents a player's action in a Poker game.
type AllInAction struct {
	Position   int // The position of the player in the table making the action.
	Desciption string
}

func NewAllInAction(position int) ActionIf {
	return &AllInAction{
		Position:   position,
		Desciption: "ALL-IN",
	}
}

func (pa *AllInAction) Execute(gs *GameState, pm *PlayerManager) {
	// Log executing all-in action from player
	fmt.Println("Executing ALL-IN action from player", pm.GetPlayer(pa.Position).Name())

	player := pm.GetPlayer(pa.Position)
	if player.Status() != WaitForAct {
		// Log warning the player is not allowed to fold, the action is invalid
		fmt.Println("Player", player.Name(), "is not allowed to ", pa.Desciption, ", the action is invalid")
		return
	}
	allInAmount := player.Chips()
	if allInAmount > gs.CurrentBet {
		player.UpdateCurrentBet(allInAmount + player.CurrentBet())
		player.TakeChips(allInAmount)
		player.UpdateStatus(AlledIn)
		gs.pot.AddToPot(player.Position(), allInAmount)
		gs.CurrentBet = player.CurrentBet()

		// Update all player status to Active and NextPlayer to WaitForAct
		for _, p := range pm.GetListOfOtherPlayers(pa.Position, Called, Raised, Checked) {
			p.UpdateStatus(Active)
		}
		if np := pm.NextPlayer(pa.Position, Active); np != nil {
			np.UpdateStatus(WaitForAct)
		}
	} else {
		// Have to slip side pot
	}
}

func (pa *AllInAction) What() string {
	return "ALLIN"
}
