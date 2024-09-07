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
	Players        []*msgpb.Player
	CommunityCards []*msgpb.Card
	DealerPosition int
	CurrentRound   msgpb.RoundStateType
	Pot            int
	// Current user
	YourTablePosition int
	YourUsernameID    string
	YourPrivateState  *msgpb.PeerState

	sync.Mutex
}

var UI_MODEL_DATA = Model{
	MaxPlayers: 0,
	Players:    make([]*msgpb.Player, 0),
}
