package engine

import "fmt"

type PlayerManager struct {
	numberOfSlots int
	players       map[int]Player
}

func NewPlayerManager(maxSlots int) *PlayerManager {
	return &PlayerManager{
		numberOfSlots: maxSlots,
		players:       make(map[int]Player),
	}
}

func (pm *PlayerManager) UpdateMaxNoOfSlot(sn int) {
	if sn > pm.numberOfSlots {
		newPlayers := make(map[int]Player, sn)
		for i := 0; i < pm.numberOfSlots; i++ {
			newPlayers[i] = pm.players[i]
		}
	} else if sn < pm.numberOfSlots && sn >= pm.GetNumberOfPlayers() {
		newPlayers := make(map[int]Player, sn)
		// Copy not nil players to new map
		j := 0
		for i := 0; i < pm.numberOfSlots; i++ {
			if pm.players[i] != nil {
				newPlayers[j] = pm.players[i]
				j++
			}
		}
		pm.players = newPlayers
	} else {
		// Log invalid number of slots
		fmt.Printf("Invalid number of slots=%d\n", sn)
		return

	}
	pm.numberOfSlots = sn
}

func (pm *PlayerManager) GetMaxNoSlot() int {
	return pm.numberOfSlots
}

func (pm *PlayerManager) AddPlayer(reqSlot int, p Player) {
	// If requested slot is in an empty slot
	if pm.players[reqSlot] == nil {
		pm.players[reqSlot] = p
		// Log player has been added
		fmt.Printf("Player %s has been added to slot %d. Total players: %d\n",
			p.Name(), reqSlot, pm.GetNumberOfPlayers())
	} else {
		// Log requets slot is not available
		fmt.Printf("Slot %d is not available\n", reqSlot)
	}
}

func (pm *PlayerManager) RemovePlayer(reqSlot int) {
	if pm.players[reqSlot] != nil {
		pm.players[reqSlot] = nil
		// Log player has been removed
		fmt.Printf("Player has been removed from slot %d. Total players: %d\n",
			reqSlot, pm.GetNumberOfPlayers())
	}
}

func (pm *PlayerManager) GetNumberOfPlayers() int {
	count := 0
	for _, p := range pm.players {
		if p != nil {
			count++
		}
	}
	return count
}

func (pm *PlayerManager) GetPlayer(reqSlot int) Player {
	if pm.players[reqSlot] != nil {
		return pm.players[reqSlot]
	}
	// Log player not found
	fmt.Printf("Player not found at slot=%d\n", reqSlot)
	return nil
}

// | 0 | 1 | 2 | 3 | 4 | 5 |
// number of slots = 6, from slot = 2 , it < 8
// 2, 3, 4, 5, 0, 1
func (pm *PlayerManager) NextPlayer(fromSlot int, status PlayerStatus) Player {
	// Log next from slot
	fmt.Printf("Find next player from slot=%d with status=%v\n", fromSlot, status)
	fromSlot += 1

	for it := fromSlot; it < pm.numberOfSlots+fromSlot; it++ {
		slot := it % pm.numberOfSlots
		if pm.players[slot] != nil {
			fmt.Printf("Found player %s is at slot=%d with status=%v\n",
				pm.players[slot].Name(), slot, pm.players[slot].Status())
			if pm.players[slot].Status() == status {
				return pm.players[slot]
			}
		}
	}
	// Log could not find any player
	fmt.Printf("Could not find any player at slot=%d with status=%v\n", fromSlot, status)
	return nil
}

func (pm *PlayerManager) GetListOfOtherPlayers(exceptSlot int, expect ...PlayerStatus) []Player {
	result := make([]Player, 0)
	for slot, p := range pm.players {
		if p != nil && slot != exceptSlot {
			for _, s := range expect {
				if p.Status() == s {
					result = append(result, p)
				}
			}
		}
	}
	return result
}

func (pm *PlayerManager) IsAllPlayersActed() bool {
	// Check if all players have acted
	for _, p := range pm.players {
		if p != nil && p.Status() == WaitForAct || p.Status() == Active {
			return false
		}
	}
	return true
}

func (pm *PlayerManager) ResetForNewGame() {
	for _, p := range pm.players {
		if p != nil {
			p.ResetForNewGame()
		}
	}
}

func (pm *PlayerManager) ResetForNewRound() {
	for _, p := range pm.players {
		if p != nil {
			p.ResetForNewRound()
		}
	}
}

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
	ShowHand() *Hand
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
	position   int // slot no
	chips      int
	hand       Hand
	status     PlayerStatus
	currentBet int

	// Suggested actions
	suggestAction []ActionType
}

// NewOnlinePlayer creates a new online player.
func NewOnlinePlayer(name string, id, slot int) *OnlinePlayer {
	return &OnlinePlayer{
		name:     name,
		id:       id,
		position: slot,
		chips:    0,
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
}

func (p *OnlinePlayer) AddChips(amount int) {
	p.chips += amount
}

func (p *OnlinePlayer) ResetForNewRound() {
	if (p.status != Folded) && (p.status != AlledIn) {
		fmt.Printf("Player %s is reset for new round\n", p.name)
		p.status = Active
	}
	p.currentBet = 0
}

func (p *OnlinePlayer) ResetForNewGame() {
	// Print player name
	fmt.Printf("Player %s is reset for new game\n", p.name)
	p.status = Active
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
