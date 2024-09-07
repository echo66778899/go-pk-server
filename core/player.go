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
	// For game engine
	UpdatePosition(int)
	DealCard(Card, int)
	HasNewCards() bool
	UpdateCurrentBet(int)
	UpdateStatus(msgpb.PlayerStatusType)
	AddWonChips(int)
	GetChipForBet(int)
	TakeChips(int)
	AddChips(int)
	UpdateInvalidAction([]msgpb.PlayerGameActionType)
	ResetForNewRound()
	ResetForNewGame()
	ResetPlayerState()
	// Tracking player state
	Position() int
	CurrentBet() int
	ShowHand() *Hand
	Status() msgpb.PlayerStatusType
	ID() uint64
	Name() string
	Chips() int
	ChipChange() int
	UnsupportActs() []msgpb.PlayerGameActionType

	// Notifies the player of the game state
	NotifyPlayerIfNewHand()
}

type OnlinePlayer struct {
	// Required fields
	name      string
	gid       uint64
	connAgent Agent

	// Player state
	position int // slot no
	chips    int
	status   msgpb.PlayerStatusType
	// Round state
	hand             Hand
	chipChangeAmount int // chips won or lost in the current round
	currentBet       int
	isNewCard        bool
	// Invalid actions for the player
	invalidAction []msgpb.PlayerGameActionType
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
	if p.status == msgpb.PlayerStatusType_NoStack {
		p.status = msgpb.PlayerStatusType_Playing
	}
}

func (p *OnlinePlayer) UpdateInvalidAction(invalidActions []msgpb.PlayerGameActionType) {
	p.invalidAction = invalidActions
}

func (p *OnlinePlayer) ResetForNewRound() {
	if (p.status != msgpb.PlayerStatusType_Fold) && (p.status != msgpb.PlayerStatusType_AllIn) {
		mylog.Debugf("Player %s is reset for new round\n", p.name)
		p.status = msgpb.PlayerStatusType_Playing
	}
	p.currentBet = 0
	p.invalidAction = nil
}

func (p *OnlinePlayer) ResetForNewGame() {
	// Print player name
	mylog.Debugf("Player %s is reset for new game\n", p.name)
	// round state
	p.hand.Reset()
	p.chipChangeAmount = 0
	p.currentBet = 0
	p.isNewCard = false
	// invalid actions for ui
	p.invalidAction = nil
}

// ResetPlayerState resets the player's state as new.
func (p *OnlinePlayer) ResetPlayerState() {
	// Print player name
	mylog.Debugf("Player %s's state is reset as new\n", p.name)
	// player state
	p.position = 0
	p.chips = 0
	p.status = msgpb.PlayerStatusType_NoStack
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

func (p *OnlinePlayer) Status() msgpb.PlayerStatusType {
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

func (p *OnlinePlayer) UnsupportActs() []msgpb.PlayerGameActionType {
	return p.invalidAction
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

func (p *OnlinePlayer) NewReAct(ActionType msgpb.PlayerGameActionType, Amount int) ActionIf {
	return &PlayerAction{
		PlayerPosition: p.position,
		ActionType:     ActionType,
		Amount:         Amount,
	}
}

// func (p *OnlinePlayer) RandomSuggestionAction() PlayerAction {
// // check if nil list
// if len(p.suggestAction) == 0 || p.chips == 0 || p.status == msgpb.PlayerStatusType_Fold || p.status != msgpb.PlayerStatusType_Wait4Act {
// 	return PlayerAction{
// 		PlayerPosition: p.position,
// 		ActionType:     Unknown,
// 	}
// }
// // Log random suggestion action from the list
// len := len(p.suggestAction)
// fmt.Printf("Player %s's suggested %d actions: %v\n", p.name, len, p.suggestAction)

// // Randomly select an action from suggested actions list
// index := rand.Intn(len - 1)
// action := p.suggestAction[index]

// if action == Raise {
// 	// Randomly select a raise amount between big blind and max chips
// 	raiseAmount := rand.Intn(p.chips / 2)
// 	// mod by 20
// 	raiseAmount = raiseAmount - (raiseAmount % 20)
// 	if raiseAmount < 20 {
// 		raiseAmount = 20
// 	}
// 	return PlayerAction{
// 		PlayerPosition: p.position,
// 		ActionType:     action,
// 		Amount:         raiseAmount,
// 	}
// }

// return PlayerAction{
// 	PlayerPosition: p.position,
// 	ActionType:     action,
// }
// }
