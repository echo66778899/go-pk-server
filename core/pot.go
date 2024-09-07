package engine

import (
	mylog "go-pk-server/log"
	"sort"
)

// SidePot struct represents each side pot.
type SidePot struct {
	Amount  int
	Players []int // Players who joined the pot
}

// Pot represents the pot in a Poker game.
type Pot struct {
	amount       int
	playerAmount map[int]int
}

func NewPot() Pot {
	return Pot{
		amount:       0,
		playerAmount: make(map[int]int),
	}
}

func (p *Pot) AddToPot(position, amount int) {
	// Log the amount added to the current pot and the pot before the addition and after the addition
	mylog.Infof("Player idx %d adding %d to the pot. Pot before: %d, Pot after: %d\n",
		position, amount, p.amount, p.amount+amount)
	p.amount += amount
	p.playerAmount[position] += amount
}

func (p *Pot) ResetPot() {
	// Log the pot before and after reset
	mylog.Infof("Resetting the pot. Pot before: %d, Pot after: %d\n", p.amount, 0)
	p.amount = 0
	p.playerAmount = make(map[int]int)
}

func (p *Pot) Total() int {
	return p.amount
}

func (p *Pot) Size() int {
	return p.amount
}

func (p *Pot) PlayerAmount(position int) int {
	return p.playerAmount[position]
}

func (p *Pot) CalculateSidePots() []SidePot {
	// Get slice of player from map
	players := make([]struct {
		TablePos int
		Bet      int
	}, 0)

	for pos, bet := range p.playerAmount {
		players = append(players, struct {
			TablePos int
			Bet      int
		}{TablePos: pos, Bet: bet})
	}

	// Sort players by their bet amounts (smallest to largest)
	sort.Slice(players, func(i, j int) bool {
		return players[i].Bet < players[j].Bet
	})

	var sidePots []SidePot
	totalPlayers := len(players)

	// Keep track of the total chips bet in each round
	previousBet := 0

	for _, player := range players {
		// Calculate how much this player contributes to the current side pot
		currentBet := player.Bet

		// If this player's bet is higher than the previous one, we calculate a side pot
		if currentBet > previousBet {
			// Calculate the size of this side pot
			sidePotAmount := (currentBet - previousBet) * totalPlayers

			// Add this side pot to the list
			potPlayers := make([]int, 0)
			for _, p := range players[:totalPlayers] { // Only active players in the current side pot
				potPlayers = append(potPlayers, p.TablePos)
			}

			sidePots = append(sidePots, SidePot{
				Amount:  sidePotAmount,
				Players: potPlayers,
			})

			// Update the previous bet
			previousBet = currentBet
		}

		// Reduce the number of players for the next side pot
		totalPlayers--
	}

	return sidePots
}
