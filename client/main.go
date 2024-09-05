package main

import (
	"fmt"
	"go-pk-server/client/handler"
	"go-pk-server/client/network"
	"go-pk-server/client/ui"
	msgpb "go-pk-server/gen"
	"log"
	"os"
)

var (
	i    int = 1000
	TEST     = true
	// pointsChan         = make(chan int)
	keyboardEventsChan = make(chan ui.KeyboardEvent)
)

var (
	table      *ui.Table
	board      *ui.Cards
	centerText *ui.TextBox
	dealerIcon *ui.DealerWidget
	playerWg   *ui.PlayersGroup
	btnCtrl    *ui.ButtonCtrlCenter
)

func testGame() {
	centerText.Text = fmt.Sprintf("i = %d", i)
	dealerIcon.IndexUI(i%6, 6)

	board.SetCards([]msgpb.Card{
		{Suit: msgpb.SuitType_HEARTS, Rank: msgpb.RankType_ACE},
		{Suit: msgpb.SuitType_SPADES, Rank: msgpb.RankType_KING},
		{Suit: msgpb.SuitType_DIAMONDS, Rank: msgpb.RankType_QUEEN},
		{Suit: msgpb.SuitType_CLUBS, Rank: msgpb.RankType_JACK},
		{Suit: msgpb.SuitType_HEARTS, Rank: msgpb.RankType_FIVE},
	})

	gs := msgpb.GameState{
		Players: []*msgpb.Player{
			{
				Name:          "player1",
				Chips:         1500,
				TablePosition: 0,
			},
			{
				Name:          "player2",
				Chips:         2000,
				TablePosition: 1,
			},
			{
				Name:          "player3",
				Chips:         4000,
				TablePosition: 2,
			},
			{
				Name:          "Hai",
				Chips:         3000,
				TablePosition: 4,
			},
		},
	}
	playerWg.UpdateState(gs.Players, true)
	if i%2 == 0 {
		playerWg.UpdatePocketPair(nil)
	} else {
		playerWg.UpdatePocketPair(&msgpb.PeerState{
			PlayerCards: []*msgpb.Card{
				&msgpb.Card{Suit: msgpb.SuitType_HEARTS, Rank: msgpb.RankType_DEUCE},
				&msgpb.Card{Suit: msgpb.SuitType_SPADES, Rank: msgpb.RankType_SEVEN},
			},
		})
	}
}

func initWindows() {
	// Table
	table = ui.NewTable()
	table.SetRect(ui.TABLE_CENTER_X-ui.TABLE_RADIUS_X, ui.TABLE_CENTER_Y-ui.TABLE_RADIUS_Y, 2*ui.TABLE_RADIUS_X, 2*ui.TABLE_RADIUS_Y)

	// New board of cards
	board = ui.NewCards()
	board.SetTitle("Community Cards")
	board.SetCoodinate(ui.COMMUNITY_CARDS_X, ui.COMMUNITY_CARDS_Y)

	// New central text box
	centerText = ui.NewParagraph()
	centerText.SetRect(ui.TABLE_CENTER_X-10, ui.TABLE_CENTER_Y+2, ui.TABLE_CENTER_X-10+21, ui.TABLE_CENTER_Y+5)

	// Dealer Icon
	dealerIcon = ui.NewDealerWidget()
	dealerIcon.SetRect(ui.TABLE_CENTER_X-10, ui.TABLE_CENTER_Y+5, ui.TABLE_CENTER_X-10+21, ui.TABLE_CENTER_Y+8)

	// Players slider
	playerWg = ui.NewPlayersGroup()
	playerWg.SetMaxOtherPlayers(6)

	// Buttons
	btnCtrl = ui.NewButtonCtrlCenter()
	btnCtrl.InitButtonPosition()
}

func render() {
	uiItems := []ui.Drawable{table, board, centerText, dealerIcon}
	uiItems = append(uiItems, playerWg.GetAllItems()...)
	uiItems = append(uiItems, btnCtrl.GetDisplayingButton()...)
	ui.Render(uiItems...)
}

func getRoomInfoInput() (playerName, room, passcode, sessId string) {
	if TEST {
		return "Hai Phan", "2222", "1234", "0"
	}

	// Enter the authentication details
	fmt.Print("Enter your name    : ")
	_, err := fmt.Scanln(&playerName)

	if err != nil {
		log.Println("Failed to read input:", err)
		os.Exit(1)
	}

	fmt.Print("Enter room number  : ")
	_, err = fmt.Scanln(&room)

	if err != nil {
		log.Println("Failed to read input:", err)
		os.Exit(1)
	}

	fmt.Print("Enter your passcode: ")
	_, err = fmt.Scanln(&passcode)

	if err != nil {
		log.Println("Failed to read input:", err)
		os.Exit(1)
	}

	fmt.Print("Enter your session : ")
	_, err = fmt.Scanln(&sessId)

	if err != nil {
		log.Println("Failed to read input:", err)
		os.Exit(1)
	}

	return
}

func buttonEnterEvtHandler(bt ui.ButtonType) {
	switch bt {
	case ui.BNT_FoldButton:
	case ui.BNT_CheckButton:
	case ui.BNT_CallButton:
	case ui.BNT_RaiseButton:
	case ui.BNT_AllInButton:
	case ui.BNT_JoinTableButton:
	case ui.BNT_StartGameButton:
	case ui.BNT_LeaveGameButton:
	case ui.BNT_RequestChipButton:
	default:
	}
}

func main() {
	// Connect to the WebSocket server
	agent := network.NewAgent()
	if agent.Connect() {
		// Send a message to the server
		n, r, p, s := getRoomInfoInput()
		msg := msgpb.ClientMessage{
			Message: &msgpb.ClientMessage_JoinRoom{
				JoinRoom: &msgpb.JoinRoom{
					NameId:    n,
					Room:      r,
					Passcode:  p,
					SessionId: s,
				},
			},
		}
		agent.Send(&msg)
	}
	defer agent.Close()

	// Initialize the UI
	ui.Init()
	defer ui.Deinit()

	// Initialize the windows and widgets
	initWindows()

	btnCtrl.SetOutsideHandler(buttonEnterEvtHandler)

	// new chan for keyboard
	go ui.ListenToKeyboard(keyboardEventsChan)

	// Game looop
	for {
		// Process input
		select {
		case ev := <-keyboardEventsChan:
			switch ev.EventType {
			case ui.LEFT:
				i--
				btnCtrl.MoveLeft()
			case ui.RIGHT:
				i++
				btnCtrl.MoveRight()
			case ui.UP:
				btnCtrl.MoveUp()
			case ui.DOWN:
				btnCtrl.MoveDown()
			case ui.ENTER:
				btnCtrl.Enter()
			case ui.SPACE:
				if btnCtrl.CtrlEnabled {
					btnCtrl.EnableButtonCtrl(false)
				} else {
					btnCtrl.EnableButtonCtrl(true)
				}
			case ui.BACKSPACE:
				// Toggle button control
				btnCtrl.ToggleMenu()
			case ui.END:
				return
			}
		case m := <-agent.ReceivingMessage():
			handler.HandleServerMessage(m)
		}

		// Update ux elements
		ui.CurrentPlayer.CurrentPlayerPossition = 4
		testGame()

		// Render ux elements
		render()
	}

	// // Connect to the WebSocket server
	// ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to server: %v", err)
	// }
	// defer ws.Close()

	// // Send a message to the server
	// msg := sync_msg.AuthenMessage{}
	// {

	// commMsg := sync_msg.CommunicationMessage{
	// 	Type:    sync_msg.AuthMsgType,
	// 	Payload: msg,
	// }

	// err = ws.WriteJSON(commMsg)
	// if err != nil {
	// 	log.Fatalf("Failed to send message: %v", err)
	// }

	// fmt.Printf("Sent authen: %s\n", commMsg)

	// // Read will not block if there is no message
	// go func() {
	// 	for {
	// 		var msg sync_msg.CommunicationMessage
	// 		err = ws.ReadJSON(&msg)
	// 		if err != nil {
	// 			log.Fatalf("Failed to read message: %v", err)
	// 		}

	// 		fmt.Printf("\nReceived: %s\n", msg)

	// 		if msg.Type == sync_msg.ErrorMsgType {
	// 			log.Fatalf("Error message received: %s", msg.Payload)
	// 		}

	// 		if msg.Type == sync_msg.SyncGameStateMsgType {
	// 			// Render the game state
	// 			fmt.Printf("Game State: %v\n", msg.Payload)
	// 			ui.PrintBoardFromGameSyncState(&msg)
	// 		}
	// 	}
	// }()

	// var joinTable bool

	// for {
	// 	var msg sync_msg.CommunicationMessage
	// 	fmt.Println("Menu Options:")
	// 	if !joinTable {
	// 		fmt.Println("0. Join a slot in table")
	// 	}
	// 	fmt.Println("1. Start Game")
	// 	fmt.Println("2. Next Game")
	// 	fmt.Println("3. Player Action")
	// 	fmt.Println("4. Request Buy-in")
	// 	fmt.Println("5. Exit")

	// 	var choice int
	// 	fmt.Print("Enter your choice: ")
	// 	_, err := fmt.Scanln(&choice)
	// 	if err != nil {
	// 		log.Println("Failed to read input:", err)
	// 		os.Exit(1)
	// 	}

	// 	switch choice {
	// 	case 0:
	// 		// Scan for slot number
	// 		fmt.Print("Enter slot number: ")
	// 		var slot string
	// 		_, err := fmt.Scanln(&slot)
	// 		if err != nil {
	// 			log.Println("Failed to read input:", err)
	// 			os.Exit(1)
	// 		}
	// 		msg.Type = sync_msg.CtrlMsgType
	// 		msg.Payload = sync_msg.ControlMessage{ControlType: "join_slot", Data: slot}
	// 		joinTable = true
	// 	case 1:
	// 		msg.Type = sync_msg.CtrlMsgType
	// 		msg.Payload = sync_msg.ControlMessage{ControlType: "start_game", Data: ""}
	// 	case 2:
	// 		msg.Type = sync_msg.CtrlMsgType
	// 		msg.Payload = sync_msg.ControlMessage{ControlType: "next_game", Data: ""}
	// 	case 3:
	// 		fmt.Print("Enter player action (call, check, fold, raise, allin): ")
	// 		var action string
	// 		_, err := fmt.Scanln(&action)
	// 		if err != nil {
	// 			log.Println("Failed to read input:", err)
	// 			os.Exit(1)
	// 		}
	// 		switch action {
	// 		case "call", "check", "fold", "allin":
	// 			msg.Type = sync_msg.PlayerActMsgType
	// 			msg.Payload = sync_msg.PlayerMessage{ActionName: action, Value: 0}
	// 		case "raise":
	// 			fmt.Print("Enter raise amount: ")
	// 			var value int
	// 			_, err := fmt.Scanln(&value)
	// 			if err != nil {
	// 				log.Println("Failed to read input:", err)
	// 				os.Exit(1)
	// 			}
	// 			msg.Type = sync_msg.PlayerActMsgType
	// 			msg.Payload = sync_msg.PlayerMessage{ActionName: action, Value: value}
	// 		default:
	// 			fmt.Println("Invalid action. Please try again.")
	// 		}
	// 	case 4:
	// 		msg.Type = sync_msg.CtrlMsgType
	// 		msg.Payload = sync_msg.ControlMessage{ControlType: "request_buyin", Data: ""}
	// 	case 5:
	// 		fmt.Println("Exiting...")
	// 		return
	// 	default:
	// 		fmt.Println("Invalid choice. Please try again.")
	// 	}

	// 	err = ws.WriteJSON(msg)
	// 	if err != nil {
	// 		log.Fatalf("Failed to send message: %v", err)
	// 	}
	// 	fmt.Printf("Sent: %v\n", msg)
	// }
}
