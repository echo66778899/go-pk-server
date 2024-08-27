package engine

import (
	"context"
	"fmt"
)

// GameEngine represents the game engine.
type GameEngineIf interface {
	StartGame()
	StopGame()
	PlayerJoin(roomId int, player Player)
	PlayerLeft(roomId int, player Player)
	Ready()
	PlayerAction(playerIdx int, action PlayerAction)
}

type EngineState int

const (
	RoomCreated EngineState = iota
	WaitForPlayers
	DealCards
	Preflop
	Flop
	Turn
	River
	Showdown
	WaitForNextRound
)

// Overwrite string method for EngineState
func (e EngineState) String() string {
	return [...]string{"RoomCreated", "WaitForPlayers", "DealCards", "Preflop", "Flop", "Turn", "River", "Showdown", "WaitForNextRound"}[e]
}

type GameEvent int

const (
	Unspecified GameEvent = iota
	PlayerJoined
	PlayerLeft
	PlayerReady
	NextEngState
	GameEnded
	PlayerActed
	NextRound
	UpdatePlayer
)

type EventData struct {
	PlayerInfo Player
	PlayerAct  PlayerAction
}

type Event struct {
	Type      GameEvent
	EventData EventData
}

type GameEngine struct {
	ID          int
	State       GameState
	PlayingDeck *Deck
	eState      EngineState

	ctx          context.Context
	cancel       context.CancelFunc
	EventChannel chan Event
}

// NewGameEngine creates a new instance of the game engine.
func NewGameEngine() GameEngineIf {
	// Add your initialization code here
	return &GameEngine{
		ID:           1,
		PlayingDeck:  NewDeck(),
		eState:       RoomCreated,
		EventChannel: make(chan Event, 100), // Change to buffered channel with capacity 10
	}
}

// StartGame starts the game.
func (g *GameEngine) StartGame() {
	// Log for starting the game
	fmt.Println("Starting the game")
	event := Event{Type: NextEngState}
	g.EventChannel <- event

	// Run the game engine with go routine
	g.ctx, g.cancel = context.WithCancel(context.Background())
	go g.RunEngine(g.ctx)
}

func (g *GameEngine) StopGame() {
	g.cancel()
}

func (g *GameEngine) PlayerJoin(roomId int, player Player) {
	// Log the parameters
	fmt.Printf("Player %s joined room %d\n", player.Name, roomId)
	// log debug eState and g.ID
	fmt.Printf("eState: %v, ID: %v\n", g.eState, g.ID)
	if g.eState == WaitForPlayers && roomId == g.ID {
		// Log the player joining the game
		fmt.Printf("%s joined the game\n", player.Name)
		event := Event{Type: PlayerJoined, EventData: EventData{PlayerInfo: player}}
		g.EventChannel <- event
	}
}

func (g *GameEngine) PlayerLeft(roomId int, player Player) {
	if g.eState == WaitForPlayers && roomId == g.ID {
		fmt.Printf("%s left the game\n", player.Name)
		event := Event{Type: PlayerLeft, EventData: EventData{PlayerInfo: player}}
		g.EventChannel <- event
	}
}

// PerformAction performs the specified action for the given player.
func (g *GameEngine) PlayerAction(playerIdx int, action PlayerAction) {
	// Log the action
	fmt.Printf("Player %d performs action: %v\n", playerIdx, action)
	event := Event{Type: PlayerActed, EventData: EventData{PlayerAct: action}}
	g.EventChannel <- event
}

func (g *GameEngine) Ready() {
	event := Event{Type: PlayerReady}
	g.EventChannel <- event
}

// RunEngine runs the game engine in a loop.
func (g *GameEngine) RunEngine(ctx context.Context) {
	for {
		select {
		case event := <-g.EventChannel:
			g.RunGameLogic(event)
		case <-ctx.Done():
			// Game ended
			return
		}
	}
}

func (g *GameEngine) NextEState(nexteState EngineState) {
	// log the state transition
	fmt.Printf("Transitioning from %v to %v\n", g.eState, nexteState)
	g.eState = nexteState
	event := Event{Type: NextEngState}
	g.EventChannel <- event
}

// GetPlayerIndex returns the index of the player with the given ID.
func (g *GameEngine) GetPlayerIndex(playerID int) int {
	for i, player := range g.State.Players {
		if player.ID == playerID {
			return i
		}
	}
	return -1
}

func (g *GameEngine) RunGameLogic(event Event) {
	switch g.eState {
	case RoomCreated:
		// Room is created log
		g.HandleRoomCreated()
	case WaitForPlayers:
		// Wait for players to join
		if event.Type == PlayerJoined || event.Type == PlayerLeft {
			g.HandleWaitForPlayers(event.EventData.PlayerInfo)
		} else if event.Type == PlayerReady {
			g.ResetGameState()
			g.ResetPlayerState()
			g.NextEState(DealCards)
		}
	case DealCards:
		// Deal cards to players
		g.DealCards()
		g.NextEState(Preflop)
	case Preflop:
		// Handle player actions for betting at the preflop round
		if event.Type == PlayerActed {
			if g.HandlePlayerActions(event.EventData.PlayerAct) {
				g.DealFlop()
				g.ResetRoundBet()
				g.NextEState(Flop)
			}
		}
	case Flop:
		// Flop round
		if event.Type == PlayerActed {
			if g.HandlePlayerActions(event.EventData.PlayerAct) {
				g.DealTurn()
				g.ResetRoundBet()
				g.NextEState(Turn)
			}
		}
	case Turn:
		// Turn round
		if event.Type == PlayerActed {
			if g.HandlePlayerActions(event.EventData.PlayerAct) {
				g.DealRiver()
				g.ResetRoundBet()
				g.NextEState(River)
			}
		}
	case River:
		// River round
		if event.Type == PlayerActed {
			if g.HandlePlayerActions(event.EventData.PlayerAct) {
				g.ResetRoundBet()
				g.NextEState(Showdown)
			}
		}
	case Showdown:
		// Showdown
		g.EvaluateHands()
		g.SummarizeRound()
		g.NextEState(WaitForNextRound)
	case WaitForNextRound:
		// Game over
	}
}

// HandleDealCards handles the DealCards state.
func (g *GameEngine) DealCards() {
	g.State.CurrentPlayerIdx = g.State.CurrentButtonIdx
	// Shuffle the deck before dealing
	g.PlayingDeck.Shuffle()
	g.PlayingDeck.Cut()

	// Deal cards to players in a round-robin fashion
	for i := 0; i < 2; i++ {
		for j := 0; j < len(g.State.Players); j++ {
			player := &g.State.Players[j]
			player.Hand.SetCard(g.PlayingDeck.Draw(), i)
		}
	}

	// Log all the players' hands
	for _, player := range g.State.Players {
		player.PrintHand()
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
	fmt.Printf("Flow Board: [%s]\n", g.State.CommunityCards.String())
}

func (g *GameEngine) DealTurn() {
	// Burn a card
	_ = g.PlayingDeck.Draw()
	// Add a card to the community cards
	g.State.CommunityCards.AddCard(g.PlayingDeck.Draw())
	// Print the community cards
	fmt.Printf("Turn Board: [%s]\n", g.State.CommunityCards.String())
}

func (g *GameEngine) DealRiver() {
	// Burn a card
	_ = g.PlayingDeck.Draw()
	// Add a card to the community cards
	g.State.CommunityCards.AddCard(g.PlayingDeck.Draw())
	fmt.Printf("River Board: [%s]\n", g.State.CommunityCards.String())
}

func (g *GameEngine) ResetGameState() {
	g.State.CommunityCards.Reset()
	g.State.CurrentButtonIdx = (g.State.CurrentButtonIdx + 1) % len(g.State.Players)
	g.State.CurrentPlayerIdx = g.State.CurrentButtonIdx
	g.State.CurrentBet = 0
	g.State.NumPlayingPlayer = len(g.State.Players)
	g.State.Pots.ResetPot()
}

func (g *GameEngine) ResetPlayerState() {
	for i := range g.State.Players {
		g.State.Players[i].HasFolded = false
		g.State.Players[i].HasActed = false
		g.State.Players[i].Hand.Reset()
	}
	g.ResetRoundBet()
}

func (g *GameEngine) ResetRoundBet() {
	g.State.CurrentBet = 0
	g.State.CurrentPlayerIdx = g.State.CurrentButtonIdx
	for i := range g.State.Players {
		g.State.Players[i].CurrentBet = 0
		g.State.Players[i].HasActed = false
	}
}

// HandleRoomCreated handles the RoomCreated state.
func (g *GameEngine) HandleRoomCreated() {
	// Add your logic for the RoomCreated state here

	g.NextEState(WaitForPlayers)
}

// HandleWaitForPlayers handles the WaitForPlayers state.
func (g *GameEngine) HandleWaitForPlayers(player Player) {
	g.State.Players = append(g.State.Players, player)
}

// HandlePlayerActions handles the Preflop state.
func (g *GameEngine) HandlePlayerActions(action PlayerAction) bool {
	isDoneBetting := false
	// log the player idx and current player idx
	fmt.Printf("PlayerIdx: %v, CurrentPlayerIdx: %v\n", action.PlayerIdx, g.State.CurrentPlayerIdx)
	g.HandleBetting(action)
	g.State.CurrentPlayerIdx = (g.State.CurrentPlayerIdx + 1) % len(g.State.Players)
	// Check if all players have acted
	isDoneBetting = true
	for _, player := range g.State.Players {
		if !player.HasFolded && !player.HasActed {
			isDoneBetting = false
			break
		}
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
	var winningPlayer *Player

	for i := range g.State.Players {
		player := &g.State.Players[i]
		if player.HasFolded {
			continue
		}
		rank := player.Hand.CalcBestHand(&g.State.CommunityCards)
		if winningPlayer == nil || rank > winningPlayer.Hand.calRank {
			winningPlayer = player
		} else if rank == winningPlayer.Hand.calRank {
			// Compare tiebreakers
			if result := compareTiebreakers(player.Hand.bestTiebreaker, winningPlayer.Hand.bestTiebreaker); result > 0 {
				winningPlayer = player
			} else if result == 0 {
				// Split pot
				winningPlayer = nil

			}
		}
	}
	if winningPlayer != nil {
		winningPlayer.Chips += g.State.Pots.GetPotAmount()
	}

	// Print the winner
	fmt.Printf("The winner is %s with hand: [%s] (%s)\n", winningPlayer.Name, winningPlayer.Hand.BestHand(), winningPlayer.Hand.HandRankString())
}

// IsGameOver checks if the game is over.
func (g *GameEngine) IsGameOver() bool {
	// Add your game over condition logic here
	return false
}

func (g *GameEngine) HandleBetting(action PlayerAction) {
	// log the action
	fmt.Printf("Handle betting from Player %d: %v\n", action.PlayerIdx, action)
	player := &g.State.Players[action.PlayerIdx]
	switch action.Type {
	case Fold:
		player.HasFolded = true
		fmt.Printf("%s folds\n", player.Name)
	case Call:
		callAmount := g.State.CurrentBet - player.CurrentBet
		if player.Chips >= callAmount {
			player.Chips -= callAmount
			player.CurrentBet += callAmount
			g.State.Pots.AddToPot(callAmount)
			fmt.Printf("%s calls with %d chips\n", player.Name, callAmount)
			// Update the hasActed flag for the player
			player.HasActed = true
		} else {
			fmt.Printf("%s doesn't have enough chips to call\n", player.Name)
			// Player goes all-in
			player.CurrentBet += player.Chips

			// Slip pot
			g.State.Pots.AddToPot(player.Chips)

			player.Chips = 0
			fmt.Printf("%s goes all-in with %d chips\n", player.Name, player.CurrentBet)
			// Update the hasActed flag for the player
			player.HasActed = true
		}
	case Raise, Bet:
		raiseAmount := action.Bet
		if raiseAmount > player.Chips {
			fmt.Printf("%s doesn't have enough chips to raise\n", player.Name)
		} else if raiseAmount < g.State.CurrentBet {
			fmt.Printf("%s's raise is too small\n", player.Name)
		} else {
			actionLog := "bets"
			if g.State.CurrentBet != 0 {
				actionLog = "raises to"
			}
			raiseAmount += g.State.CurrentBet - player.CurrentBet
			player.Chips -= raiseAmount
			player.CurrentBet += raiseAmount
			g.State.CurrentBet = player.CurrentBet
			g.State.Pots.AddToPot(raiseAmount)
			fmt.Printf("%s %s %d chips\n", player.Name, actionLog, g.State.CurrentBet)
			player.HasActed = true

			// Reset the hasActed flag for all players expect this player
			for i := range g.State.Players {
				if i != action.PlayerIdx {
					g.State.Players[i].HasActed = false
				}
			}
		}
	case Check:
		if g.State.CurrentBet == 0 {
			fmt.Printf("%s checks\n", player.Name)
			// Update the hasActed flag for the player
			player.HasActed = true
		} else {
			fmt.Printf("%s cannot check because there's a bet\n", player.Name)
		}
	}
}

func (g *GameEngine) SummarizeRound() {
	// Log the round summary
	fmt.Println("Round summary:")
	for _, player := range g.State.Players {
		if player.HasFolded {
			fmt.Printf("Player %s: Fold\n", player.Name)
		} else {
			fmt.Printf("Player %s: %s (%s)\n", player.Name, player.Hand.BestHand(), player.Hand.HandRankString())
		}
		fmt.Printf("Player %s has %d chips\n", player.Name, player.Chips)
	}
}
