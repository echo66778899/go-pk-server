package main

import (
	engine "go-pk-server/core"
)

func main() {
	game := engine.NewGameEngine()
	game.StartEngine()

	// Create 4 players
	playerA := engine.NewOnlinePlayer("A", 123, 0)
	playerB := engine.NewOnlinePlayer("B", 456, 1)
	playerC := engine.NewOnlinePlayer("C", 789, 3)
	playerD := engine.NewOnlinePlayer("D", 101, 5)
	playerE := engine.NewOnlinePlayer("E", 121, 6)

	// Add chips to players
	playerA.AddChips(1000)
	playerB.AddChips(1000)
	playerC.AddChips(1000)
	playerD.AddChips(1000)
	playerE.AddChips(1000)

	game.StartGame()
	// For loop to simulate players joining the game and start the game and send actions
	game.PlayerJoin(playerA)
	game.PlayerJoin(playerE)
	game.PlayerJoin(playerC)
	game.PlayerJoin(playerB)
	game.PlayerJoin(playerD)

	game.Ready()

	// A dealer, B small blind, C big blind, D UTG and E is the last player
	// Preflop, UTG action first
	game.PlayerAction(engine.NewCallAction(playerD.Position()))
	game.PlayerAction(engine.NewCallAction(playerE.Position()))
	game.PlayerAction(engine.NewCallAction(playerA.Position()))
	game.PlayerAction(engine.NewCallAction(playerB.Position()))
	game.PlayerAction(engine.NewFoldAction(playerC.Position()))

	// Round Flop
	game.PlayerAction(engine.NewCheckAction(playerB.Position()))
	game.PlayerAction(engine.NewCheckAction(playerD.Position()))
	game.PlayerAction(engine.NewFoldAction(playerE.Position()))
	game.PlayerAction(engine.NewCheckAction(playerA.Position()))

	// Round Turn
	game.PlayerAction(engine.NewRaiseAction(playerB.Position(), 100))
	game.PlayerAction(engine.NewCallAction(playerD.Position()))
	game.PlayerAction(engine.NewRaiseAction(playerA.Position(), 200))
	game.PlayerAction(engine.NewCallAction(playerB.Position()))
	game.PlayerAction(engine.NewFoldAction(playerD.Position()))

	// Round River
	game.PlayerAction(engine.NewCheckAction(playerB.Position()))
	game.PlayerAction(engine.NewAllInAction(playerA.Position()))
	game.PlayerAction(engine.NewCallAction(playerB.Position()))
}
