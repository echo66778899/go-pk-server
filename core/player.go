package engine

import (
	"fmt"
	"math/rand"

	msgpb "go-pk-server/gen"
	mylog "go-pk-server/log"
)

// Agent interface
type Agent interface {
	NotifiesChanges(message *msgpb.ServerMessage)
	DirectNotify(nameId string, message *msgpb.ServerMessage)
}

type Player interface {
	// For game engine
	UpdatePosition(int)
	DealCard(Card, int)
	HasNewCards() bool
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
	NotifyPlayerIfNewHand()
}

type OnlinePlayer struct {
	name string
	gid  uint64
	//networkClient *wsnetwork.NetworkClient

	// Player state
	position   int // slot no
	chips      int
	hand       Hand
	status     PlayerStatus
	currentBet int

	// internal state
	isNewCard bool

	// Suggested actions
	suggestAction []PlayerActType

	// Connection agent
	connAgent Agent
}

// NewOnlinePlayer creates a new online player.
func NewOnlinePlayer(name string, connAgent Agent, gid uint64) *OnlinePlayer {
	return &OnlinePlayer{
		name:      name,
		gid:       gid,
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
	if idx == 1 {
		p.isNewCard = true
	}
}

// HasNewCards returns true if the player has new cards. Otherwise, it returns false.
// It also resets the isNewCard flag to false.
func (p *OnlinePlayer) HasNewCards() bool {
	ret := p.isNewCard
	p.isNewCard = false
	return ret
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
	// If player is sat out, set status to playing
	if p.status == SatOut {
		p.status = Playing
	}
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
	p.hand.Reset()
	p.isNewCard = false
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
	return p.gid
}

func (p *OnlinePlayer) Name() string {
	return p.name
}

func (p *OnlinePlayer) Chips() int {
	return p.chips
}

func (p *OnlinePlayer) NotifyPlayerIfNewHand() {
	mylog.Infof("Sync player %s's new cards\n", p.name)
	// Send the player's hand to the player
	handMsg := msgpb.ServerMessage{
		Message: &msgpb.ServerMessage_PeerState{
			PeerState: &msgpb.PeerState{
				TablePos:    int32(p.position),
				PlayerCards: make([]*msgpb.Card, 0),
			},
		},
	}
	// Add the player's hand to the message
	for _, card := range p.hand.Cards() {
		handMsg.GetPeerState().PlayerCards = append(handMsg.GetPeerState().PlayerCards, &msgpb.Card{
			Suit: msgpb.SuitType(card.Suit),
			Rank: msgpb.RankType(card.Value),
		})
	}
	p.connAgent.DirectNotify(p.name, &handMsg)
}

func (p *OnlinePlayer) RandomSuggestionAction() PlayerAction {
	// check if nil list
	if len(p.suggestAction) == 0 || p.chips == 0 || p.status == Folded || p.status != Wait4Act {
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
