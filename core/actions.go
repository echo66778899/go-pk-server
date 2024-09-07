package engine

import msgpb "go-pk-server/gen"

// Action represents a player's action in a Poker game.
type ActionIf interface {
	// For game engine
	FromWho() int
	WhatAction() msgpb.PlayerGameActionType
	HowMuch() int
}

type PlayerAction struct {
	// Common fields for all actions
	PlayerPosition int
	ActionType     msgpb.PlayerGameActionType
	Amount         int
}

func NewPlayerAction(position int, actionType msgpb.PlayerGameActionType, amount int) ActionIf {
	return &PlayerAction{
		PlayerPosition: position,
		ActionType:     actionType,
		Amount:         amount,
	}
}

func (pa *PlayerAction) FromWho() int {
	return pa.PlayerPosition
}

func (pa *PlayerAction) WhatAction() msgpb.PlayerGameActionType {
	return pa.ActionType
}

func (pa *PlayerAction) HowMuch() int {
	return pa.Amount
}
