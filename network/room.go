package network

import (
	msgpb "go-pk-server/gen"
	mylog "go-pk-server/log"
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	RoomId   int
	Passcode string
	AdminId  string

	// Reference to the room's game
	settingGetter   func() *msgpb.GameSetting
	gameStateGetter func() *msgpb.GameState

	// List of clients in the room
	people map[userId]*Client
	// Broadcast message to all clients in the room
	broadcastChan chan *msgpb.ServerMessage

	mtx sync.Mutex
}

func NewRoom(id int, pass string, admin string) *Room {
	return &Room{
		RoomId:   id,
		Passcode: pass,
		AdminId:  admin,
		people:   make(map[userId]*Client),
	}
}

func (r *Room) Serve() {
	mylog.Debug("Starting broadcast for group:", r.RoomId)
	r.broadcastChan = make(chan *msgpb.ServerMessage, 10)
	go r.handleBroadcastMessage()
}

func (r *Room) SetSettingGetter(getter func() *msgpb.GameSetting, stateGetter func() *msgpb.GameState) {
	r.settingGetter = getter
	r.gameStateGetter = stateGetter
}

// Satisfy the PublicRoom interface
func (r *Room) BroadcastMessageToYourRoom(msg *msgpb.ServerMessage) {
	r.NotifiesChanges(msg)
}

func (r *Room) AddClient(c *Client) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.people[userId(c.Username)] = c
	// Send the game setting to the client
	if r.settingGetter != nil {
		setting := r.settingGetter()
		err := c.send(&msgpb.ServerMessage{
			Message: &msgpb.ServerMessage_GameSetting{
				GameSetting: setting,
			},
		})
		if err != nil {
			mylog.Errorf("Failed to send game setting to %s: %v", c.Username, err)
		}
	}
	// Send the game state to the client
	if r.gameStateGetter != nil {
		state := r.gameStateGetter()
		err := c.send(&msgpb.ServerMessage{
			Message: &msgpb.ServerMessage_GameState{
				GameState: state,
			},
		})
		if err != nil {
			mylog.Errorf("Failed to send game state to %s: %v", c.Username, err)
		}
	}
}

func (r *Room) RemoveClient(conn *websocket.Conn) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	userNameId := ""
	for _, c := range r.people {
		if c.conn.ws == conn {
			userNameId = c.Username
			c.handleDisconnect()
			break
		}
	}
	delete(r.people, userId(userNameId))
	mylog.Infof("Client of [%s] disconnected from the room", userNameId)
}

func (r *Room) CheckPasscode(pass string) bool {
	return r.Passcode == pass
}

func (r *Room) CheckUsername(username userId) bool {
	_, ok := r.people[userId(username)]
	return ok
}

func (r *Room) UpdateConnection(username userId, wsconn *websocket.Conn) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	if c, ok := r.people[userId(username)]; ok {
		c.conn.ws = wsconn
	}
}

func (r *Room) NotifiesChanges(message *msgpb.ServerMessage) {
	if message == nil {
		return
	}
	// Log the message
	mylog.Debugf("Broadcast message to all player in room %d\n", r.RoomId)
	// Send the message
	r.broadcastChan <- message
}

func (r *Room) DirectNotify(nameId string, message *msgpb.ServerMessage) {
	// Log the message
	mylog.Debugf("Direct message to player %s in room %d\n", nameId, r.RoomId)
	// Send the message
	if c, ok := r.people[userId(nameId)]; ok {
		_ = c.send(message)
	}
}

func (r *Room) handleBroadcastMessage() {
	for {
		message := <-r.broadcastChan

		// Broadcast message to all clients in the group
		mylog.Debugf("Broadcasting message to group %d: %v\n", r.RoomId, message)

		// for client in room
		for _, c := range r.people {
			err := c.send(message)

			if err != nil {
				mylog.Errorf("Error broadcasting message to %s:%v", c.Username, err)
			}
		}
	}
}
