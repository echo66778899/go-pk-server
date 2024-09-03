package main

import (
	"fmt"
	"log"
	"os"

	sync_msg "go-pk-server/msg"

	"github.com/gorilla/websocket"
)

type CMessage struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	// Connect to the WebSocket server
	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer ws.Close()

	// Send a message to the server
	msg := sync_msg.AuthenMessage{}
	{
		// Enter the authentication details
		fmt.Print("Enter your username: ")
		_, err := fmt.Scanln(&msg.Username)

		if err != nil {
			log.Println("Failed to read input:", err)
			os.Exit(1)
		}

		fmt.Print("Enter your room: ")
		_, err = fmt.Scanln(&msg.Room)

		if err != nil {
			log.Println("Failed to read input:", err)
			os.Exit(1)
		}

		fmt.Print("Enter your passcode: ")
		_, err = fmt.Scanln(&msg.Passcode)

		if err != nil {
			log.Println("Failed to read input:", err)
			os.Exit(1)
		}

		fmt.Print("Enter your session: ")
		_, err = fmt.Scanln(&msg.Session)

		if err != nil {
			log.Println("Failed to read input:", err)
			os.Exit(1)
		}
	}

	commMsg := sync_msg.CommunicationMessage{
		Type:    sync_msg.AuthMsgType,
		Payload: msg,
	}

	err = ws.WriteJSON(commMsg)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	fmt.Printf("Sent authen: %s\n", commMsg)

	// Read will not block if there is no message
	go func() {
		for {
			var msg sync_msg.CommunicationMessage
			err = ws.ReadJSON(&msg)
			if err != nil {
				log.Fatalf("Failed to read message: %v", err)
			}

			fmt.Printf("\nReceived: %s\n", msg)
		}
	}()

	var joinTable bool

	for {
		var msg sync_msg.CommunicationMessage
		fmt.Println("Menu Options:")
		if !joinTable {
			fmt.Println("0. Join a slot in table")
		}
		fmt.Println("1. Start Game")
		fmt.Println("2. Next Game")
		fmt.Println("3. Player Action")
		fmt.Println("4. Request Buy-in")
		fmt.Println("5. Exit")

		var choice int
		fmt.Print("Enter your choice: ")
		_, err := fmt.Scanln(&choice)
		if err != nil {
			log.Println("Failed to read input:", err)
			os.Exit(1)
		}

		switch choice {
		case 0:
			// Scan for slot number
			fmt.Print("Enter slot number: ")
			var slot string
			_, err := fmt.Scanln(&slot)
			if err != nil {
				log.Println("Failed to read input:", err)
				os.Exit(1)
			}
			msg.Type = sync_msg.CtrlMsgType
			msg.Payload = sync_msg.ControlMessage{ControlType: "join_slot", Data: slot}
			joinTable = true
		case 1:
			msg.Type = sync_msg.CtrlMsgType
			msg.Payload = sync_msg.ControlMessage{ControlType: "start_game", Data: ""}
		case 2:
			msg.Type = sync_msg.CtrlMsgType
			msg.Payload = sync_msg.ControlMessage{ControlType: "next_game", Data: ""}
		case 3:
			fmt.Print("Enter player action (call, check, fold, raise, allin): ")
			var action string
			_, err := fmt.Scanln(&action)
			if err != nil {
				log.Println("Failed to read input:", err)
				os.Exit(1)
			}
			switch action {
			case "call", "check", "fold", "allin":
				msg.Type = sync_msg.PlayerActMsgType
				msg.Payload = sync_msg.PlayerMessage{ActionName: action, Value: 0}
			case "raise":
				fmt.Print("Enter raise amount: ")
				var value int
				_, err := fmt.Scanln(&value)
				if err != nil {
					log.Println("Failed to read input:", err)
					os.Exit(1)
				}
			default:
				fmt.Println("Invalid action. Please try again.")
			}
		case 4:
			msg.Type = sync_msg.CtrlMsgType
			msg.Payload = sync_msg.ControlMessage{ControlType: "request_buyin", Data: ""}
		case 5:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}

		err = ws.WriteJSON(msg)
		if err != nil {
			log.Fatalf("Failed to send message: %v", err)
		}
		fmt.Printf("Sent: %v\n", msg)
	}
}
