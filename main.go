package main

import (
	engine "go-pk-server/core"
)

func main() {
	game := engine.NewGameEngine()
	game.StartEngine()

	// Create 4 players
	player1 := engine.NewOnlinePlayer("A", 123)
	player2 := engine.NewOnlinePlayer("B", 456)
	player3 := engine.NewOnlinePlayer("C", 789)
	player4 := engine.NewOnlinePlayer("D", 101)

	// Add chips to players
	player1.AddChips(1000)
	player2.AddChips(1000)
	player3.AddChips(1000)
	player4.AddChips(1000)

	game.StartGame()
	// For loop to simulate players joining the game and start the game and send actions
	game.PlayerJoin(player1)
	game.PlayerJoin(player2)
	game.PlayerJoin(player3)
	game.PlayerJoin(player4)

	game.Ready()
	game.DumpGameState()

	// Player 1 action
	game.PlayerAction(engine.NewCheckAction(player1.Position()))
	game.PlayerAction(engine.NewBetAction(player2.Position(), 10))
	game.PlayerAction(engine.NewFoldAction(player3.Position(), 0))
	game.PlayerAction(engine.NewCallAction(player4.Position()))
	game.PlayerAction(engine.NewCallAction(player1.Position()))

	game.DumpGameState()

	// Round Flop
	game.PlayerAction(engine.NewCheckAction(player1.Position()))
	game.PlayerAction(engine.NewCheckAction(player4.Position()))
	game.PlayerAction(engine.NewCheckAction(player1.Position()))
	game.PlayerAction(engine.NewCheckAction(player2.Position()))

	// Round Turn
	game.PlayerAction(engine.NewCheckAction(player4.Position()))
	game.PlayerAction(engine.NewCheckAction(player1.Position()))
	game.PlayerAction(engine.NewCheckAction(player2.Position()))
	game.PlayerAction(engine.NewCheckAction(player4.Position()))
}
