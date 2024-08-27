package engine

type Pot struct {
	amount int
}

func (p *Pot) AddToPot(amount int) {
	p.amount += amount
}

func (p *Pot) ResetPot() {
	p.amount = 0
}

func (p *Pot) GetPotAmount() int {
	return p.amount
}
