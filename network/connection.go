package network

import (
	"crypto/tls"
	msgpb "go-pk-server/gen"
	mylog "go-pk-server/log"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type userId string
type groupId uint64

type ConnectionManager struct {
	upgrader websocket.Upgrader

	rooms map[groupId]*Room
	mutex sync.Mutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for simplicity (be cautious in production)
			}},
		rooms: make(map[groupId]*Room),
	}
}

func (cm *ConnectionManager) AddRoom(gId groupId, r *Room) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	if r == nil {
		log.Fatalf("Room is nil")
	}
	cm.rooms[gId] = r
}

func (cm *ConnectionManager) StartServer(addr string) error {
	http.HandleFunc("/ws", cm.serveWebSocket)

	// Load the self-signed certificate and key
	certFile := "assets/cert.pem"
	keyFile := "assets/key.pem"

	// Configure the HTTPS server to use TLS
	server := &http.Server{
		Addr: ":8888", // HTTPS port (you can use :443 for production)
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12, // Use TLS 1.2 or higher
		},
	}

	if addr != "" {
		server.Addr = addr
	}

	mylog.Infof("Starting WSS server on https://%s\n", server.Addr)

	// Start the server with TLS using the self-signed certificate
	if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
		mylog.Fatalf("Failed to start server: %v", err)
	}
	return nil
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
	c, room := cm.ProcessCheckin(conn)
	if room == nil {
		mylog.Error("Failed to register client to a designated room")
		return
	}
	defer room.RemoveClient(conn)

	for {
		msgType, blob, err := c.conn.ws.ReadMessage()
		if err != nil {
			mylog.Errorf("Error reading message: %v", err)
			break
		}
		if msgType != websocket.BinaryMessage {
			mylog.Error("Invalid message type")
			continue
		}

		// Unmarshal the message
		message := &msgpb.ClientMessage{}
		if err := proto.Unmarshal(blob, message); err != nil {
			mylog.Errorf("Failed to unmarshal proto: %v", err)
			continue
		}

		// Dispatch message to the appropriate handler
		if c != nil {
			// Log received message
			mylog.Debugf("Received from player %s, message: %v", message.GetMessage(), c.Username)
			c.handleMessage(message)
		}
	}
}

// ==================
// Internal functions
// ==================

func (cm *ConnectionManager) ProcessCheckin(conn *websocket.Conn) (c *Client, r *Room) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	c, r = nil, nil

	// And timeout if no message is received in 10 seconds
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	msgType, blob, err := conn.ReadMessage()
	if err != nil {
		mylog.Error("Failed to read auth message:", err)
		return
	}

	if msgType != websocket.BinaryMessage {
		mylog.Error("Invalid authen message type")
		return
	}

	// Authenticate the client
	var message msgpb.ClientMessage

	if err := proto.Unmarshal(blob, &message); err != nil {
		log.Fatalf("Failed to unmarshal proto: %v", err)
		return
	}

	// Print the message
	mylog.Debug("Received authen message", message.GetMessage())

	jr, ok := message.GetMessage().(*msgpb.ClientMessage_JoinRoom)

	if !ok {
		mylog.Error("Invalid message ClientMessage_JoinRoom")
		return
	}

	roomStr := jr.JoinRoom.Room
	mylog.Debug("Room:", roomStr)

	// convert room string to uint64
	roomNo, err := strconv.ParseUint(roomStr, 10, 64)
	if err != nil {
		mylog.Errorf("Parsing invalid room number: %v", err.Error())
	}

	gId := groupId(roomNo)
	nameId := userId(jr.JoinRoom.NameId)
	passcode := jr.JoinRoom.Passcode
	session := jr.JoinRoom.SessionId

	// Check if the room is valid
	// Find the room
	room := cm.rooms[gId]
	if room == nil {
		mylog.Error("No existing room")
		return
	}

	// Check if the passcode is correct
	if room.CheckPasscode(passcode) == false {
		mylog.Error("Invalid passcode")
		return
	}

	// Check if username exists
	if room.CheckUsername(nameId) {
		if session != "" { // Reconnect
			// Log updated connection
			mylog.Infof("Player %s reconnected", nameId)
			room.UpdateConnection(nameId, conn)
			panic("Not expected updated connection")
		} else {
			mylog.Error("Username exists")
			return
		}
	}

	// Register new the client
	c = newConnectedClient(string(nameId), uint64(gId), room, conn)
	room.AddClient(c)
	r = room

	if r != nil && c != nil {
		mylog.Infof("Client at %v connected to room %d under a name [%s]", conn.RemoteAddr().String(), r.RoomId, nameId)
		conn.SetReadDeadline(time.Time{}) // Reset the read deadline
	}
	return c, r
}
