package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	// Connect to the WebSocket server
	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer ws.Close()

	// Send a message to the server
	message := "Hello, Server!"
	err = ws.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	fmt.Printf("Sent: %s\n", message)

	// Read the server's response
	_, response, err := ws.ReadMessage()
	if err != nil {
		log.Fatalf("Failed to read message: %v", err)
	}

	fmt.Printf("Received: %s\n", response)

	// Keep the client running to allow for multiple exchanges (optional)
	for {
		fmt.Print("Enter a message: ")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Println("Failed to read input:", err)
			os.Exit(1)
		}

		err = ws.WriteMessage(websocket.TextMessage, []byte(input))
		if err != nil {
			log.Fatalf("Failed to send message: %v", err)
		}

		_, response, err = ws.ReadMessage()
		if err != nil {
			log.Fatalf("Failed to read message: %v", err)
		}

		fmt.Printf("Received: %s\n", response)
	}
}
