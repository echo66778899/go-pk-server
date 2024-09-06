package network

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	msgpb "go-pk-server/gen"
)

type Agent struct {
	ws    *websocket.Conn
	rxMsg chan *msgpb.ServerMessage
	txMsg chan *msgpb.ClientMessage

	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewAgent() *Agent {
	return &Agent{
		ws:    nil,
		rxMsg: make(chan *msgpb.ServerMessage, 30),
		txMsg: make(chan *msgpb.ClientMessage, 10),
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

	ctx := context.Background()
	ctx, a.cancel = context.WithCancel(ctx)
	go a.handleTxMsgRoutine(ctx, &a.wg)
	go a.handleRxMsgRoutine(ctx, &a.wg)

	return true
}

func (a *Agent) Close() {
	a.cancel()
	a.wg.Wait()

	if a.ws != nil {
		a.ws.Close()
		a.ws = nil
	}
}

func (a *Agent) SendingMessage(cm *msgpb.ClientMessage) {
	log.Printf("Sending message: %v", cm)
	a.txMsg <- cm
}

func (a *Agent) ReceivingMessage() chan *msgpb.ServerMessage {
	return a.rxMsg
}

func (a *Agent) handleTxMsgRoutine(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for {
		select {
		case msg := <-a.txMsg:
			e := a.sendMsg(msg)
			if e != nil {
				panic(e)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (a *Agent) handleRxMsgRoutine(ctx context.Context, wg *sync.WaitGroup) {
	//wg.Add(1)
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, e := a.receiveMsg()
		if e != nil {
			panic(e)
		}
		a.rxMsg <- msg
	}
}

func (a *Agent) sendMsg(msg *msgpb.ClientMessage) error {
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

func (a *Agent) receiveMsg() (*msgpb.ServerMessage, error) {
	// Use unblocking call
	//a.ws.SetReadDeadline(time.Now().Add(1 * time.Second))

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
