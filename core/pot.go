package engine

import "fmt"

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
	fmt.Printf("Adding %d to the pot. Pot before: %d, Pot after: %d\n",
		amount, p.amount, p.amount+amount)
	p.amount += amount
	p.playerAmount[position] += amount
}

func (p *Pot) ResetPot() {
	// Log the pot before and after reset
	fmt.Printf("Resetting the pot. Pot before: %d, Pot after: %d\n", p.amount, 0)
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
