package engine

import "fmt"

// Game setting struct
type GameSetting struct {
	NumPlayers int
	StackSize  int
	SmallBlind int
	BigBlind   int
}

// Round state
type RoundState int

const (
	PreFlop RoundState = iota
	Flop
	Turn
	River
	Showdown
)

// Game logic state
type GameState struct {
	Players        []Player
	Pots           Pot
	ButtonPosition int
	CommunityCards CommunityCards

	CurrentRound     RoundState
	CurrentBet       int
	NumPlayingPlayer int

	GameSetting GameSetting
}

// Add player to the game
func (gs *GameState) AddPlayer(player Player) {
	gs.Players = append(gs.Players, player)
}

// Get player by position
func (gs *GameState) GetPlayerByPosition(position int) Player {
	for _, player := range gs.Players {
		if player.Position() == position {
			return player
		}
	}
	return nil
}

// Get player by ID
func (gs *GameState) GetPlayerByID(id int) Player {
	for _, player := range gs.Players {
		if player.ID() == id {
			return player
		}
	}
	return nil
}

// Get the player in button position
func (gs *GameState) GetButtonPlayer() Player {
	for _, player := range gs.Players {
		if player.Position() == gs.ButtonPosition {
			return player
		}
	}
	return nil
}

// Get the player in small blind position
func (gs *GameState) GetSmallBlindPlayer() Player {
	for _, player := range gs.Players {
		if player.Position() == (gs.ButtonPosition+1)%len(gs.Players) {
			return player
		}
	}
	return nil
}

// Get the player in big blind position
func (gs *GameState) GetBigBlindPlayer() Player {
	for _, player := range gs.Players {
		if player.Position() == (gs.ButtonPosition+2)%len(gs.Players) {
			return player
		}
	}
	return nil
}

// Get the next player in the game
func (gs *GameState) NextPlayer(currentPlayer Player) Player {
	for i, player := range gs.Players {
		if player.Position() == currentPlayer.Position() {
			if i == len(gs.Players)-1 {
				return gs.Players[0]
			}
			return gs.Players[i+1]
		}
	}
	return nil
}

// Get the player in the next position
func (gs *GameState) NextActivePlayer(position int) Player {
	// Find the player
	for i := 1; i <= len(gs.Players); i++ {
		position = (position + i) % len(gs.Players)
		player := gs.PlayerInPreviousPosition(position)
		if player.Status() == Active {
			return player
		}
	}
	return nil
}

// Get the previous player in the game
func (gs *GameState) PreviousPlayer(currentPlayer Player) Player {
	for i, player := range gs.Players {
		if player.Position() == currentPlayer.Position() {
			if i == 0 {
				return gs.Players[len(gs.Players)-1]
			}
			return gs.Players[i-1]
		}
	}
	return nil
}

// Get the player in the previous position
func (gs *GameState) PlayerInPreviousPosition(position int) Player {
	for _, player := range gs.Players {
		if player.Position() == position {
			return player
		}
	}
	return nil
}

// Get the player in still playing
func (gs *GameState) GetPlayingPlayers() []Player {
	players := []Player{}
	for _, player := range gs.Players {
		if player.Status() == WaitForAct {
			players = append(players, player)
		}
	}
	return players
}

// Add cards to the community cards
func (gs *GameState) AddCommunityCards(cards []Card) {
	gs.CommunityCards.Cards = append(gs.CommunityCards.Cards, cards...)
}

// Reset the community cards and the game state
func (gs *GameState) ResetBeforePlay() {
	gs.CommunityCards.Cards = []Card{}
	gs.CurrentRound = PreFlop
	gs.NumPlayingPlayer = len(gs.Players)
	gs.Pots = Pot{}

	if history := gs.Pots.History(); len(history) > 0 {
		gs.ButtonPosition++
		if gs.ButtonPosition >= gs.NumPlayingPlayer {
			gs.ButtonPosition = 0
		}
	} else {
		gs.ButtonPosition = 0
	}

	for _, player := range gs.Players {
		if player != nil {
			player.ResetForNewGame()
			fmt.Println(player.Name(), player.Chips())
		}
	}

	// Log all state
	fmt.Println("Resetting game state")
	fmt.Println("Button position:", gs.ButtonPosition)
	fmt.Println("Players:", gs.NumPlayingPlayer)
	for _, player := range gs.Players {
		fmt.Println(player.Name(), player.Chips())
	}
	fmt.Println("====================================")

}

// Reset round state
func (gs *GameState) ResetBettingState() {
	gs.CurrentBet = 0
}

// Reset the game state
func (gs *GameState) Reset() {
	gs.Players = []Player{}
	gs.Pots = Pot{}
	gs.ButtonPosition = 0
	gs.CommunityCards = CommunityCards{}
	gs.CurrentRound = PreFlop
	gs.CurrentBet = 0
	gs.NumPlayingPlayer = 0
}

func (gs *GameState) AllPlayersActed() bool {
	playersActed := 0
	for _, player := range gs.Players {
		if player.Status() != WaitForAct {
			playersActed++
		}
	}
	return playersActed == len(gs.Players)
}
