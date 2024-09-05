package engine

import (
	"fmt"
	"go-pk-server/msg"
	"math/rand"
)

// Agent interface
type Agent interface {
	NotifiesChanges(gId uint64, message *msg.CommunicationMessage)
}

type Player interface {
	// For game engine
	UpdatePosition(int)
	DealCard(Card, int)
	UpdateCurrentBet(int)
	UpdateStatus(PlayerStatus)
	TakeChips(int)
	AddChips(int)
	UpdateSuggestions([]PlayerActType)
	ResetForNewRound()
	ResetForNewGame()
	// Tracking player state
	Position() int
	CurrentBet() int
	ShowHand() *Hand
	Status() PlayerStatus
	ID() uint64
	Name() string
	Chips() int

	// Notifies the player of the game state
	NotifyGameState(gs *GameState, tm *TableManager)
}

type OnlinePlayer struct {
	name string
	id   uint64
	//networkClient *wsnetwork.NetworkClient

	// Player state
	position   int // slot no
	chips      int
	hand       Hand
	status     PlayerStatus
	currentBet int

	// Suggested actions
	suggestAction []PlayerActType

	// Connection agent
	connAgent Agent
}

// NewOnlinePlayer creates a new online player.
func NewOnlinePlayer(name string, connAgent Agent, id uint64) *OnlinePlayer {
	return &OnlinePlayer{
		name:      name,
		id:        id,
		connAgent: connAgent,
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
	fmt.Printf("Player %s's status is updated from %v to %v\n", p.name, p.status, status)
	p.status = status
}

func (p *OnlinePlayer) TakeChips(amount int) {
	p.chips -= amount
	if p.chips < 0 {
		panic("Player counld not has negative chips")
	}
	fmt.Printf("Player %s's remaining chips: %d\n", p.name, p.chips)
}

func (p *OnlinePlayer) AddChips(amount int) {
	p.chips += amount
	fmt.Printf("Player %s's chips are added by %d. Total chips: %d\n", p.name, amount, p.chips)
}

func (p *OnlinePlayer) UpdateSuggestions(suggestions []PlayerActType) {
	p.suggestAction = suggestions
}

func (p *OnlinePlayer) ResetForNewRound() {
	if (p.status != Folded) && (p.status != AlledIn) {
		fmt.Printf("Player %s is reset for new round\n", p.name)
		p.status = Playing
	}
	p.currentBet = 0
}

func (p *OnlinePlayer) ResetForNewGame() {
	// Print player name
	fmt.Printf("Player %s is reset for new game\n", p.name)
	if p.chips > 0 {
		p.status = Playing
	} else {
		p.status = SatOut
	}
	p.currentBet = 0
}

func (p *OnlinePlayer) Position() int {
	return p.position
}

func (p *OnlinePlayer) CurrentBet() int {
	return p.currentBet
}

func (p *OnlinePlayer) ShowHand() *Hand {
	return &p.hand
}

func (p *OnlinePlayer) Status() PlayerStatus {
	return p.status
}

func (p *OnlinePlayer) ID() uint64 {
	return p.id
}

func (p *OnlinePlayer) Name() string {
	return p.name
}

func (p *OnlinePlayer) Chips() int {
	return p.chips
}

func (p *OnlinePlayer) NotifyGameState(gs *GameState, tm *TableManager) {
	var players []msg.PlayerState
	for _, player := range tm.players {
		players = append(players, msg.PlayerState{
			Name:   player.Name(),
			Slot:   player.Position(),
			Chips:  player.Chips(),
			Bet:    player.CurrentBet(),
			Status: player.Status().String(),
		})
	}

	var communityCards []msg.Card
	for _, card := range gs.cc.Cards {
		communityCards = append(communityCards, msg.Card{
			Suit:  int(card.Suit),
			Value: int(card.Value),
		})
	}

	var playerHand []msg.Card
	for _, card := range p.hand.Cards() {
		playerHand = append(playerHand, msg.Card{
			Suit:  int(card.Suit),
			Value: int(card.Value),
		})
	}

	var message = msg.CommunicationMessage{
		Type: msg.SyncGameStateMsgType,
		Payload: msg.SyncGameStateMessage{
			CommunityCards: communityCards,
			Players:        players,
		},
	}

	fmt.Printf("Player %s is notified of the game state: %v\n", p.name, message)

	p.connAgent.NotifiesChanges(p.id, &message)
}

func (p *OnlinePlayer) RandomSuggestionAction() PlayerAction {
	// check if nil list
	if len(p.suggestAction) == 0 || p.chips == 0 || p.status == Folded || p.status != WaitForAct {
		return PlayerAction{
			PlayerPosition: p.position,
			ActionType:     Unknown,
		}
	}
	// Log random suggestion action from the list
	len := len(p.suggestAction)
	fmt.Printf("Player %s's suggested %d actions: %v\n", p.name, len, p.suggestAction)

	// Randomly select an action from suggested actions list
	index := rand.Intn(len - 1)
	action := p.suggestAction[index]

	if action == Raise {
		// Randomly select a raise amount between big blind and max chips
		raiseAmount := rand.Intn(p.chips / 2)
		// mod by 20
		raiseAmount = raiseAmount - (raiseAmount % 20)
		if raiseAmount < 20 {
			raiseAmount = 20
		}
		return PlayerAction{
			PlayerPosition: p.position,
			ActionType:     action,
			Amount:         raiseAmount,
		}
	}

	return PlayerAction{
		PlayerPosition: p.position,
		ActionType:     action,
	}
}

func (p *OnlinePlayer) NewReAct(ActionType PlayerActType, Amount int) ActionIf {
	return &PlayerAction{
		PlayerPosition: p.position,
		ActionType:     ActionType,
		Amount:         Amount,
	}
}
