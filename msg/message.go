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
