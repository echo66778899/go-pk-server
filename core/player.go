package engine

type PlayerStatus int

const (
	Active PlayerStatus = iota
	WaitForAct
	Checked
	Called
	Betted
	Raised
	Folded
	AlledIn
)

type Player interface {
	// For game engine
	UpdatePosition(int)
	DealCard(Card, int)
	UpdateCurrentBet(int)
	UpdateStatus(PlayerStatus)
	TakeChips(int)
	AddChips(int)
	ResetForNewRound()
	ResetForNewGame()
	// Tracking player state
	Position() int
	CurrentBet() int
	ShowHand() Hand
	Status() PlayerStatus
	ID() int
	Name() string
	Chips() int

	// Notifies the player of the game state
	NotifyGameState(gs GameState)
}

type OnlinePlayer struct {
	name string
	id   int
	//networkClient *wsnetwork.NetworkClient

	// Player state
	position   int
	chips      int
	hand       Hand
	status     PlayerStatus
	currentBet int
}

// NewOnlinePlayer creates a new online player.
func NewOnlinePlayer(name string, id int) *OnlinePlayer {
	return &OnlinePlayer{
		name: name,
		id:   id,
		//networkClient: networkClient,
	}
}

// UpdatePosition updates the player's position in the table.
func (p *OnlinePlayer) UpdatePosition(position int) {
	p.position = position
}

// Implement the Player interface
func (p *OnlinePlayer) DealCard(card Card, idx int) {
	p.hand.SetCard(card, idx)
}

func (p *OnlinePlayer) UpdateCurrentBet(bet int) {
	p.currentBet = bet
}

func (p *OnlinePlayer) UpdateStatus(status PlayerStatus) {
	p.status = status
}

func (p *OnlinePlayer) TakeChips(amount int) {
	p.chips -= amount
}

func (p *OnlinePlayer) AddChips(amount int) {
	p.chips += amount
}

func (p *OnlinePlayer) ResetForNewRound() {
	p.status = WaitForAct
	p.currentBet = 0
}

func (p *OnlinePlayer) ResetForNewGame() {
	p.status = Active
	p.currentBet = 0
}

func (p *OnlinePlayer) Position() int {
	return p.position
}

func (p *OnlinePlayer) CurrentBet() int {
	return p.currentBet
}

func (p *OnlinePlayer) ShowHand() Hand {
	return p.hand
}

func (p *OnlinePlayer) Status() PlayerStatus {
	return p.status
}

func (p *OnlinePlayer) ID() int {
	return p.id
}

func (p *OnlinePlayer) Name() string {
	return p.name
}

func (p *OnlinePlayer) Chips() int {
	return p.chips
}

func (p *OnlinePlayer) NotifyGameState(gs GameState) {
	//p.networkClient.NotifyGameState(gs)
}

// func (p *OnlinePlayer) NotifyAction(action ActionIf) {
// 	p.networkClient.NotifyAction(action)
// }
