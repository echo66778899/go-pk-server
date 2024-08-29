package engine

import (
	"context"
	"fmt"
)

var DebugMode = true

type ActionType int

const (
	Unspecified ActionType = iota
	PlayerJoined
	PlayerLeft
	PlayerReady
	NextEngState
	GameEnded
	PlayerActed
	NextRound
	UpdatePlayer
)

type Action struct {
	// Common fields for all actions
	Type ActionType
	// Possible fields for an action
	PlayerAct  ActionIf
	PlayerInfo Player
}

// GameEngine represents the game engine.
type GameEngineIf interface {
	StartEngine()
	StopEngine()
	PlayerJoin(player Player)
	StartGame()
	Ready()
	PlayerAction(action ActionIf)
}

type EngineState int

const (
	RoomCreated EngineState = iota
	WaitForPlayers
	WaitForPlayerActions
	WaitForNextRound
)

// Overwrite string method for EngineState
func (e EngineState) String() string {
	return [...]string{"RoomCreated", "WaitForPlayers", "WaitForPlayerActions", "WaitForNextRound"}[e]
}

type GameEngine struct {
	gameSessionID int
	State         GameState
	playerMgr     *PlayerManager
	game          *Game

	eState      EngineState
	eventDriven bool

	ctx           context.Context
	cancel        context.CancelFunc
	ActionChannel chan Action
}

// NewGameEngine creates a new instance of the game engine.
func NewGameEngine() GameEngineIf {
	// Add your initialization code here
	return &GameEngine{
		gameSessionID: 1,
		eState:        RoomCreated,
		eventDriven:   false,
		State:         GameState{},
		playerMgr:     NewPlayerManager(8),
		ActionChannel: make(chan Action, 10), // Change to buffered channel with capacity 10

	}
}

// StartEngine starts the game.
func (g *GameEngine) StartEngine() {
	// Run the game engine with go routine
	if g.eventDriven {
		g.ctx, g.cancel = context.WithCancel(context.Background())
		go g.EngineLoop(g.ctx)
	}
}

func (g *GameEngine) StopEngine() {
	if g.eventDriven {
		g.cancel()
	}
}

func (g *GameEngine) PlayerJoin(player Player) {
	action := Action{Type: PlayerJoined, PlayerInfo: player}
	g.processActions(action)
}

func (g *GameEngine) StartGame() {
	// Log the game start
	fmt.Println("Game started")
	// send action to start the game
	act := Action{Type: ActionType(WaitForPlayers)}
	g.processActions(act)
}

// PerformAction performs the specified action for the given player.
func (g *GameEngine) PlayerAction(action ActionIf) {
	// Send action to game engine
	act := Action{Type: PlayerActed, PlayerAct: action}
	g.processActions(act)
}

func (g *GameEngine) Ready() {
	action := Action{Type: PlayerReady}
	g.processActions(action)
}

func (g *GameEngine) processActions(action Action) {
	if g.eventDriven {
		g.ActionChannel <- action
	} else {
		g.RunGameEngine(action)
	}
}

// EngineLoop runs the game engine in a loop.
func (g *GameEngine) EngineLoop(ctx context.Context) {
	for {
		select {
		case action := <-g.ActionChannel:
			g.RunGameEngine(action)
		case <-ctx.Done():
			// Game ended
			return
		}
	}
}

func (g *GameEngine) RunGameEngine(action Action) {
	switch g.eState {
	case RoomCreated:
		// Room is created log
		g.HandleRoomCreated()
		g.eState = WaitForPlayers
	case WaitForPlayers:
		// Wait for players to join
		if action.Type == PlayerJoined && action.PlayerInfo != nil {
			g.HandleWaitForPlayers(action.PlayerInfo)
		} else if action.Type == PlayerReady {
			g.game = NewGame(GameSetting{
				NumPlayers:   g.playerMgr.numberOfSlots,
				MaxStackSize: 1000,
				MinStackSize: 100,
				SmallBlind:   10,
				BigBlind:     20,
			}, g.playerMgr, NewDeck())
			// Play the game
			g.game.Play()
			g.eState = WaitForPlayerActions
		}
	case WaitForPlayerActions:
		if action.Type == PlayerActed {
			g.game.HandleActions(action.PlayerAct)
		}
	case WaitForNextRound:
		// Game over
	}

	g.NotifyGameState()
}

func (g *GameEngine) NotifyGameState() {
	// Todo: If game state has changed, notify the clients

	// Log the current game state
	fmt.Printf("===============\nCurrent game state: %v\n", g.eState)
}

// HandleRoomCreated handles the RoomCreated state.
func (g *GameEngine) HandleRoomCreated() {
}

// HandleWaitForPlayers handles the WaitForPlayers state.
func (g *GameEngine) HandleWaitForPlayers(player Player) {
	fmt.Printf("Player %s joined the game\n", player.Name())
	g.playerMgr.AddPlayer(player.Position(), player)
	// Log the player count
}

// HandleShowdown handles the Showdown state.
func (g *GameEngine) HandleShowdown() {
	// Add your logic for the Showdown state here
}

// HandleGameOver handles the GameOver state.
func (g *GameEngine) EvaluateHands() {
	// Compare hands and determine the winner
	// Find the best hand among all players

	// Print the winner
	// fmt.Printf("The winner is %s with hand: [%s] (%s)\n", winningPlayer.Name, winningPlayer.Hand.BestHand(), winningPlayer.Hand.HandRankingString())
}

// IsGameOver checks if the game is over.
func (g *GameEngine) IsGameOver() bool {
	// Add your game over condition logic here
	return false
}

func (g *GameEngine) SummarizeRound() {
	// Log the round summary
	fmt.Println("Round summary:")
	// 	for _, player := range g.State.Players {
	// 		if player.Status() == Folded {
	// 			fmt.Printf("Player %s: Fold\n", player.Name)
	// 		} else {
	// 			fmt.Printf("Player %s: %s (%s)\n", player.Name, player.ShowHand(), player.ShowHand())
	// 		}
	// 		fmt.Printf("Player %s has %d chips\n", player.Name, player.Chips)
	// 	}
}
