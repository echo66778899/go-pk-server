package main

import (
	"fmt"
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

	// Append players to a slice
	players := []*engine.OnlinePlayer{playerA, playerB, playerC, playerD, playerE}
	for _, p := range players {
		// Add chips to players
		p.AddChips(1000)
		game.PlayerJoin(p)
	}

	game.StartGame()

	// A dealer, B small blind, C big blind, D UTG and E is the last player
	// Preflop, UTG action first

	for i := 0; i < 1000; i++ {
		for _, p := range players {
			act := p.RandomSuggestionAction()
			if act.ActionType == engine.Unknown {
				continue
			}
			fmt.Printf("Player %s's chose action: [%v]\n", p.Name(), act)
			game.PlayerAction(&act)
		}
	}
}
