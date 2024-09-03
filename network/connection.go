package network

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type ConnectionManager struct {
	upgrader websocket.Upgrader
	clients  map[*websocket.Conn]bool
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		upgrader: websocket.Upgrader{},
		clients:  make(map[*websocket.Conn]bool),
	}
}

func (cm *ConnectionManager) ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := cm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}

	cm.clients[conn] = true

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			delete(cm.clients, conn)
			break
		}
	}

	conn.Close()
}

func (cm *ConnectionManager) StartServer(addr string) error {
	http.HandleFunc("/ws", cm.ServeWebSocket)
	return http.ListenAndServe(addr, nil)
}

// Broadcast sends a message to all connected clients.
func (cm *ConnectionManager) Broadcast(message []byte) {
	for conn := range cm.clients {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Failed to write message:", err)
			conn.Close()
			delete(cm.clients, conn)
		}
	}
}

// ReceiveMessages receives messages from all connected clients.
func (cm *ConnectionManager) ReceiveMessages() {
	for conn := range cm.clients {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read message:", err)
			conn.Close()
			delete(cm.clients, conn)
		}

		log.Println("Received message:", string(message))
	}
}
