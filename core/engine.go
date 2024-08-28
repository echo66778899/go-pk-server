package engine

import (
	"context"
	"fmt"
)

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

	DumpGameState()
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
	PlayingDeck   *Deck
	eState        EngineState
	eventDriven   bool

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
		State: GameState{GameSetting: GameSetting{
			SmallBlind: 10,
			BigBlind:   20,
		}},
		PlayingDeck:   NewDeck(),
		ActionChannel: make(chan Action, 10), // Change to buffered channel with capacity 10

	}
}

// StartEngine starts the game.
func (g *GameEngine) StartEngine() {
	// Run the game engine with go routine
	if g.eventDriven {
		g.ctx, g.cancel = context.WithCancel(context.Background())
		go g.RunEngine(g.ctx)
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
		g.UpdateGameState(action)
	}
}

// RunEngine runs the game engine in a loop.
func (g *GameEngine) RunEngine(ctx context.Context) {
	for {
		select {
		case action := <-g.ActionChannel:
			g.UpdateGameState(action)
		case <-ctx.Done():
			// Game ended
			return
		}
	}
}

func (g *GameEngine) UpdateGameState(action Action) {
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
			g.NewGame()
			g.TakeBlinds()
			g.DealCards()
			g.eState = WaitForPlayerActions
		}
	case WaitForPlayerActions:
		if action.Type == PlayerActed {
			g.HandleGameLogic(action.PlayerAct)
		}
	case WaitForNextRound:
		// Game over
	}

	g.NotifyGameState()
}

func (g *GameEngine) NewGame() {
	// Reset the game state and related states to be ready for a new game
	g.State.ResetBeforePlay()
}

func (g *GameEngine) TakeBlinds() {
	// Take the blinds from the players
	sbPlayer := g.State.GetSmallBlindPlayer()
	bbPlayer := g.State.GetBigBlindPlayer()
	sbPlayer.TakeChips(g.State.GameSetting.SmallBlind)
	sbPlayer.UpdateCurrentBet(g.State.GameSetting.SmallBlind)
	bbPlayer.TakeChips(g.State.GameSetting.BigBlind)
	bbPlayer.UpdateCurrentBet(g.State.GameSetting.BigBlind)
	// Add the blinds to the pot
	g.State.Pots.AddToPot(sbPlayer.CurrentBet())
	g.State.Pots.AddToPot(bbPlayer.CurrentBet())
	// Update state
	sbPlayer.UpdateStatus(Betted)
	bbPlayer.UpdateStatus(Betted)
	g.State.NextActivePlayer(g.State.ButtonPosition).UpdateStatus(WaitForAct)
}

func (g *GameEngine) NotifyGameState() {
	// Log the current game state
	fmt.Printf("===============\nCurrent game state: %v\n", g.eState)
}

func (g *GameEngine) HandleGameLogic(player ActionIf) {
	switch g.State.CurrentRound {
	case PreFlop:
		// Handle player actions for betting at the preflop round
		if g.HandlePlayerActions(player) {
			g.DealFlop()
			g.State.ResetBettingState()
			g.State.CurrentRound = Flop
		}
	case Flop:
		// Flop round
		if g.HandlePlayerActions(player) {
			g.DealTurn()
			g.State.ResetBettingState()
			g.State.CurrentRound = Turn
		}
	case Turn:
		// Turn round
		if g.HandlePlayerActions(player) {
			g.DealRiver()
			g.State.ResetBettingState()
			g.State.CurrentRound = River
		}
	case River:
		// River round
		if g.HandlePlayerActions(player) {
			g.State.CurrentRound = Showdown
		}
	case Showdown:
		// Showdown
		g.EvaluateHands()
		g.SummarizeRound()
	}
}

// HandleDealCards handles the DealCards state.
func (g *GameEngine) DealCards() {
	// Shuffle the deck before dealing
	g.PlayingDeck.Shuffle()
	g.PlayingDeck.CutTheCard()

	// Deal cards to players in a round-robin fashion
	for i := 0; i < 2; i++ {
		for j := 0; j < len(g.State.Players); j++ {
			player := g.State.Players[j]
			player.DealCard(g.PlayingDeck.Draw(), i)
		}
	}

	// Log all the players' hands
	for _, player := range g.State.Players {
		if player == nil {
			// log error
			continue
		}
		fmt.Printf("%s's hand: [%s]\n", player.Name(), player.ShowHand().String())
	}
}

func (g *GameEngine) DealFlop() {
	// Burn a card
	_ = g.PlayingDeck.Draw()
	// Add 3 cards to the community cards
	g.State.CommunityCards.AddCard(g.PlayingDeck.Draw())
	g.State.CommunityCards.AddCard(g.PlayingDeck.Draw())
	g.State.CommunityCards.AddCard(g.PlayingDeck.Draw())
	// Print the community cards with for loop
	fmt.Printf("===============\nFlop Board: %s\n===============\n",
		g.State.CommunityCards.String())
}

func (g *GameEngine) DealTurn() {
	// Burn a card
	_ = g.PlayingDeck.Draw()
	// Add a card to the community cards
	g.State.CommunityCards.AddCard(g.PlayingDeck.Draw())
	// Print the community cards
	fmt.Printf("===============\nTurn Board: %s\n===============\n",
		g.State.CommunityCards.String())
}

func (g *GameEngine) DealRiver() {
	// Burn a card
	_ = g.PlayingDeck.Draw()
	// Add a card to the community cards
	g.State.CommunityCards.AddCard(g.PlayingDeck.Draw())
	fmt.Printf("===============\nRiver Board: %s\n===============\n",
		g.State.CommunityCards.String())
}

// HandleRoomCreated handles the RoomCreated state.
func (g *GameEngine) HandleRoomCreated() {
}

// HandleWaitForPlayers handles the WaitForPlayers state.
func (g *GameEngine) HandleWaitForPlayers(player Player) {
	fmt.Printf("Player %s joined the game\n", player.Name())
	g.State.Players = append(g.State.Players, player)
	// Assign the player to the game position
	player.UpdatePosition(len(g.State.Players) - 1)
	// Log the player count
	fmt.Printf("Number of players: %d\n", len(g.State.Players))
}

// HandlePlayerActions handles the Preflop state.
func (g *GameEngine) HandlePlayerActions(action ActionIf) bool {
	isDoneBetting := false
	// Log the player action
	fmt.Printf("Handling %s action\n", action.What())

	action.Execute(&g.State)
	// If all players have acted, return true
	if g.State.AllPlayersActed() {
		isDoneBetting = true
	}
	return isDoneBetting
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
	// fmt.Printf("The winner is %s with hand: [%s] (%s)\n", winningPlayer.Name, winningPlayer.Hand.BestHand(), winningPlayer.Hand.HandRankString())
}

// IsGameOver checks if the game is over.
func (g *GameEngine) IsGameOver() bool {
	// Add your game over condition logic here
	return false
}

func (g *GameEngine) SummarizeRound() {
	// Log the round summary
	fmt.Println("Round summary:")
	g.State.Pots.ResetPot()
	// 	for _, player := range g.State.Players {
	// 		if player.Status() == Folded {
	// 			fmt.Printf("Player %s: Fold\n", player.Name)
	// 		} else {
	// 			fmt.Printf("Player %s: %s (%s)\n", player.Name, player.ShowHand(), player.ShowHand())
	// 		}
	// 		fmt.Printf("Player %s has %d chips\n", player.Name, player.Chips)
	// 	}
}

// DumpGameState dumps the current game state.
func (g *GameEngine) DumpGameState() {
	// Log the game state
	fmt.Println("!!!!!!!Dumping game state!!!!!!!!!!")
	fmt.Printf("Game state: %v\n", g.State)
	fmt.Printf("Button Player: %v\n", g.State.GetButtonPlayer().Name())
	fmt.Printf("Small Blind Player: %v\n", g.State.GetSmallBlindPlayer().Name())
	fmt.Printf("Big Blind Player: %v\n", g.State.GetBigBlindPlayer().Name())
	fmt.Printf("Deck: %v\n", g.PlayingDeck)
	fmt.Println("Players:")
	for _, player := range g.State.Players {
		fmt.Printf("%s\n", player)
	}

}
