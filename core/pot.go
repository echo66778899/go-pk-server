package engine

type Pot struct {
	amount  int
	history []int
}

func (p *Pot) AddToPot(amount int) {
	p.amount += amount
}

func (p *Pot) ResetPot() {
	p.history = append(p.history, p.amount)
	p.amount = 0
}

func (p *Pot) GetPotAmount() int {
	return p.amount
}

func (p *Pot) History() []int {
	return p.history
}
