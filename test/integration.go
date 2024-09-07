// Main to run the integration tests

package main

import (
	"fmt"
	pk_eng "go-pk-server/core"
)

func NewPlayerAction(position int, actionType pk_eng.msgpb.PlayerGameActionType, amount int) pk_eng.ActionIf {
	return &pk_eng.PlayerAction{
		PlayerPosition: position,
		ActionType:     actionType,
		Amount:         amount,
	}
}

const (
	A = 0
	B = 1
	C = 2
	D = 3
	E = 4
)

func main() {
	game := pk_eng.NewGameEngine()
	game.StartEngine(false)

	// Create 4 players
	playerA := pk_eng.NewOnlinePlayer("A", 123)
	playerA.UpdatePosition(0)
	playerB := pk_eng.NewOnlinePlayer("B", 456)
	playerB.UpdatePosition(1)
	playerC := pk_eng.NewOnlinePlayer("C", 789)
	playerC.UpdatePosition(3)
	playerD := pk_eng.NewOnlinePlayer("D", 101)
	playerD.UpdatePosition(5)
	playerE := pk_eng.NewOnlinePlayer("E", 121)
	playerE.UpdatePosition(6)

	// Append players to a slice
	players := []*pk_eng.OnlinePlayer{playerA, playerB, playerC, playerD, playerE}
	for _, p := range players {
		// Add chips to players
		p.AddChips(1000)
		game.PlayerJoin(p)
	}

	game.StartGame()

	// Test win at preflop
	TestWinAtPreflop(game, players)
	// Test win at flop + 2 continue games
	TestWinAtFlop(game, players)
	// Test win at turn + 3 continue games
	TestWinAtTurn(game, players)
}

func TestWinAtPreflop(g pk_eng.GameEngineIf, players []*pk_eng.OnlinePlayer) {
	fmt.Println("TestWinAtPreflop")
	defer fmt.Println("End TestWinAtPreflop")

	g.NextGame()
	// Preflop, UTG action first
	// A dealer, B small blind, C big blind, D UTG and E cutoff
	// UTG call
	g.PlayerAction(players[D].NewReAct(pk_eng.Call, 20))
	g.PlayerAction(players[E].NewReAct(pk_eng.Fold, 0))
	g.PlayerAction(players[A].NewReAct(pk_eng.Raise, 100))
	g.PlayerAction(players[B].NewReAct(pk_eng.Call, 100))
	g.PlayerAction(players[C].NewReAct(pk_eng.Raise, 300))
	g.PlayerAction(players[D].NewReAct(pk_eng.Fold, 0))
	g.PlayerAction(players[A].NewReAct(pk_eng.Fold, 0))
	g.PlayerAction(players[B].NewReAct(pk_eng.Fold, 0))

	// All players chips
	total := 0
	for _, p := range players {
		fmt.Printf("Test=Player %s chips: %d\n", p.Name(), p.Chips())
		total += p.Chips()
	}
	if total != 5000 {
		panic("Total chips should be 5000")
	}
}

func TestWinAtFlop(g pk_eng.GameEngineIf, players []*pk_eng.OnlinePlayer) {
	fmt.Println("TestWinAtFlop")
	defer fmt.Println("End TestWinAtFlop")

	g.NextGame()
	// Flop, UTG action first
	// B dealer, C small blind, D big blind, E UTG and A cutoff
	// UTG check
	g.PlayerAction(players[E].NewReAct(pk_eng.Call, 0))
	g.PlayerAction(players[A].NewReAct(pk_eng.Call, 0))
	g.PlayerAction(players[B].NewReAct(pk_eng.Call, 0))
	g.PlayerAction(players[C].NewReAct(pk_eng.Call, 0))
	g.PlayerAction(players[D].NewReAct(pk_eng.Check, 0))

	// Flop dealed
	// C first to act
	g.PlayerAction(players[C].NewReAct(pk_eng.Raise, 100))
	g.PlayerAction(players[D].NewReAct(pk_eng.Fold, 0))
	g.PlayerAction(players[E].NewReAct(pk_eng.Fold, 0))
	g.PlayerAction(players[A].NewReAct(pk_eng.Fold, 0))
	g.PlayerAction(players[B].NewReAct(pk_eng.Fold, 0))

	// All players chips
	total := 0
	for _, p := range players {
		fmt.Printf("Test=Player %s chips: %d\n", p.Name(), p.Chips())
		total += p.Chips()
	}
	if total != 5000 {
		panic("Total chips should be 5000")
	}
}

func TestWinAtTurn(g pk_eng.GameEngineIf, players []*pk_eng.OnlinePlayer) {
	fmt.Println("TestWinAtTurn")
	defer fmt.Println("End TestWinAtTurn")

	g.NextGame()
	// Turn, UTG action first
	// C dealer, D small blind, E big blind, A UTG and B cutoff
	// UTG check
	g.PlayerAction(players[A].NewReAct(pk_eng.Call, 0))
	g.PlayerAction(players[B].NewReAct(pk_eng.Call, 0))
	g.PlayerAction(players[C].NewReAct(pk_eng.Call, 0))
	g.PlayerAction(players[D].NewReAct(pk_eng.Call, 0))
	g.PlayerAction(players[E].NewReAct(pk_eng.Check, 0))

	// Flop dealed
	// D first to act
	g.PlayerAction(players[D].NewReAct(pk_eng.Raise, 200))
	g.PlayerAction(players[E].NewReAct(pk_eng.Call, 0))
	g.PlayerAction(players[A].NewReAct(pk_eng.Call, 0))
	g.PlayerAction(players[B].NewReAct(pk_eng.Fold, 0))
	g.PlayerAction(players[C].NewReAct(pk_eng.Fold, 0))

	// Turn, D action first
	g.PlayerAction(players[D].NewReAct(pk_eng.Raise, 100))
	g.PlayerAction(players[E].NewReAct(pk_eng.Fold, 0))
	g.PlayerAction(players[A].NewReAct(pk_eng.Fold, 0))

	// All players chips
	total := 0
	for _, p := range players {
		fmt.Printf("Test=Player %s chips: %d\n", p.Name(), p.Chips())
		total += p.Chips()
	}
	if total != 5000 {
		panic("Total chips should be 5000")
	}
}

func TestWinAtRiver(g pk_eng.GameEngineIf, players []*pk_eng.OnlinePlayer) {
	fmt.Println("TestWinAtTurn")
	defer fmt.Println("End TestWinAtTurn")

	g.NextGame()
	// Turn, UTG action first
	// C dealer, D small blind, E big blind, A UTG and B cutoff
	// UTG check
	g.PlayerAction(players[A].NewReAct(pk_eng.Call, 0))
	g.PlayerAction(players[B].NewReAct(pk_eng.Call, 0))
	g.PlayerAction(players[C].NewReAct(pk_eng.Call, 0))
	g.PlayerAction(players[D].NewReAct(pk_eng.Call, 0))
	g.PlayerAction(players[E].NewReAct(pk_eng.Check, 0))

	// Flop dealed
	// D first to act
	g.PlayerAction(players[D].NewReAct(pk_eng.Raise, 200))
	g.PlayerAction(players[E].NewReAct(pk_eng.Call, 0))
	g.PlayerAction(players[A].NewReAct(pk_eng.Call, 0))
	g.PlayerAction(players[B].NewReAct(pk_eng.Fold, 0))
	g.PlayerAction(players[C].NewReAct(pk_eng.Fold, 0))

	// Turn, D action first
	g.PlayerAction(players[D].NewReAct(pk_eng.Raise, 100))
	g.PlayerAction(players[E].NewReAct(pk_eng.Call, 0))
	g.PlayerAction(players[A].NewReAct(pk_eng.Fold, 0))

	// River, D action first
	g.PlayerAction(players[D].NewReAct(pk_eng.Check, 100))
	g.PlayerAction(players[E].NewReAct(pk_eng.Fold, 0))

	// All players chips
	total := 0
	for _, p := range players {
		fmt.Printf("Test=Player %s chips: %d\n", p.Name(), p.Chips())
		total += p.Chips()
	}
	if total != 5000 {
		panic("Total chips should be 5000")
	}
}
