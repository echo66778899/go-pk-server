package main

import (
	"fmt"
	"go-pk-server/client/handler"
	"go-pk-server/client/network"
	"go-pk-server/client/ui"
	msgpb "go-pk-server/gen"
	"log"
	"os"
	"time"
)

var (
	SKIP_LOGIN_ROOM    = false
	AUTO_ENTER_REQUEST = false
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
	l          *ui.List
	rankDisp   *ui.RankingText
)

func initClient() {
	// Table
	table = ui.NewTable()
	table.SetRect(ui.TABLE_CENTER_X-ui.TABLE_RADIUS_X, ui.TABLE_CENTER_Y-ui.TABLE_RADIUS_Y, 2*ui.TABLE_RADIUS_X, 2*ui.TABLE_RADIUS_Y)

	// New board of cards
	board = ui.NewCards()
	board.SetTitle("Community Cards")
	board.SetCoodinate(ui.COMMUNITY_CARDS_X, ui.COMMUNITY_CARDS_Y)

	// New central text box
	centerText = ui.NewTextBox()
	centerText.SetRect(ui.TABLE_CENTER_X-12, ui.TABLE_CENTER_Y+2, ui.TABLE_CENTER_X-12+24, ui.TABLE_CENTER_Y+5)

	// Dealer Icon
	dealerIcon = ui.NewDealerWidget()
	dealerIcon.SetRect(ui.TABLE_CENTER_X-10, ui.TABLE_CENTER_Y+5, ui.TABLE_CENTER_X-10+21, ui.TABLE_CENTER_Y+8)

	// Players slider
	playerWg = ui.NewPlayersGroup()

	// Buttons
	btnCtrl = ui.NewButtonCtrlCenter()
	btnCtrl.InitButtonPosition()

	// List balance info
	l = ui.NewList()
	// l.Rows = []string{
	// 	"[0] github.com/gizak/termui/v3",
	// 	"[1] [你好，世界](fg:blue)",
	// 	"[2] [こんにちは世界](fg:red)",
	// 	"[3] [color](fg:white,bg:green) output",
	// 	"[4] output.go",
	// 	"[5] random_out.go",
	// 	"[6] dashboard.go",
	// 	"[7] foo",
	// 	"[8] bar",
	// 	"[9] baz",
	// }33
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetRect(ui.BALANCE_INFO_X, ui.BALANCE_INFO_Y, ui.BALANCE_INFO_X+29, 12)

	// Ranking display
	rankDisp = ui.NewRankingText()
}

func render() {
	uiItems := []ui.Drawable{table, board, centerText, dealerIcon, l}
	uiItems = append(uiItems, playerWg.GetAllItems()...)
	uiItems = append(uiItems, btnCtrl.GetDisplayingButton()...)
	uiItems = append(uiItems, rankDisp.GetDisplayingTexts()...)
	ui.Render(uiItems...)
}

// Just for testing
func autoChoseSlotAndBuyIn(selection string) {

	profiles := map[string]int{
		"1": 2,
		"2": 4,
		"3": 0,
	}

	if AUTO_ENTER_REQUEST {
		join := msgpb.ClientMessage{
			Message: &msgpb.ClientMessage_JoinGame{
				JoinGame: &msgpb.JoinGame{
					ChooseSlot: int32(profiles[selection]),
				},
			},
		}
		agent.SendingMessage(&join)
		time.Sleep(100 * time.Millisecond)
		agent.SendingMessage(factoryCtrlMessage("request_buyin",
			int32(profiles[selection])))
	}
}

func getRoomInfoInput(selection string) (playerName, room, passcode, sessId string) {
	if SKIP_LOGIN_ROOM {
		profiles := map[string]struct{ name, room, passcode, sessId string }{
			"1": {"Hai Phan", "2", "1", "0"},
			"2": {"Quynh Mi", "2", "1", "0"},
			"3": {"Thalaba", "2", "1", "0"},
		}
		if _, ok := profiles[selection]; !ok {
			log.Println("Profile not found")
			os.Exit(128)
		}
		return profiles[selection].name, profiles[selection].room,
			profiles[selection].passcode, profiles[selection].sessId
	}

	// Enter the authentication details
	fmt.Println("+-------------------------------------------+")
	fmt.Println("|  Welcome to the Poker Game Client v1.0.0  |")
	fmt.Println("+-------------------------------------------+")
	fmt.Print("-> Enter your name (no space): ")

	_, err := fmt.Scanln(&playerName)
	if err != nil {
		log.Println("Failed to read input:", err)
		os.Exit(1)
	}

	fmt.Print("-> Enter room number: ")
	_, err = fmt.Scanln(&room)

	if err != nil {
		log.Println("Failed to read input:", err)
		os.Exit(1)
	}

	fmt.Print("-> Enter room pass: ")
	_, err = fmt.Scanln(&passcode)

	if err != nil {
		log.Println("Failed to read input:", err)
		os.Exit(1)
	}

	// fmt.Print("-> Enter your session : ")
	// _, err = fmt.Scanln(&sessId)

	// if err != nil {
	// 	log.Println("Failed to read input:", err)
	// 	os.Exit(1)
	// }

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

func factoryCtrlMessage(ctrl string, opts ...int32) *msgpb.ClientMessage {
	return &msgpb.ClientMessage{
		Message: &msgpb.ClientMessage_ControlAction{
			ControlAction: &msgpb.ControlAction{
				ControlType: ctrl,
				Options:     opts,
			},
		},
	}
}

func buttonEnterEvtHandler(opt ...int) {
	log.Printf("Button Enter Event Handler: %v", opt)

	bt := ui.UIButtonType(opt[0])
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
	case ui.BNT_PauseGameButton:
		//agent.SendingMessage(factoryCtrlMessage("pause_game"))
	case ui.BNT_StartGameButton:
		agent.SendingMessage(factoryCtrlMessage("start_game"))
	case ui.BNT_LeaveGameButton:
		agent.SendingMessage(factoryCtrlMessage("leave_game",
			int32(ui.UI_MODEL_DATA.YourTablePosition)))
	case ui.BNT_RequestBuyinButton:
		agent.SendingMessage(factoryCtrlMessage("request_buyin",
			int32(ui.UI_MODEL_DATA.YourTablePosition)))
	case ui.BNT_PaybackBuyinButton:
		agent.SendingMessage(factoryCtrlMessage("payback_buyin",
			int32(ui.UI_MODEL_DATA.YourTablePosition)))
	case ui.BNT_JoinSlotButton:
		if len(opt) > 1 {
			join := msgpb.ClientMessage{
				Message: &msgpb.ClientMessage_JoinGame{
					JoinGame: &msgpb.JoinGame{
						ChooseSlot: int32(opt[1]),
					},
				},
			}
			agent.SendingMessage(&join)
		}
	default:
	}
}

// Main can be run with the following command:
// go run main.go 1
// go run main.go 2
func main() {
	// Get the profile number from the command line
	args := os.Args

	// The first argument (args[0]) is the program's name
	log.Println("Program Name:", args[0])
	dest, profile := "", ""
	// Check if there are additional arguments
	if len(args) > 1 {
		log.Println("Arguments passed:")
		for i, arg := range args[1:] {
			// Convert the argument to an integer
			log.Printf("Arg %d: %s\n", i+1, arg)
		}
		// convert the argument to an integer
		dest = args[1]
	} else {
		log.Println("No arguments were passed.")
		SKIP_LOGIN_ROOM = false
	}

	if SKIP_LOGIN_ROOM {
		profile = dest
		dest = "localhost:8088"
	}

	// Connect to the server
	if agent.Connect(dest) {
		// Send a message to the server
		n, r, p, s := getRoomInfoInput(profile)
		ui.UI_MODEL_DATA.YourLoginUsernameID = n

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

	// Set up the log file
	logFile, err := os.OpenFile(ui.UI_MODEL_DATA.YourLoginUsernameID+"client.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()
	// Redirect log output to the file
	log.SetOutput(logFile)

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

	// if testing auto chose slot and request buyin
	autoChoseSlotAndBuyIn(profile)

	// Game looop
	for {
		// Log
		log.Printf("WAITING FOR NEW EVENT...\n")
		// Process input and server message
		select {
		case ev := <-keyboardEventsChan:
			switch ev.EventType {
			case ui.LEFT:
				btnCtrl.MoveLeft()
			case ui.RIGHT:
				btnCtrl.MoveRight()
			case ui.UP:
				btnCtrl.MoveUp()
			case ui.DOWN:
				btnCtrl.MoveDown()
			case ui.ENTER:
				btnCtrl.Enter()
			case ui.START_GAME:
				agent.SendingMessage(factoryCtrlMessage("start_game"))
			case ui.PAUSE_GAME:
				agent.SendingMessage(factoryCtrlMessage("pause_game"))
			case ui.LEAVE_GAME:
				agent.SendingMessage(factoryCtrlMessage("leave_game",
					int32(ui.UI_MODEL_DATA.YourTablePosition)))
			case ui.TAKE_BUYIN:
				agent.SendingMessage(factoryCtrlMessage("request_buyin",
					int32(ui.UI_MODEL_DATA.YourTablePosition)))
			case ui.GIVE_BUYIN:
				agent.SendingMessage(factoryCtrlMessage("payback_buyin",
					int32(ui.UI_MODEL_DATA.YourTablePosition)))
			case ui.FOLD:
				agent.SendingMessage(factoryAction("fold"))
			case ui.CHECK_OR_CALL:
				action := "check"
				if ui.UI_MODEL_DATA.YourPlayerState != nil {
					for _, a := range ui.UI_MODEL_DATA.YourPlayerState.NoActions {
						if a == msgpb.PlayerGameActionType_CHECK {
							action = "call"
						}
					}
				}
				agent.SendingMessage(factoryAction(action))
			case ui.RAISE:
				// Raise min to double the current bet
				agent.SendingMessage(factoryAction("raise", ui.UI_MODEL_DATA.CurrentBet))
			case ui.ALL_IN:
				agent.SendingMessage(factoryAction("allin"))

			// Menu buttons for testing
			case ui.SPACE:
				agent.SendingMessage(factoryCtrlMessage("sync_game_state"))
			case ui.BACKSPACE:
				agent.SendingMessage(factoryCtrlMessage("request_game_end"))
			case ui.MENU1:
				btnCtrl.SetMenu(ui.ButtonMenuType_PLAYING_BTN)
				btnCtrl.EnableButtonCtrl(true)
			case ui.MENU2:
				btnCtrl.SetMenu(ui.ButtonMenuType_CTRL_BTN)
				btnCtrl.EnableButtonCtrl(true)
			case ui.MENU3:
				btnCtrl.SetMenu(ui.ButtonMenuType_SLOTS_BTN)
				btnCtrl.EnableButtonCtrl(true)
			case ui.EXIT:
				return
			default:
				fmt.Println("Received unsupport key:", ev.Key)
				continue
			}
		case m := <-agent.ReceivingMessage():
			h.HandleServerMessage(m)
		}

		// Update UI element state based on the model
		btnCtrl.SetMenu(ui.UI_MODEL_DATA.ActiveButtonMenu)
		btnCtrl.EnableButtonCtrl(ui.UI_MODEL_DATA.IsButtonEnabled)
		btnCtrl.UpdateState()
		btnCtrl.DisableListButton(ui.UI_MODEL_DATA.YourPlayerState)

		playerWg.UpdateState(true)
		playerWg.UpdateGroupPlayers(ui.UI_MODEL_DATA.MaxPlayers)
		playerWg.UpdatePocketPair(ui.UI_MODEL_DATA.YourPrivateState)

		// If current round is not SHOWDOWN, update all player's pocket pair
		if ui.UI_MODEL_DATA.CurrentRound == msgpb.RoundStateType_SHOWDOWN &&
			ui.UI_MODEL_DATA.Result != nil {

			// Update the ranking display
			rankDisp.UpdateTextsBasedPlayers()

			for _, p := range ui.UI_MODEL_DATA.Result.ShowingCards {
				if p != nil {
					playerWg.UpdatePocketPairAtPosition(int(p.TablePos), p)
					rankDisp.UpdateTextAtPosition(int(p.TablePos), p)
				}
			}
		} else if ui.UI_MODEL_DATA.CurrentRound == msgpb.RoundStateType_INITIAL ||
			ui.UI_MODEL_DATA.CurrentRound == msgpb.RoundStateType_PREFLOP {
			// Clear all the pocket pair
			rankDisp.ClearAllTexts()
		}

		board.SetCards(ui.UI_MODEL_DATA.CommunityCards)
		centerText.Text = fmt.Sprintf("%s [POT:%d]\n", ui.UI_MODEL_DATA.CurrentRound.String(), ui.UI_MODEL_DATA.Pot)

		dealerIcon.IndexUI(ui.UI_MODEL_DATA.DealerPosition, ui.UI_MODEL_DATA.MaxPlayers)
		dealerIcon.SetVisible(ui.UI_MODEL_DATA.IsDealerVisible)

		l.Rows = ui.UI_MODEL_DATA.PlayersBalance
		//testGame()

		// Render UI elements
		render()
	}
}
