package engine

type PlayerActType int

const (
	Unknown PlayerActType = iota
	Fold
	Check
	Call
	Raise
	AllIn
)

func (pat PlayerActType) String() string {
	return [...]string{"UNKNOWN", "FOLD", "CHECK", "CALL", "RAISE", "ALL-IN"}[pat]
}

// Action represents a player's action in a Poker game.
type ActionIf interface {
	// For game engine
	FromWho() int
	WhatAction() PlayerActType
	HowMuch() int
}

type PlayerAction struct {
	// Common fields for all actions
	PlayerPosition int
	ActionType     PlayerActType
	Amount         int
}

func NewPlayerAction(position int, actionType PlayerActType, amount int) ActionIf {
	return &PlayerAction{
		PlayerPosition: position,
		ActionType:     actionType,
		Amount:         amount,
	}
}

func (pa *PlayerAction) FromWho() int {
	return pa.PlayerPosition
}

func (pa *PlayerAction) WhatAction() PlayerActType {
	return pa.ActionType
}

func (pa *PlayerAction) HowMuch() int {
	return pa.Amount
}
