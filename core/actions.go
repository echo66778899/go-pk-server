package engine

import "fmt"

// Action represents a player's action in a Poker game.
type ActionIf interface {
	// For game engine
	Execute(gs *GameState)
	What() string
}

// ========================================
// Action represents a player's action in a Poker game.
type FoldAction struct {
	Position int // The position of the player in the table making the action.
	Bet      int // The amount of chips bet (if applicable).
}

// NewPlayerAction creates a new PlayerAction.
func NewFoldAction(position, bet int) ActionIf {
	return &FoldAction{
		Position: position,
		Bet:      bet,
	}
}

// Execute executes the player's action.
func (pa *FoldAction) Execute(gs *GameState) {
	player := gs.GetPlayerByPosition(pa.Position)
	player.UpdateStatus(Folded)
	// Change state
	gs.NumPlayingPlayer--
	if np := gs.NextActivePlayer(pa.Position); np != nil {
		np.UpdateStatus(WaitForAct)
	}
}

// What returns the name of the action.
func (pa *FoldAction) What() string {
	return "Fold"
}

// ========================================
// Action represents a player's action in a Poker game.
type CheckAction struct {
	Position int // The position of the player in the table making the action.
}

func NewCheckAction(position int) ActionIf {
	return &CheckAction{
		Position: position,
	}
}

func (pa *CheckAction) Execute(gs *GameState) {
	player := gs.GetPlayerByPosition(pa.Position)
	// Check if the player is allowed to check
	if player.CurrentBet() == gs.CurrentBet {
		player.UpdateStatus(Checked)
		if np := gs.NextActivePlayer(pa.Position); np != nil {
			np.UpdateStatus(WaitForAct)
		}
	} else {
		// Log info the player name is not allowed to check, the action is invalid
		fmt.Println("Player", player.Name(), "is not allowed to check, the action is invalid")
	}
}

func (pa *CheckAction) What() string {
	return "Check"
}

// ========================================
// Action represents a player's action in a Poker game.
type CallAction struct {
	Position int // The position of the player in the table making the action.
}

func NewCallAction(position int) ActionIf {
	return &CallAction{
		Position: position,
	}
}

func (pa *CallAction) Execute(gs *GameState) {
	player := gs.GetPlayerByPosition(pa.Position)
	// If the player chip is less than the current bet, the player is all-in
	if player.Chips() > gs.CurrentBet {
		callChip := gs.CurrentBet - player.CurrentBet()
		player.UpdateCurrentBet(callChip)
		player.TakeChips(callChip)
		player.UpdateStatus(Called)
		gs.Pots.AddToPot(callChip)
		gs.CurrentBet = player.CurrentBet()
		if np := gs.NextActivePlayer(pa.Position); np != nil {
			np.UpdateStatus(WaitForAct)
		}
	} else {
		// Log warning the player should go all-in
		fmt.Println("Player", player.Name(), "should go all-in rather than call")
	}
}

func (pa *CallAction) What() string {
	return "Call"
}

// ========================================
// Action represents a player's action in a Poker game.
type BetAction struct {
	Position int // The position of the player in the table making the action.
	Bet      int // The amount of chips bet (if applicable).
}

func NewBetAction(position, bet int) ActionIf {
	return &BetAction{
		Position: position,
		Bet:      bet,
	}
}

func (pa *BetAction) Execute(gs *GameState) {
	player := gs.GetPlayerByPosition(pa.Position)
	betAmount := pa.Bet
	// If the player chip is less than the current bet, the player is all-in
	if player.Chips() > gs.CurrentBet/2 {
		betAmount += gs.CurrentBet - player.CurrentBet()
		player.UpdateCurrentBet(betAmount)
		player.TakeChips(betAmount)
		player.UpdateStatus(Betted)
		gs.Pots.AddToPot(betAmount)
		gs.CurrentBet = player.CurrentBet()

		// Update all player status to WaitForAct
		for _, p := range gs.Players {
			if p.Position() != pa.Position && p.Status() != Folded {
				p.UpdateStatus(WaitForAct)
			}
		}
	} else {
		// Log warning the player should go all-in
		fmt.Println("Player", player.Name(), "should go all-in rather than bet")
	}
}

func (pa *BetAction) What() string {
	return "Bet"
}

// ========================================
// Action represents a player's action in a Poker game.
type RaiseAction struct {
	Position int // The position of the player in the table making the action.
	Bet      int // The amount of chips bet (if applicable).
}

func NewRaiseAction(position, bet int) ActionIf {
	return &RaiseAction{
		Position: position,
		Bet:      bet,
	}
}

func (pa *RaiseAction) Execute(gs *GameState) {
	player := gs.GetPlayerByPosition(pa.Position)
	raiseAmount := pa.Bet
	// If the player chip is less than the current bet, the player is all-in
	if player.Chips() > gs.CurrentBet/2 {
		raiseAmount += gs.CurrentBet - player.CurrentBet()
		player.UpdateCurrentBet(raiseAmount)
		player.TakeChips(raiseAmount)
		player.UpdateStatus(Raised)
		gs.Pots.AddToPot(raiseAmount)
		gs.CurrentBet = player.CurrentBet()

		// Update all player status to WaitForAct
		for _, p := range gs.Players {
			if p.Position() != pa.Position && p.Status() != Folded {
				p.UpdateStatus(WaitForAct)
			}
		}
	} else {
		// Log warning the player should go all-in
		fmt.Println("Player", player.Name(), "should go all-in rather than Raise")
	}
}

func (pa *RaiseAction) What() string {
	return "Raise"
}

// ========================================
// Action represents a player's action in a Poker game.
type AllInAction struct {
	Position int // The position of the player in the table making the action.
}

func NewAllInAction(position int) ActionIf {
	return &AllInAction{
		Position: position,
	}
}

func (pa *AllInAction) Execute(gs *GameState) {
	player := gs.GetPlayerByPosition(pa.Position)
	allInAmount := player.Chips()
	if allInAmount > gs.CurrentBet {
		player.UpdateCurrentBet(allInAmount + player.CurrentBet())
		player.TakeChips(allInAmount)
		player.UpdateStatus(AlledIn)
		gs.Pots.AddToPot(allInAmount)
		gs.CurrentBet = player.CurrentBet()

		// Update all player status to WaitForAct
		for _, p := range gs.Players {
			if p.Position() != pa.Position && p.Status() != Folded {
				p.UpdateStatus(WaitForAct)
			}
		}
	} else {
		// Have to slip side pot
	}
}

func (pa *AllInAction) What() string {
	return "AllIn"
}
