package network

import (
	mylog "go-pk-server/log"
	"go-pk-server/msg"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type userId string
type groupId uint64

type ConnectionManager struct {
	upgrader websocket.Upgrader

	rooms     map[groupId]map[userId]*Client
	broadcast map[groupId]chan msg.CommunicationMessage // Broadcast channel
	mutex     sync.Mutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for simplicity (be cautious in production)
			}},
		rooms:     make(map[groupId]map[userId]*Client),
		broadcast: make(map[groupId]chan msg.CommunicationMessage),
	}
}

func (cm *ConnectionManager) CreateRoom(gId groupId) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	if _, ok := cm.rooms[gId]; !ok {
		cm.rooms[gId] = make(map[userId]*Client)
	}
}

func (cm *ConnectionManager) StartServer(addr string) error {
	http.HandleFunc("/ws", cm.serveWebSocket)
	return http.ListenAndServe(addr, nil)
}

// Function to handle WebSocket connections
func (cm *ConnectionManager) serveWebSocket(w http.ResponseWriter, r *http.Request) {
	// Log the request
	mylog.Debug("Received connection from:", r.RemoteAddr)

	conn, err := cm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		mylog.Error("Failed to upgrade connection:", err)
		return
	}

	// Remove client when this function returns
	defer func() {
		err := conn.Close()
		if err != nil {
			mylog.Error(err)
		}
	}()

	// Handle Register new client to room
	c := cm.ProcessCheckin(conn)
	if c == nil {
		mylog.Error("Failed to register client")
		return
	}

	defer c.handleDisconnect()

	for {
		// Log waiting for message
		mylog.Debugf("Waiting for message from player [%s]", c.Username)
		var message msg.CommunicationMessage
		err := c.ws.ReadJSON(&message)
		if err != nil {
			mylog.Errorf("Error reading message: %v", err)
			cm.RemoveClient(c)
			break
		}
		mylog.Debugf("Received from player %s, message=%v\n", c.Username, message)
		// Dispatch message to the appropriate handler

		if c != nil {
			c.handleMessage(message)
		}
	}
}

func (cm *ConnectionManager) handleBroadcastMessage(gId groupId) {
	for {
		message := <-cm.broadcast[gId]

		// Broadcast message to all clients in the group
		mylog.Debugf("Broadcasting message to group %d: %v\n", gId, message)

		// for client in room
		for _, c := range cm.rooms[gId] {
			conn := c.ws
			if conn == nil {
				mylog.Errorf("Player %s not found\n", c.Username)
				continue
			}
			err := conn.WriteJSON(message)
			if err != nil {
				conn.Close()
				mylog.Error("Error broadcasting message:", err)
				c.ws = nil
				// Call client disconnect
				c.handleDisconnect()
			}
		}
	}
}

// ==================
// Internal functions
// ==================

func (cm *ConnectionManager) ProcessCheckin(conn *websocket.Conn) (client *Client) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	client = nil
	// Authenticate the client
	var message msg.CommunicationMessage

	// And timeout if no message is received in 10 seconds
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	err := conn.ReadJSON(&message)
	if err != nil {
		mylog.Error("Failed to read auth message:", err)
		return
	}

	if message.Type != msg.AuthMsgType || message.Payload == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Invalid message"))
		mylog.Error("Invalid authen message")
		return
	}

	// Print the message
	mylog.Debugf("Received authen message: %v\n", message)

	checkIn, ok := message.Payload.(map[string]interface{})
	if !ok {
		conn.WriteMessage(websocket.TextMessage, []byte("Invalid message"))
		mylog.Error("Invalid authen message")
		return
	}

	roomStr := checkIn["room"].(string)

	mylog.Debug("Room:", roomStr)
	// convert room string to uint64
	roomNo, err := strconv.ParseUint(roomStr, 10, 64)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Invalid room"))
		mylog.Errorf("Parsing invalid room number: %v", err.Error())
	}

	gId := groupId(roomNo)
	nameId := userId(checkIn["username"].(string))
	passcode := checkIn["passcode"].(string)
	session := checkIn["session"].(string)

	// Check if the room is valid
	// Find the room
	if cm.rooms[gId] == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Invalid room"))
		mylog.Error("No existing room")
		return
	}

	// Check if the passcode is correct
	if passcode != roomMap[gId] {
		conn.WriteMessage(websocket.TextMessage, []byte("Invalid passcode"))
		mylog.Error("Invalid passcode")
		return
	}

	// Check if username exists
	if cm.rooms[gId][nameId] != nil {
		if session != "" { // Reconnect
			// Log updated connection
			mylog.Infof("Player %s reconnected", nameId)
			client = cm.rooms[gId][nameId]
			client.ws = conn
		} else {
			conn.WriteMessage(websocket.TextMessage, []byte("Username exists"))
			mylog.Error("Username exists")
			return
		}
	}

	// Register new the client
	client = newConnectedClient(string(nameId), uint64(gId), cm, conn)
	cm.rooms[gId][nameId] = client

	if _, ok := cm.broadcast[gId]; !ok {
		cm.broadcast[gId] = make(chan msg.CommunicationMessage)
		mylog.Debug("Starting broadcast for group:", gId)
		go cm.handleBroadcastMessage(gId)
	}

	if client != nil {
		mylog.Infof("Client connected from %v. client username: %s", conn.RemoteAddr().String(), nameId)
		conn.SetReadDeadline(time.Time{}) // Reset the read deadline
	}
	return client
}

func (cm *ConnectionManager) RemoveClient(client *Client) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	// Log the disconnection
	mylog.Infof("Client from player %s disconnected", client.Username)

	delete(cm.rooms[groupId(client.GroupId)], userId(client.Username))

	mylog.Debugf("Number of clients in room %d: %d\n", client.GroupId, len(cm.rooms[groupId(client.GroupId)]))
}

func (cm *ConnectionManager) NotifiesChanges(gId uint64, message *msg.CommunicationMessage) {
	// Log the message
	mylog.Debugf("Sending message to player: %v\n", message)

	// Send the message
	cm.broadcast[groupId(gId)] <- *message
}

func (cm *ConnectionManager) BroadcastMessage(b []byte, gId groupId) {
	var message msg.CommunicationMessage
	message.Type = msg.SyncGameStateMsgType
	message.Payload = string(b)

	cm.broadcast[gId] <- message
}
