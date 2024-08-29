package network

import (
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	conn          *websocket.Conn
	sendQueue     chan []byte
	receiveQueue  chan []byte
	dispatchQueue chan []byte
	wg            sync.WaitGroup
}

func NewWebSocketClient(serverURL string) (*WebSocketClient, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	client := &WebSocketClient{
		conn:          conn,
		sendQueue:     make(chan []byte),
		receiveQueue:  make(chan []byte),
		dispatchQueue: make(chan []byte),
	}

	client.wg.Add(2)
	go client.sendRoutine()
	go client.receiveRoutine()

	return client, nil
}

func (c *WebSocketClient) sendRoutine() {
	defer c.wg.Done()

	for {
		select {
		case message := <-c.sendQueue:
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Failed to send message:", err)
			}
		}
	}
}

func (c *WebSocketClient) receiveRoutine() {
	defer c.wg.Done()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("Failed to receive message:", err)
			return
		}

		c.receiveQueue <- message
	}
}

func (c *WebSocketClient) DispatchRoutine() {
	for message := range c.receiveQueue {
		// Handle incoming message here
		fmt.Println("Received message:", string(message))
	}

	for message := range c.dispatchQueue {
		c.sendQueue <- message
	}
}

func (c *WebSocketClient) SendMessage(message []byte) {
	c.dispatchQueue <- message
}

func (c *WebSocketClient) Close() {
	close(c.sendQueue)
	close(c.receiveQueue)
	close(c.dispatchQueue)
	c.wg.Wait()
	c.conn.Close()
}
