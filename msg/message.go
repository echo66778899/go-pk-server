package msg

type CommMessageType int // Communication message type between client and server

const (
	// CommMessage represents the communication message type.
	CtrlMsgType CommMessageType = iota
	// AuthenMessage represents the authentication message type.
	AuthMsgType
	// PlayerMessage represents the player action message type.
	PlayerActMsgType
	// SyncGameState represents the game state sync message type.
	SyncGameStateMsgType
	// SyncPriaveState represents the player state sync message type.
	SyncPrivateStateMsgType
	// ErrorMsgType represents the error message type.
	ErrorMsgType
)

// CommMessageTypeMap is a map of CommMessageType to string.
func (cmt CommMessageType) String() string {
	return [...]string{"CtrlMsgType", "AuthMsgType", "PlayerActMsgType", "SyncGameStateMsgType", "ErrorMsgType"}[cmt]
}

// CommunicationMessage represents the JSON structure for communication messages between client and server.
type CommunicationMessage struct {
	Type    CommMessageType `json:"type"`
	Payload interface{}     `json:"payload"`
}

// AuthenMessage represents the authentication JSON message from the client to the server.
type AuthenMessage struct {
	Username string `json:"username"`
	Room     string `json:"room"`
	Passcode string `json:"passcode"`
	Session  string `json:"session"` // Emtpty string for first time connection
}

// control message struct for start game, end game, etc.
type ControlMessage struct {
	ControlType string `json:"control_type"`
	Data        string `json:"data"`
}

// Player Action message struct
type PlayerMessage struct {
	ActionName string `json:"action_name"`
	Value      int    `json:"value"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type PlayerState struct {
	Name   string `json:"name"`
	Slot   int    `json:"slot"`
	Chips  int    `json:"chips"`
	Bet    int    `json:"bet"`
	Status string `json:"status"`
}

type PrivateMessage struct {
	Hand []Card `json:"your_hand"`
}

type Card struct {
	Suit  int `json:"suit"`
	Value int `json:"value"`
}

type SyncGameStateMessage struct {
	CommunityCards []Card        `json:"community_cards"`
	Players        []PlayerState `json:"players"`
	Pot            int           `json:"pot"`
	CurrentPlayer  string        `json:"current_player"`
}
