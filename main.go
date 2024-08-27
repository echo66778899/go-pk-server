package main

import (
	"time"

	engine "github.com/haiphan1811/go-pk-server/game"
)

func main() {
	game := engine.NewGameEngine()
	game.StartGame()

	// Create 4 players
	player1 := engine.NewPlayer(1, "Player 1")
	player2 := engine.NewPlayer(2, "Player 2")
	player3 := engine.NewPlayer(3, "Player 3")
	player4 := engine.NewPlayer(4, "Player 4")

	time.Sleep(1 * time.Second)
	// For loop to simulate players joining the game and start the game and send actions
	game.PlayerJoin(1, player1)
	game.PlayerJoin(1, player2)
	game.PlayerJoin(1, player3)
	game.PlayerJoin(1, player4)

	game.Ready()

	time.Sleep(1 * time.Second)
	game.PlayerAction(1, engine.PlayerAction{PlayerIdx: 0, Type: engine.Bet, Bet: 100})
	game.PlayerAction(2, engine.PlayerAction{PlayerIdx: 1, Type: engine.Call})
	game.PlayerAction(3, engine.PlayerAction{PlayerIdx: 2, Type: engine.Call})
	game.PlayerAction(4, engine.PlayerAction{PlayerIdx: 3, Type: engine.Call})

	time.Sleep(1 * time.Second)
	game.PlayerAction(1, engine.PlayerAction{PlayerIdx: 0, Type: engine.Check})
	game.PlayerAction(1, engine.PlayerAction{PlayerIdx: 1, Type: engine.Bet, Bet: 100})
	game.PlayerAction(1, engine.PlayerAction{PlayerIdx: 2, Type: engine.Raise, Bet: 100})
	game.PlayerAction(1, engine.PlayerAction{PlayerIdx: 3, Type: engine.Call})
	game.PlayerAction(1, engine.PlayerAction{PlayerIdx: 0, Type: engine.Call})
	game.PlayerAction(1, engine.PlayerAction{PlayerIdx: 1, Type: engine.Call})

	time.Sleep(2 * time.Second)
	game.PlayerAction(1, engine.PlayerAction{PlayerIdx: 0, Type: engine.Check})
	game.PlayerAction(1, engine.PlayerAction{PlayerIdx: 1, Type: engine.Bet, Bet: 100})
	game.PlayerAction(1, engine.PlayerAction{PlayerIdx: 2, Type: engine.Fold})
	game.PlayerAction(1, engine.PlayerAction{PlayerIdx: 3, Type: engine.Call})
	game.PlayerAction(1, engine.PlayerAction{PlayerIdx: 0, Type: engine.Fold})

	time.Sleep(1 * time.Second)
	game.PlayerAction(1, engine.PlayerAction{PlayerIdx: 1, Type: engine.Check})
	game.PlayerAction(1, engine.PlayerAction{PlayerIdx: 3, Type: engine.Bet, Bet: 200})
	game.PlayerAction(1, engine.PlayerAction{PlayerIdx: 1, Type: engine.Call})

	// loop to avoid exit, you can remove this loop if you want
	for {
		time.Sleep(1 * time.Second)
	}
}
