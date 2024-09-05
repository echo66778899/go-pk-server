package network

import (
	"errors"
	"log"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	msgpb "go-pk-server/gen"
)

type Agent struct {
	ws    *websocket.Conn
	rxMsg chan *msgpb.ServerMessage
}

func NewAgent() *Agent {
	return &Agent{
		ws:    nil,
		rxMsg: make(chan *msgpb.ServerMessage),
	}
}

func (a *Agent) Send(msg *msgpb.ClientMessage) error {
	// Serialize (marshal) the protobuf message
	sendData, err := proto.Marshal(msg)
	if err != nil {
		log.Fatalf("Failed to marshal proto: %v", err)
	}
	// Send the response
	if err := a.ws.WriteMessage(websocket.BinaryMessage, sendData); err != nil {
		log.Fatalf("Failed to write message: %v", err)
	}

	return nil
}

func (a *Agent) ReceivingMessage() chan *msgpb.ServerMessage {
	return a.rxMsg
}

func (a *Agent) Receive() (*msgpb.ServerMessage, error) {
	// Receive a message (blocking call)
	_, blob, err := a.ws.ReadMessage()
	if err != nil {
		log.Fatalf("Failed to read message: %v", err)
		return nil, errors.New("Failed to read message")
	}

	// Deserialize (unmarshal) the protobuf message
	var serverMsg msgpb.ServerMessage
	if err := proto.Unmarshal(blob, &serverMsg); err != nil {
		log.Fatalf("Failed to unmarshal proto: %v", err)
		return nil, errors.New("Failed to unmarshal proto")
	}

	return &serverMsg, nil
}

func (a *Agent) Close() {
	if a.ws != nil {
		a.ws.Close()
		a.ws = nil
	}
}

func (a *Agent) Connect() bool {
	// Connect to the WebSocket server
	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
		return false
	}
	a.ws = ws
	return true
}
