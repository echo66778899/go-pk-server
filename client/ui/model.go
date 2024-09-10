package ui

import (
	msgpb "go-pk-server/gen"
	"sync"
)

const (
	ROOM_LOGIN = "room_login"
	IN_GAME    = "in_game"
)

// Model is a struct that holds data for the UI rendering.
type Model struct {
	MaxPlayers     int
	Players        []*msgpb.PlayerState
	CommunityCards []*msgpb.Card
	DealerPosition int
	CurrentRound   msgpb.RoundStateType
	Pot            int
	CurrentBet     int
	Result         *msgpb.Result
	PlayersBalance []string
	// Current user
	YourTablePosition   int
	YourLoginUsernameID string
	YourPrivateState    *msgpb.PeerState
	YourTurn            bool
	YourPlayerState     *msgpb.PlayerState

	// For Dealer Icon
	IsDealerVisible bool

	// For UI buttons control
	ActiveButtonMenu UIButtonMenuType
	IsButtonEnabled  bool
	IsButtonsVisible bool

	sync.Mutex
}

var UI_MODEL_DATA = Model{
	MaxPlayers: 0,
	Players:    make([]*msgpb.PlayerState, 0),
}
