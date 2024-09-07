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
	PlayerJoined
	PlayerLeft
	PlayerReady
	GameStarted
	GameEnded
	PlayerActed
	NextGame
)

func (at GaneInputType) String() string {
	return [...]string{"Unspecified", "PlayerJoined", "PlayerLeft",
		"PlayerReady", "GameStarted", "GameEnded", "PlayerActed", "NextGame"}[at]
}

type Input struct {
	// Common fields for all actions
	Type       GaneInputType `json:"type"`
	PlayerInfo Player        `json:"player_info"`
	// Possible fields for an input
	PlayerAct ActionIf
}

type NotifyGameStateReason int

const (
	NotifyGameStateReason_DONT_NOFITY   NotifyGameStateReason = 0
	NofityGameStateReason_ALL           NotifyGameStateReason = 1
	NotifyGameStateReason_NEW_ROUND     NotifyGameStateReason = 2
	NotifyGameStateReason_NEW_GAME      NotifyGameStateReason = 3
	NotifyGameStateReason_UPDATE_PLAYER NotifyGameStateReason = 4
	NotifyGameStateReason_NEW_ACTION    NotifyGameStateReason = 5
)

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
	StartGame()
	NextGame()
	Ready()
	PlayerAction(input ActionIf)
	ChangeSetting(setting *msgpb.GameSetting)
	GetGameSetting() *msgpb.GameSetting
	SyncGameState() *msgpb.GameState
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

func (g *GameEngine) StartGame() {
	// Log the game start
	fmt.Println("Game started")
	// send input to start the game
	act := Input{Type: GaneInputType(GameStarted)}
	g.processActions(act)
}

func (g *GameEngine) NextGame() {
	// Log the next game
	fmt.Println("Next game")
	// send input to start the next game
	act := Input{Type: NextGame}
	g.processActions(act)
}

// PerformAction performs the specified input for the given player.
func (g *GameEngine) PlayerAction(input ActionIf) {
	// Send input to game engine
	act := Input{Type: PlayerActed, PlayerAct: input}
	g.processActions(act)
}

func (g *GameEngine) Ready() {
	input := Input{Type: PlayerReady}
	g.processActions(input)
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
			g.HandleJoiningPlayer(input.PlayerInfo)
			g.NeedNtfAndReason(NotifyGameStateReason_UPDATE_PLAYER)
		case PlayerLeft:
			g.HandleLeavingPlayer(input.PlayerInfo)
			g.NeedNtfAndReason(NotifyGameStateReason_UPDATE_PLAYER)
		case PlayerReady, GameStarted:
			// Play the game
			if g.game.Play() {
				g.eState = EngineState_PLAYING
				g.NeedNtfAndReason(NotifyGameStateReason_NEW_GAME)
			}
		}
	case EngineState_PLAYING:
		switch input.Type {
		case PlayerActed:
			g.game.HandleActions(input.PlayerAct)
			g.NeedNtfAndReason(NotifyGameStateReason_NEW_ACTION)
		case NextGame:
		}
	case EngineState_PAUSED:
		switch input.Type {
		case NextGame:
			g.eState = EngineState_PLAYING
		}
	}

	// Try to notify the game state each time the game engine is run
	g.NotifyGameState()
}

func (g *GameEngine) NeedNtfAndReason(reason NotifyGameStateReason) {
	mylog.Infof("Need to notify the game state: %v\n", reason)
	g.ntfReason = reason
}

func (g *GameEngine) NotifyGameState() {
	switch g.ntfReason {
	case NotifyGameStateReason_DONT_NOFITY:
		return
	case NofityGameStateReason_ALL,
		NotifyGameStateReason_NEW_ROUND,
		NotifyGameStateReason_UPDATE_PLAYER,
		NotifyGameStateReason_NEW_ACTION:
		// Special case: Notifying
	case NotifyGameStateReason_NEW_GAME:
		// Notify directly to all players if they have new cards
		for _, player := range g.playerMgr.players {
			if player != nil && player.HasNewCards() {
				player.NotifyPlayerIfNewHand()
			}
		}
	}

	if g.room != nil {
		g.room.BroadcastMessageToYourRoom(&msgpb.ServerMessage{
			Message: &msgpb.ServerMessage_GameState{
				GameState: g.SyncGameState(),
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
			MaxPlayers:  6,
			MinPlayers:  2,
			SmallBlind:  10,
			BigBlind:    20,
			TimePerTurn: 0, // 0 means no limit
		},
		g.playerMgr,
		NewDeck(),
		g.gotoState,
	)
}

// HandleWaitForPlayers handles the EngineState_WAIT_FOR_PLAYING state.
func (g *GameEngine) HandleJoiningPlayer(player Player) {
	if player != nil {
		fmt.Printf("Player %s joined the game\n", player.Name())
		g.playerMgr.AddPlayer(player.Position(), player)
	}
}

func (g *GameEngine) HandleLeavingPlayer(player Player) {
	if player != nil {
		g.playerMgr.RemovePlayer(player.Position())
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

// SyncGameState synchronizes the game state.
func (g *GameEngine) SyncGameState() *msgpb.GameState {
	// Create a message to sync the game state

	fakeGameState := &msgpb.GameState{
		Players:        make([]*msgpb.Player, 0),
		PotSize:        1000,
		DealerId:       0,
		CommunityCards: make([]*msgpb.Card, 0),
		CurrentBet:     0,
		CurrentRound:   msgpb.RoundStateType_PREFLOP,
		FinalResult: &msgpb.Result{
			WinnerPosition: 3,
			WonChip:        1000,
			ShowingCards: []*msgpb.PeerState{
				{
					TablePos: 4,
					PlayerCards: []*msgpb.Card{
						{Suit: msgpb.SuitType_SPADES, Rank: msgpb.RankType_ACE},
						{Suit: msgpb.SuitType_DIAMONDS, Rank: msgpb.RankType_KING},
					},
				},
				{
					TablePos: 4,
					PlayerCards: []*msgpb.Card{
						{Suit: msgpb.SuitType_SPADES, Rank: msgpb.RankType_ACE},
						{Suit: msgpb.SuitType_DIAMONDS, Rank: msgpb.RankType_KING},
					},
				},
			},
		},
	}
	fakeGameState.CommunityCards = []*msgpb.Card{
		{Suit: msgpb.SuitType_SPADES, Rank: msgpb.RankType_ACE},
		{Suit: msgpb.SuitType_DIAMONDS, Rank: msgpb.RankType_KING},
		{Suit: msgpb.SuitType_CLUBS, Rank: msgpb.RankType_QUEEN},
		{Suit: msgpb.SuitType_HEARTS, Rank: msgpb.RankType_JACK},
		{Suit: msgpb.SuitType_SPADES, Rank: msgpb.RankType_TEN},
	}
	fakeGameState.Players = []*msgpb.Player{
		{
			Name:          "player1",
			Chips:         1500,
			TablePosition: 0,
			Status:        "Playing",
		},
		{
			Name:          "player2",
			Chips:         2000,
			TablePosition: 1,
			Status:        "Wait4Act",
		},
		{
			Name:          "player3",
			Chips:         4000,
			TablePosition: 2,
			Status:        "Fold",
		},
		{
			Name:          "player4",
			Chips:         3000,
			TablePosition: 4,
			Status:        "Check",
		},
		{
			Name:          "player5",
			Chips:         4800,
			TablePosition: 5,
			Status:        "Raise",
			CurrentBet:    200,
		},
	}

	return fakeGameState

	syncMsg := &msgpb.GameState{
		Players:        make([]*msgpb.Player, 0),
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
			syncMsg.Players = append(syncMsg.Players, &msgpb.Player{
				Name:          player.Name(),
				TablePosition: int32(player.Position()),
				Chips:         int32(player.Chips()),
				IsDealer:      player.Position() == g.game.gs.ButtonPosition,
				Status:        player.Status().String(),
				CurrentBet:    int32(player.CurrentBet()),
			})
		}
	}

	return syncMsg
}
