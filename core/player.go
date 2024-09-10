package engine

import (
	msgpb "go-pk-server/gen"
	mylog "go-pk-server/log"
)

// Agent interface
type Agent interface {
	NotifiesChanges(message *msgpb.ServerMessage)
	DirectNotify(nameId string, message *msgpb.ServerMessage)
}

type Player interface {
	// For game control
	UpdatePosition(int)
	DealCard(*msgpb.Card, int)
	HasPocketCards() bool
	DropPocketCards()
	UpdateCurrentBet(int)
	UpdateStatus(msgpb.PlayerStatusType)
	GetChipForBet(int)
	AddWonChips(int)
	TakeChips(int)
	AddChips(int)
	UpdateInvalidAction([]msgpb.PlayerGameActionType)
	PrepareForNewRound()
	PrepareNewGame()
	ResetPlayerState()
	IsSpectating() bool

	// Tracking player state
	Position() int
	IsDealer() bool
	CurrentBet() int
	ShowHand() *Hand
	Status() msgpb.PlayerStatusType
	RoomID() uint64
	Name() string
	Chips() int
	ChipChange() int
	UnsupportAction() []msgpb.PlayerGameActionType

	// Notifies the player of the game state
	NotifyPlayerIfNewHand()
}

type OnlinePlayer struct {
	// Required fields
	name      string
	gid       uint64
	connAgent Agent

	// Player option control
	isSpectating bool

	// Player game state
	position int // slot no
	chips    int

	// All rounds state
	isDealer         bool
	hand             Hand
	chipChangeAmount int // chips won or lost in the current round
	// Eeach round state
	currentBet int
	invalidPGA []msgpb.PlayerGameActionType
	// player status
	status msgpb.PlayerStatusType
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
func (p *OnlinePlayer) DealCard(card *msgpb.Card, idx int) {
	p.hand.SetCard(card, idx)
}

// HasPocketCards returns true if the player has new cards.
func (p *OnlinePlayer) HasPocketCards() bool {
	return p.hand.HasCards()
}

func (p *OnlinePlayer) DropPocketCards() {
	p.hand.Reset()
}

func (p *OnlinePlayer) UpdateCurrentBet(bet int) {
	p.currentBet = bet
}

func (p *OnlinePlayer) UpdateStatus(status msgpb.PlayerStatusType) {
	mylog.Debugf("Player %s's status is updated from %v to %v\n", p.name, p.status, status)
	p.status = status
}

func (p *OnlinePlayer) AddWonChips(amount int) {
	p.chips += amount
	p.chipChangeAmount += amount
	mylog.Debugf("Player %s's won +%d chips. Total chips: %d\n", p.name, amount, p.chips)
}

func (p *OnlinePlayer) ChipChange() int {
	return p.chipChangeAmount
}

func (p *OnlinePlayer) GetChipForBet(amount int) {
	p.chips -= amount
	p.chipChangeAmount -= amount
	if p.chips < 0 {
		panic("Player counld not has negative chips")
	}
	mylog.Debugf("Player %s's remaining chips after betting: %d\n", p.name, p.chips)
}

func (p *OnlinePlayer) TakeChips(amount int) {
	p.chips -= amount
	if p.chips < 0 {
		panic("Player counld not has negative chips")
	}
	mylog.Debugf("Player %s's remaining chips: %d\n", p.name, p.chips)
}

func (p *OnlinePlayer) AddChips(amount int) {
	p.chips += amount
	mylog.Debugf("Player %s's chips are added by %d. Total chips: %d\n", p.name, amount, p.chips)
	// If player is sat out, set status to playing
	if p.status == msgpb.PlayerStatusType_Sat_Out {
		p.status = msgpb.PlayerStatusType_Playing
	}
}

func (p *OnlinePlayer) UpdateInvalidAction(invalidActions []msgpb.PlayerGameActionType) {
	p.invalidPGA = invalidActions
}

func (p *OnlinePlayer) PrepareForNewRound() {
	// each round state
	p.currentBet = 0
	p.invalidPGA = nil
	// player status
	previousStatus := p.status
	switch p.status {
	case msgpb.PlayerStatusType_Fold:
		// Player is folded, reset the status to Spectating for new round
		p.status = msgpb.PlayerStatusType_Spectating
	case msgpb.PlayerStatusType_Wait4Act,
		msgpb.PlayerStatusType_Check,
		msgpb.PlayerStatusType_Call,
		msgpb.PlayerStatusType_Raise:
		// Player is playing, reset the status to Playing for new round
		p.status = msgpb.PlayerStatusType_Playing
	case msgpb.PlayerStatusType_AllIn:
		// Player is all-in, Let they show status
	}
	mylog.Debugf("Player %s has status %s, reset to %s for new round\n",
		p.name, previousStatus, p.status)
}

func (p *OnlinePlayer) PrepareNewGame() {
	// Print player name
	mylog.Debugf("Player [%s] state is reset for new game\n", p.name)
	// all round state
	p.isDealer = false
	p.hand.Reset()
	p.chipChangeAmount = 0
	// each round state
	p.currentBet = 0
	p.invalidPGA = nil
	// player status
	if (p.status != msgpb.PlayerStatusType_Sat_Out) &&
		(p.status != msgpb.PlayerStatusType_Spectating) {
		mylog.Debugf("Player %s's status is reset from %s to %s for new game\n",
			p.name, p.status, msgpb.PlayerStatusType_Playing)
		p.status = msgpb.PlayerStatusType_Playing
	}
}

// ResetPlayerState resets the player's state as new.
func (p *OnlinePlayer) ResetPlayerState() {
	// Print player name
	mylog.Debugf("Player %s's state is reset as new\n", p.name)
	// player state
	p.position = 0
	p.chips = 0
	p.status = msgpb.PlayerStatusType_Sat_Out
}

func (p *OnlinePlayer) IsSpectating() bool {
	return p.isSpectating
}

func (p *OnlinePlayer) Position() int {
	return p.position
}

func (p *OnlinePlayer) IsDealer() bool {
	return p.isDealer
}

func (p *OnlinePlayer) CurrentBet() int {
	return p.currentBet
}

func (p *OnlinePlayer) ShowHand() *Hand {
	return &p.hand
}

func (p *OnlinePlayer) Status() msgpb.PlayerStatusType {
	return p.status
}

func (p *OnlinePlayer) RoomID() uint64 {
	return p.gid
}

func (p *OnlinePlayer) Name() string {
	return p.name
}

func (p *OnlinePlayer) Chips() int {
	return p.chips
}

func (p *OnlinePlayer) UnsupportAction() []msgpb.PlayerGameActionType {
	return p.invalidPGA
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
		handMsg.GetPeerState().PlayerCards = append(handMsg.GetPeerState().PlayerCards, card)
	}
	p.connAgent.DirectNotify(p.name, &handMsg)
}

func (p *OnlinePlayer) NewReAct(ActionType msgpb.PlayerGameActionType, Amount int) ActionIf {
	return &PlayerAction{
		PlayerPosition: p.position,
		ActionType:     ActionType,
		Amount:         Amount,
	}
}
