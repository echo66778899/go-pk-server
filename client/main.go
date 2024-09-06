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
	TEST     = false
	// pointsChan         = make(chan int)
	keyboardEventsChan = make(chan ui.KeyboardEvent)
)

// Connect to the WebSocket server
var agent = network.NewAgent()

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

	// ui.UI_MODEL_DATA.Players = []*msgpb.Player{
	// 	// {
	// 	// 	Name:          "player1",
	// 	// 	Chips:         1500,
	// 	// 	TablePosition: 0,
	// 	// },
	// 	{
	// 		Name:          "player2",
	// 		Chips:         2000,
	// 		TablePosition: 1,
	// 	},
	// 	{
	// 		Name:          "player3",
	// 		Chips:         4000,
	// 		TablePosition: 2,
	// 	},
	// 	{
	// 		Name:          "Hai",
	// 		Chips:         3000,
	// 		TablePosition: 4,
	// 	},
	// }
	playerWg.UpdateState(true)
	if i%2 == 0 {
		playerWg.UpdatePocketPair(nil)
	} else {
		playerWg.UpdatePocketPair(&msgpb.PeerState{
			PlayerCards: []*msgpb.Card{
				{Suit: msgpb.SuitType_HEARTS, Rank: msgpb.RankType_DEUCE},
				{Suit: msgpb.SuitType_SPADES, Rank: msgpb.RankType_SEVEN},
			},
		})
	}
	btnCtrl.UpdateState()
}

func initClient() {
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

func factoryAction(action string, amount ...int) *msgpb.ClientMessage {
	a := 0
	if len(amount) > 0 {
		a = amount[0]
	}

	return &msgpb.ClientMessage{
		Message: &msgpb.ClientMessage_PlayerAction{
			PlayerAction: &msgpb.PlayerAction{
				ActionType:  action,
				RaiseAmount: int32(a),
			},
		},
	}
}

func factoryCtrlMessage(ctrl string) *msgpb.ClientMessage {
	return &msgpb.ClientMessage{
		Message: &msgpb.ClientMessage_ControlMessage{
			ControlMessage: ctrl,
		},
	}
}

func buttonEnterEvtHandler(opt ...int) {
	log.Printf("Button Enter Event Handler: %v", opt)

	bt := ui.ButtonType(opt[0])
	switch bt {
	case ui.BNT_FoldButton:
		agent.SendingMessage(factoryAction("fold"))
	case ui.BNT_CheckButton:
		agent.SendingMessage(factoryAction("check"))
	case ui.BNT_CallButton:
		agent.SendingMessage(factoryAction("call"))
	case ui.BNT_RaiseButton:
		if len(opt) > 1 {
			agent.SendingMessage(factoryAction("raise", opt[1]))
		}
	case ui.BNT_AllInButton:
		agent.SendingMessage(factoryAction("allin"))
	case ui.BNT_JoinTableButton:
	case ui.BNT_StartGameButton:
		agent.SendingMessage(factoryCtrlMessage("start_game"))
	case ui.BNT_LeaveGameButton:
	case ui.BNT_RequestChipButton:
		agent.SendingMessage(factoryCtrlMessage("request_buyin"))
	case ui.BNT_SlotButton:
		if len(opt) > 1 {
			join := msgpb.ClientMessage{
				Message: &msgpb.ClientMessage_JoinGame{
					JoinGame: &msgpb.JoinGame{
						ChooseSlot: int32(opt[1]),
					},
				},
			}
			agent.SendingMessage(&join)
			ui.UI_MODEL_DATA.YourTablePosition = opt[1]
		}
	default:
	}
}

func main() {
	logFile, err := os.OpenFile("client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Redirect log output to the file
	log.SetOutput(logFile)

	if agent.Connect() {
		// Send a message to the server
		n, r, p, s := getRoomInfoInput()
		ui.UI_MODEL_DATA.YourUsernameID = n

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
		agent.SendingMessage(&msg)
	}
	defer agent.Close()

	// Initialize the UI
	ui.Init()
	defer ui.Deinit()

	// Initialize the windows and widgets
	initClient()

	// set the button handler
	btnCtrl.SetUserButtonInteractionHandler(buttonEnterEvtHandler)

	// new chan for keyboard
	go ui.ListenToKeyboard(keyboardEventsChan)

	// handler for server message
	h := handler.Handler{UI_Model: &ui.UI_MODEL_DATA}

	// Game looop
	for {
		// Process input and server message
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
				agent.SendingMessage(factoryCtrlMessage("sync_game_state"))
			case ui.BACKSPACE:
				if btnCtrl.CtrlEnabled {
					btnCtrl.EnableButtonCtrl(false)
				} else {
					btnCtrl.EnableButtonCtrl(true)
				}
			case ui.MENU1:
				btnCtrl.SetMenu(ui.ButtonMenuType_PLAYING_BTN)
				btnCtrl.EnableButtonCtrl(true)
			case ui.MENU2:
				btnCtrl.SetMenu(ui.ButtonMenuType_CTRL_BTN)
				btnCtrl.EnableButtonCtrl(true)
			case ui.MENU3:
				btnCtrl.SetMenu(ui.ButtonMenuType_SLOTS_BTN)
				btnCtrl.EnableButtonCtrl(true)
			case ui.END:
				return
			}
		case m := <-agent.ReceivingMessage():
			h.HandleServerMessage(m)
		}

		// Check player state to enable/disable button control
		for _, p := range ui.UI_MODEL_DATA.Players {
			if p.TablePosition == int32(ui.UI_MODEL_DATA.YourTablePosition) {
				if p.Status == "Wait4Act" {
					btnCtrl.EnableButtonCtrl(true)
				}
				break
			}
		}

		// Update UI state based on the model
		playerWg.UpdateState(true)
		playerWg.UpdateGroupPlayers(ui.UI_MODEL_DATA.MaxPlayers)
		playerWg.UpdatePocketPair(ui.UI_MODEL_DATA.YourPrivateState)
		board.SetCards(ui.UI_MODEL_DATA.CommunityCards)
		btnCtrl.UpdateState()
		dealerIcon.IndexUI(ui.UI_MODEL_DATA.DealerPosition, ui.UI_MODEL_DATA.MaxPlayers)

		//testGame()

		// Render UI elements
		render()
	}
}
