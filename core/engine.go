package engine

import (
	"context"
	"fmt"
	msgpb "go-pk-server/gen"
	mylog "go-pk-server/log"
)

var DebugMode = true

type GaneInputType int

const (
	Unspecified GaneInputType = iota
	Refresh
	PlayerJoined
	PlayerLeft
	PlayerReady
	GameStarted
	GameEnded
	PlayerActed
	GameControl
)

func (at GaneInputType) String() string {
	return [...]string{"Unspecified", "Refresh", "PlayerJoined", "PlayerLeft",
		"PlayerReady", "GameStarted", "GameEnded", "PlayerActed", "GameControl"}[at]
}

type ControlIf interface {
	GetControlType() string
	GetOptions() []int32
}

type Input struct {
	// Common fields for all actions
	Type       GaneInputType
	PlayerInfo Player
	// Possible fields for an input
	PlayerAct  ActionIf
	ControlAct ControlIf
}

type NotifyGameStateReason int

const (
	NotifyGameStateReason_DONT_NOFITY     NotifyGameStateReason = 0
	NofityGameStateReason_ALL             NotifyGameStateReason = 1
	NotifyGameStateReason_NEW_ROUND       NotifyGameStateReason = 2
	NotifyGameStateReason_NEW_GAME        NotifyGameStateReason = 3
	NotifyGameStateReason_UPDATE_PLAYER   NotifyGameStateReason = 4
	NotifyGameStateReason_NEW_ACTION      NotifyGameStateReason = 5
	NofityGameStateReason_SETTING_CHANGED NotifyGameStateReason = 6
	NofityGameStateReason_SYNC_BALANCE    NotifyGameStateReason = 7
)

// String method for NotifyGameStateReason
func (r NotifyGameStateReason) String() string {
	return [...]string{"NotifyGameStateReason_DONT_NOFITY", "NofityGameStateReason_ALL",
		"NotifyGameStateReason_NEW_ROUND", "NotifyGameStateReason_NEW_GAME",
		"NotifyGameStateReason_UPDATE_PLAYER", "NotifyGameStateReason_NEW_ACTION",
		"NofityGameStateReason_SETTING_CHANGED", "NofityGameStateReason_SYNC_BALANCE"}[r]
}

type PublicRoom interface {
	BroadcastMessageToYourRoom(msg *msgpb.ServerMessage)
}

// GameEngine represents the game engine.
type GameEngineIf interface {
	StartEngine(bool) // Start the game engine with event driven mode (true) or synchronous mode (false)
	StopEngine()
	SetRoomAgent(room PublicRoom)
	PlayerJoin(player Player)
	PlayerLeave(player Player)
	ControlAction(ControlIf)
	PlayerAction(input ActionIf)
	ChangeSetting(setting *msgpb.GameSetting)
	GetGameSetting() *msgpb.GameSetting
	GetGameState() *msgpb.GameState
}

type EngineState int

const (
	EngineState_INITIALIZING     EngineState = 0
	EngineState_WAIT_FOR_PLAYING EngineState = 1
	EngineState_PLAYING          EngineState = 2
	EngineState_PAUSED           EngineState = 3
)

// Overwrite string method for EngineState
func (e EngineState) String() string {
	return [...]string{"EngineState_INITIALIZING", "EngineState_WAIT_FOR_PLAYING", "EngineState_PLAYING", "EngineState_PAUSED"}[e]
}

type GameEngine struct {
	gameSessionID int
	balanceMgr    *BalanceManager
	playerMgr     *TableManager
	game          *Game
	room          PublicRoom

	ntfReason NotifyGameStateReason

	eState      EngineState
	eventDriven bool

	ctx           context.Context
	cancel        context.CancelFunc
	GameInputChan chan Input
}

var MyGame = NewGameEngine()

// NewGameEngine creates a new instance of the game engine.
func NewGameEngine() GameEngineIf {
	// Add your initialization code here
	return &GameEngine{
		gameSessionID: 1,
		eState:        EngineState_INITIALIZING,
		eventDriven:   false,
		balanceMgr:    NewBalanceManager(),
		playerMgr:     NewTableManager(),
		GameInputChan: make(chan Input, 10), // Change to buffered channel with capacity 10

	}
}

// StartEngine starts the game.
func (g *GameEngine) StartEngine(e bool) {
	g.eventDriven = e
	// Run the game engine with go routine
	if g.eventDriven {
		g.ctx, g.cancel = context.WithCancel(context.Background())
		go g.EngineLoop(g.ctx)
	}
	act := Input{Type: GaneInputType(EngineState_WAIT_FOR_PLAYING)}
	g.processActions(act)
}

func (g *GameEngine) StopEngine() {
	if g.eventDriven {
		g.cancel()
	}
}

func (g *GameEngine) SetRoomAgent(room PublicRoom) {
	g.room = room
}

func (g *GameEngine) PlayerJoin(player Player) {
	input := Input{Type: PlayerJoined, PlayerInfo: player}
	g.processActions(input)
}

func (g *GameEngine) PlayerLeave(player Player) {
	input := Input{Type: PlayerLeft, PlayerInfo: player}
	g.processActions(input)
}

func (g *GameEngine) ControlAction(ctrl ControlIf) {
	input := Input{Type: GameControl, ControlAct: ctrl}
	g.processActions(input)
}

// PerformAction performs the specified input for the given player.
func (g *GameEngine) PlayerAction(input ActionIf) {
	// Send input to game engine
	act := Input{Type: PlayerActed, PlayerAct: input}
	g.processActions(act)
}

func (g *GameEngine) processActions(input Input) {
	if g.eventDriven {
		g.GameInputChan <- input
	} else {
		g.RunGameEngine(input)
	}
}

func (g *GameEngine) gotoState(newState EngineState) {
	g.eState = newState
}

// EngineLoop runs the game engine in a loop.
func (g *GameEngine) EngineLoop(ctx context.Context) {
	for {
		select {
		case input := <-g.GameInputChan:
			mylog.Infof("Run game engine with current %s and input: %v\n", g.eState, input.Type)
			g.RunGameEngine(input)
		case <-ctx.Done():
			// Game ended
			return
		}
	}
}

func (g *GameEngine) RunGameEngine(input Input) {
	switch g.eState {
	case EngineState_INITIALIZING:
		// Room is created log
		g.Initializing()
		g.eState = EngineState_WAIT_FOR_PLAYING
	case EngineState_WAIT_FOR_PLAYING:
		// Wait for players to join, buy chip, and ready up
		switch input.Type {
		case PlayerJoined:
			g.handleJoiningPlayer(input.PlayerInfo)
		case PlayerLeft:
			g.handleLeavingPlayer(input.PlayerInfo)
		case PlayerReady, GameStarted:
			// Play the game
			if g.game.Play() {
				g.eState = EngineState_PLAYING
				g.needNtfAndReason(NotifyGameStateReason_NEW_GAME)
			}
		case GameControl:
			// Handle control message
			g.handleControlMessage(input)
		}
	case EngineState_PLAYING:
		switch input.Type {
		case PlayerActed:
			g.game.HandleActions(input.PlayerAct)
			g.needNtfAndReason(NotifyGameStateReason_NEW_ACTION)
		case PlayerLeft:
			mylog.Warnf("Player %s left during the game", input.PlayerInfo.Name())
			g.handleLeavingPlayer(input.PlayerInfo)
		case GameControl:
			g.handleControlMessage(input)
			// panic should not receive control message in playing state
			mylog.Error("Game engine should not receive control message in playing state")
		}
	case EngineState_PAUSED:
		switch input.Type {
		case GameStarted:
			g.eState = EngineState_PLAYING
		}
	}

	// Try to notify the game state each time the game engine is run
	g.NotifyGameState()
}

func (g *GameEngine) needNtfAndReason(reason NotifyGameStateReason) {
	mylog.Infof("Need to notify the game state for %v\n", reason)
	g.ntfReason = reason
}

func (g *GameEngine) NotifyGameState() {
	// Additionally, notify the other state if needed
	switch g.ntfReason {
	case NotifyGameStateReason_DONT_NOFITY:
		return
	case NofityGameStateReason_ALL,
		NotifyGameStateReason_NEW_ROUND,
		NotifyGameStateReason_UPDATE_PLAYER,
		NotifyGameStateReason_NEW_ACTION:
		// Special case: Notifying
	case NofityGameStateReason_SETTING_CHANGED:
	// Notify to all players if the setting is changed
	case NofityGameStateReason_SYNC_BALANCE:
		if g.room != nil {
			g.room.BroadcastMessageToYourRoom(&msgpb.ServerMessage{
				Message: &msgpb.ServerMessage_BalanceInfo{
					BalanceInfo: g.balanceMgr.GetBalanceSummary(),
				},
			})
		}
	case NotifyGameStateReason_NEW_GAME:
		// Notify directly to all players if they have new cards
		for _, player := range g.playerMgr.players {
			if player != nil && player.HasNewCards() {
				player.NotifyPlayerIfNewHand()
			}
		}
	}

	// Always notify the game state
	if g.room != nil {
		g.room.BroadcastMessageToYourRoom(&msgpb.ServerMessage{
			Message: &msgpb.ServerMessage_GameState{
				GameState: g.GetGameState(),
			},
		})
	} else {
		mylog.Errorf("Room agent is not set\n")
	}
	g.ntfReason = NotifyGameStateReason_DONT_NOFITY
}

// Initializing handles the EngineState_INITIALIZING state.
func (g *GameEngine) Initializing() {
	g.game = NewGame(
		&msgpb.GameSetting{
			MaxPlayers:   6,
			MinPlayers:   2,
			SmallBlind:   10,
			MinStackSize: 500,
			BigBlind:     20,
			TimePerTurn:  0, // 0 means no limit
		},
		g.playerMgr,
		NewDeck(),
		g.gotoState,
	)
}

// HandleWaitForPlayers handles the EngineState_WAIT_FOR_PLAYING state.
func (g *GameEngine) handleJoiningPlayer(player Player) {
	if player != nil {
		fmt.Printf("Player %s joined the game\n", player.Name())
		g.playerMgr.AddPlayer(player.Position(), player)
		// If first player joined, set the game button position
		if g.playerMgr.GetNumberOfPlayers() == 1 {
			g.game.gs.ButtonPosition = player.Position()
		}
	}
	// Notify the game state due to new player
	g.needNtfAndReason(NotifyGameStateReason_UPDATE_PLAYER)
}

func (g *GameEngine) handleLeavingPlayer(player Player) {
	if player != nil {
		// Notify the game state due to player left
		g.needNtfAndReason(NotifyGameStateReason_UPDATE_PLAYER)
		// Take remaining chips from the player
		remaining := player.Chips()
		if remaining > 0 {
			g.balanceMgr.ReturnStack(player.Name(), remaining)
			g.needNtfAndReason(NofityGameStateReason_SYNC_BALANCE)
		}
		// Remove and If no player left, reset the game
		if g.playerMgr.RemovePlayer(player.Position()) < 2 {
			g.eState = EngineState_WAIT_FOR_PLAYING
			g.game.ResetGame()
		}
	}
}

func (g *GameEngine) handleControlMessage(intput Input) {
	g.needNtfAndReason(NotifyGameStateReason_NEW_ACTION)
	// Perform control action
	ctrlActionName := intput.ControlAct.GetControlType()
	switch ctrlActionName {
	case "pause_game":
		mylog.Info("A player has paused the GAME")
		g.eState = EngineState_PAUSED
	case "resume_game":
		mylog.Info("A player has resumed the GAME")
		g.eState = EngineState_PLAYING
	case "ready_game":
		mylog.Info("A player has readied up")
		input := Input{Type: PlayerReady}
		g.processActions(input)
	case "start_game":
		// Log the game start
		mylog.Info("A player has started the GAME")
		// send input to start the game
		act := Input{Type: GaneInputType(GameStarted)}
		g.processActions(act)
	case "leave_game":
		mylog.Info("A player has request to leave the GAME")
		opts := intput.ControlAct.GetOptions()
		if len(opts) > 0 {
			// Find the player
			reqPlayerIdx := int(opts[0])
			player := g.playerMgr.GetPlayer(reqPlayerIdx)
			if player != nil {
				input := Input{Type: PlayerLeft, PlayerInfo: player}
				g.processActions(input)
			} else {
				mylog.Errorf("Requesting player %d for leaving not found", reqPlayerIdx)
			}
		} else {
			mylog.Errorf("Did not provide player id to leave the game")
		}
	case "request_buyin":
		mylog.Info("A player sending request to ADD 1 buyin")
		// parse optional amount
		opts := intput.ControlAct.GetOptions()
		if len(opts) > 0 {
			reqPlayerIdx := int(opts[0])
			// Find the player
			player := g.playerMgr.GetPlayer(reqPlayerIdx)
			if player == nil {
				mylog.Errorf("Requesting add chips to player %d not found ", reqPlayerIdx)
			} else {
				player.AddChips(g.balanceMgr.TakeOneBuyIn(player.Name()))
				g.needNtfAndReason(NofityGameStateReason_SYNC_BALANCE)
			}
		} else {
			mylog.Errorf("Did not provide player id to add chips")
		}
	case "payback_buyin":
		mylog.Info("A player sending request to PAYBACK 1 buyin")
		// parse optional amount
		opts := intput.ControlAct.GetOptions()
		if len(opts) > 0 {
			reqPlayerIdx := int(opts[0])
			// Find the player
			player := g.playerMgr.GetPlayer(reqPlayerIdx)
			if player == nil {
				mylog.Errorf("Requesting payback chips from player %d not found ", reqPlayerIdx)
			} else {
				// If player has more than 1 buyin, payback 1 buyin
				if player.Chips()-BUY_IN_SIZE > BUY_IN_SIZE-int(g.game.setting.MinStackSize) {
					player.TakeChips(BUY_IN_SIZE)
					g.balanceMgr.PaybackOneBuyIn(player.Name())
					g.needNtfAndReason(NofityGameStateReason_SYNC_BALANCE)
				}
			}
		} else {
			mylog.Errorf("Did not provide player id to add chips")
		}
	case "sync_balance":
	case "sync_game_state":
		mylog.Info("A player has requested to sync GAME state")
		// send refresh input to sync the game state
		act := Input{Type: GaneInputType(Refresh)}
		g.processActions(act)
	default:
		mylog.Errorf("Game engine not support control message type: %v", ctrlActionName)
		g.needNtfAndReason(NotifyGameStateReason_DONT_NOFITY)
		return
	}

}

// ChangeSetting changes the game setting.
func (g *GameEngine) ChangeSetting(setting *msgpb.GameSetting) {
	// Validate the setting
	if setting.MaxPlayers < 2 || setting.MaxPlayers > 6 {
		fmt.Println("Invalid setting: MaxPlayers")
		return
	}
}

// GetGameSetting returns the game setting.
func (g *GameEngine) GetGameSetting() *msgpb.GameSetting {
	if g.game != nil && g.game.setting != nil {
		return g.game.setting
	}
	return nil
}

// GetGameState synchronizes the game state.
func (g *GameEngine) GetGameState() *msgpb.GameState {
	// Create a message to sync the game state

	// fakeGameState := &msgpb.GameState{
	// 	Players:        make([]*msgpb.PlayerState, 0),
	// 	PotSize:        1000,
	// 	DealerId:       0,
	// 	CommunityCards: make([]*msgpb.Card, 0),
	// 	CurrentBet:     0,
	// 	CurrentRound:   msgpb.RoundStateType_PREFLOP,
	// 	FinalResult: &msgpb.Result{
	// 		WinnerPosition: 3,
	// 		WonChip:        1000,
	// 		ShowingCards: []*msgpb.PeerState{
	// 			{
	// 				TablePos: 4,
	// 				PlayerCards: []*msgpb.Card{
	// 					{Suit: msgpb.SuitType_SPADES, Rank: msgpb.RankType_ACE},
	// 					{Suit: msgpb.SuitType_DIAMONDS, Rank: msgpb.RankType_KING},
	// 				},
	// 			},
	// 			{
	// 				TablePos: 4,
	// 				PlayerCards: []*msgpb.Card{
	// 					{Suit: msgpb.SuitType_SPADES, Rank: msgpb.RankType_ACE},
	// 					{Suit: msgpb.SuitType_DIAMONDS, Rank: msgpb.RankType_KING},
	// 				},
	// 			},
	// 		},
	// 	},
	// }
	// fakeGameState.CommunityCards = []*msgpb.Card{
	// 	{Suit: msgpb.SuitType_SPADES, Rank: msgpb.RankType_ACE},
	// 	{Suit: msgpb.SuitType_DIAMONDS, Rank: msgpb.RankType_KING},
	// 	{Suit: msgpb.SuitType_CLUBS, Rank: msgpb.RankType_QUEEN},
	// 	{Suit: msgpb.SuitType_HEARTS, Rank: msgpb.RankType_JACK},
	// 	{Suit: msgpb.SuitType_SPADES, Rank: msgpb.RankType_TEN},
	// }
	// fakeGameState.Players = []*msgpb.PlayerState{
	// 	{
	// 		Name:          "player1",
	// 		Chips:         1500,
	// 		TablePosition: 0,
	// 		Status:        "Playing",
	// 	},
	// 	{
	// 		Name:          "player2",
	// 		Chips:         2000,
	// 		TablePosition: 1,
	// 		Status:        "Wait4Act",
	// 	},
	// 	{
	// 		Name:          "player3",
	// 		Chips:         4000,
	// 		TablePosition: 2,
	// 		Status:        "Fold",
	// 	},
	// 	{
	// 		Name:          "player4",
	// 		Chips:         3000,
	// 		TablePosition: 4,
	// 		Status:        "Check",
	// 	},
	// 	{
	// 		Name:          "player5",
	// 		Chips:         4800,
	// 		TablePosition: 5,
	// 		Status:        "Raise",
	// 		CurrentBet:    200,
	// 	},
	// }

	// return fakeGameState

	syncMsg := &msgpb.GameState{
		Players:        make([]*msgpb.PlayerState, 0),
		PotSize:        int32(g.game.gs.pot.Size()),
		DealerId:       int32(g.game.gs.ButtonPosition),
		CommunityCards: make([]*msgpb.Card, 0),
		CurrentBet:     int32(g.game.gs.CurrentBet),
		CurrentRound:   g.game.gs.CurrentRound,
	}
	// Add the community cards
	for _, card := range g.game.gs.cc.GetCards() {
		syncMsg.CommunityCards = append(syncMsg.CommunityCards, &msgpb.Card{Suit: msgpb.SuitType(card.Suit), Rank: msgpb.RankType(card.Value)})
	}
	// Add the players
	for _, player := range g.playerMgr.players {
		if player != nil {
			syncMsg.Players = append(syncMsg.Players, &msgpb.PlayerState{
				Name:          player.Name(),
				TablePosition: int32(player.Position()),
				Chips:         int32(player.Chips()),
				IsDealer:      player.Position() == g.game.gs.ButtonPosition,
				Status:        player.Status(),
				CurrentBet:    int32(player.CurrentBet()),
				ChangeAmount:  int32(player.ChipChange()),
				NoActions:     player.UnsupportActs(),
			})
		}
	}

	if g.game.gs.FinalResult != nil {
		syncMsg.FinalResult = &msgpb.Result{
			ShowingCards: make([]*msgpb.PeerState, 0),
		}
		// Add the showing cards
		for _, peer := range g.game.gs.FinalResult.ShowingCards {
			peerState := &msgpb.PeerState{
				TablePos:    int32(peer.TablePos),
				PlayerCards: make([]*msgpb.Card, 0),
			}
			for _, card := range peer.PlayerCards {
				peerState.PlayerCards = append(peerState.PlayerCards, &msgpb.Card{Suit: msgpb.SuitType(card.Suit), Rank: msgpb.RankType(card.Rank)})
			}
			syncMsg.FinalResult.ShowingCards = append(syncMsg.FinalResult.ShowingCards, peerState)
		}
	}

	return syncMsg
}
