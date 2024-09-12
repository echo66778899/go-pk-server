package engine

import (
	msgpb "go-pk-server/gen"
	mylog "go-pk-server/log"
)

type Table struct {
	// State of table
	numberOfSlots int
	players       map[int]Player

	// whole table state
	minStackSize  int
	dealerPos     int
	smallBlindPos int
	bigBlindPos   int
}

func NewTable() *Table {
	return &Table{
		numberOfSlots: 0,
		players:       make(map[int]Player),
	}
}

// ============================================
// Begin of functions for TableManager interface

func (thisTable *Table) GetMaxSlot() int {
	return thisTable.numberOfSlots
}

func (thisTable *Table) UpdateMaxSeat(sn int) {
	if sn > thisTable.numberOfSlots {
		newPlayers := make(map[int]Player, sn)
		for i := 0; i < thisTable.numberOfSlots; i++ {
			newPlayers[i] = thisTable.players[i]
		}
	} else if sn < thisTable.numberOfSlots && sn >= thisTable.CountSeatedPlayers() {
		newPlayers := make(map[int]Player, sn)
		// Copy not nil players to new map
		j := 0
		for i := 0; i < thisTable.numberOfSlots; i++ {
			if thisTable.players[i] != nil {
				newPlayers[j] = thisTable.players[i]
				j++
			}
		}
		thisTable.players = newPlayers
	} else {
		// Log invalid number of slots
		mylog.Debugf("Invalid number of slots=%d\n", sn)
		return
	}
	thisTable.numberOfSlots = sn
}

func (thisTable *Table) CountSeatedPlayers() int {
	count := 0
	mylog.Debugf("Players map is %+v when count playable player\n", thisTable.players)
	if thisTable.players == nil {
		mylog.Errorf("Players map is nil\n")
		return count
	}
	for _, p := range thisTable.players {
		if p != nil {
			count++
		}
	}
	return count
}

func (thisTable *Table) CountPlayablePlayers() int {
	count := 0
	mylog.Debugf("Players map is %+v when count playable player\n", thisTable.players)
	if thisTable.players == nil {
		mylog.Errorf("Players map is nil\n")
		return count
	}
	for _, p := range thisTable.players {
		if p != nil {
			if p.Status() == msgpb.PlayerStatusType_Playing ||
				p.Status() == msgpb.PlayerStatusType_SB ||
				p.Status() == msgpb.PlayerStatusType_BB {
				count++
			}
		}
	}
	return count
}

// CheckPlayersReadiness() bool
func (thisTable *Table) CheckPlayersReadiness(s *msgpb.GameSetting) bool {
	mylog.Debugf("Check player readiness with min stack=%d\n", s.MinStackSize)
	thisTable.minStackSize = int(s.MinStackSize)
	// Check if all players are ready
	for _, p := range thisTable.players {
		if p != nil {
			// If chip >= min chip
			if p.Chips() < int(s.MinStackSize) {
				mylog.Errorf("Player %s has not enough chips\n", p.Name())
				p.UpdateStatus(msgpb.PlayerStatusType_Sat_Out)
				return false
			} else {
				if p.IsSpectating() { // If player still want spectating
					p.UpdateStatus(msgpb.PlayerStatusType_Spectating)
				} else {
					p.UpdateStatus(msgpb.PlayerStatusType_Playing)
				}
				p.UpdateCurrentBet(0)
			}
		}
	}
	return true
}

func (thisTable *Table) DetermineNextButtonPosition(fromPos int) int {
	// Invalidate all positions
	thisTable.dealerPos = -1
	thisTable.smallBlindPos = -1
	thisTable.bigBlindPos = -1

	// Find next player position to be dealer
	for i := fromPos + 1; i < thisTable.numberOfSlots+fromPos; i++ {
		slot := i % thisTable.numberOfSlots
		p := thisTable.players[slot]
		if p != nil && p.Status() == msgpb.PlayerStatusType_Playing {
			thisTable.dealerPos = slot
			break
		}
	}
	// find small blind position
	for i := thisTable.dealerPos + 1; i < thisTable.numberOfSlots+thisTable.dealerPos; i++ {
		slot := i % thisTable.numberOfSlots
		p := thisTable.players[slot]
		if p != nil && p.Status() == msgpb.PlayerStatusType_Playing {
			thisTable.smallBlindPos = slot
			p.UpdateStatus(msgpb.PlayerStatusType_SB)
			break
		}
	}
	// find big blind position
	for i := thisTable.smallBlindPos + 1; i < thisTable.numberOfSlots+thisTable.smallBlindPos; i++ {
		slot := i % thisTable.numberOfSlots
		p := thisTable.players[slot]
		if p != nil && p.Status() == msgpb.PlayerStatusType_Playing {
			thisTable.bigBlindPos = slot
			p.UpdateStatus(msgpb.PlayerStatusType_BB)
			break
		}
	}
	return thisTable.dealerPos
}

func (thisTable *Table) GetBigBlindPosition() int {
	// Find next player position to be big blind
	return thisTable.bigBlindPos
}
func (thisTable *Table) GetSmallBlindPosition() int {
	// Find next player position to be small blind
	return thisTable.smallBlindPos
}

func (thisTable *Table) DealCardsToPlayers(deck *Deck) {
	// Deal 2 cards to each player
	for i := 0; i < 2; i++ {
		for _, p := range thisTable.players {
			if p != nil {
				p.DealCard(deck.Draw(), i)
			}
		}
	}
	mylog.Info("Dealing cards to players successfully")
}

// | 0 | 1 | 2 | 3 | 4 | 5 |
// number of slots = 6, from slot = 2 , it < 6 + 2
// 2, 3, 4, 5, 0, 1
func (thisTable *Table) FindNextPlayablePlayer(fromSlot int, statusMap map[msgpb.PlayerStatusType]bool) (Player, bool) {
	// Log next from slot
	mylog.Debug("----------------------------------------------")
	mylog.Debugf("Find next player from slot=%d with status=%+v\n", fromSlot, statusMap)
	mylog.Debug("----------------------------------------------")
	defer mylog.Debug("----------------------------------------------")

	for it := fromSlot + 1; it < thisTable.numberOfSlots+fromSlot; it++ {
		slot := it % thisTable.numberOfSlots
		if thisTable.players[slot] != nil && statusMap[thisTable.players[slot].Status()] {
			mylog.Debugf("Found player %s is at slot=%d with status=%v\n",
				thisTable.players[slot].Name(), slot, thisTable.players[slot].Status())
			return thisTable.players[slot], true
		}
	}
	// Log could not find any player
	mylog.Debugf("Not found any playable player from slot=%d\n", fromSlot)
	return nil, false
}

func (thisTable *Table) PrepareNewGame() {
	for _, p := range thisTable.players {
		if p != nil {
			p.PrepareNewGame()
		}
	}
}

func (thisTable *Table) PrepareForNewRound() {
	for _, p := range thisTable.players {
		if p != nil {
			p.PrepareForNewRound()
		}
	}
}

// Can return nil, if player not found. Please check nil before using
func (thisTable *Table) GetPlayerAtPosition(position int) (p Player, ok bool) {
	// Check if posistion is valid
	if position >= 0 && position < thisTable.numberOfSlots {
		if thisTable.players[position] != nil {
			return thisTable.players[position], true
		}
	}
	mylog.Errorf("Player not found at position %d\n", position)
	return nil, false
}

func (thisTable *Table) FindLastStayingPlayer() (last Player, found bool) {
	lookFor := map[msgpb.PlayerStatusType]bool{
		msgpb.PlayerStatusType_Wait4Act: true,
		msgpb.PlayerStatusType_Playing:  true,
		msgpb.PlayerStatusType_Check:    true,
		msgpb.PlayerStatusType_Call:     true,
		msgpb.PlayerStatusType_Raise:    true,
		msgpb.PlayerStatusType_AllIn:    true,
		msgpb.PlayerStatusType_SB:       true,
		msgpb.PlayerStatusType_BB:       true,
	}
	// count all players with statuses indicating they are still playing
	count := 0
	for _, p := range thisTable.players {
		if p != nil && lookFor[p.Status()] {
			last = p
			count++
		}
	}

	if count == 1 {
		// Find the last player
		mylog.Infof("Found last player %s with status=%v\n", last.Name(), last.Status())
		return last, true
	}
	// Log could not find last player
	mylog.Warn("There is no last player")
	return nil, false
}

func (thisTable *Table) UpdatePlayerStatusDueToCurrentBetIncrease(makeActPos int) {
	lookFor := map[msgpb.PlayerStatusType]bool{
		msgpb.PlayerStatusType_Check: true,
		msgpb.PlayerStatusType_Call:  true,
		msgpb.PlayerStatusType_Raise: true,
	}

	// Update player other players
	for it := makeActPos + 1; it < thisTable.numberOfSlots+makeActPos; it++ {
		slot := it % thisTable.numberOfSlots
		if thisTable.players[slot] != nil && lookFor[thisTable.players[slot].Status()] {
			thisTable.players[slot].UpdateStatus(msgpb.PlayerStatusType_Playing)
		} else {
			if thisTable.players[slot] != nil &&
				thisTable.players[slot].Status() == msgpb.PlayerStatusType_Fold {
				thisTable.players[slot].UpdateStatus(msgpb.PlayerStatusType_Spectating)
			}
		}
	}
}

func (thisTable *Table) DoAttachedFunctionToAllPlayers(f func(Player)) {
	for _, p := range thisTable.players {
		if p != nil {
			f(p)
		}
	}
}

// End of functions for TableManager interface
// ==========================================

func (thisTable *Table) AddPlayer(reqSlot int, p Player) {
	// If requested slot is in an empty slot
	if thisTable.players[reqSlot] == nil {
		thisTable.players[reqSlot] = p
		// Log player has been added
		mylog.Debugf("Player %s has been added to slot %d. Total players: %d\n",
			p.Name(), reqSlot, thisTable.CountSeatedPlayers())
		p.UpdateStatus(msgpb.PlayerStatusType_Spectating)
	} else {
		// Log requets slot is not available
		mylog.Debugf("Slot %d is not available\n", reqSlot)
	}
}

func (thisTable *Table) RemovePlayer(reqSlot int) (remaining int) {
	if thisTable.players[reqSlot] != nil {
		// Reset player's state before removing
		thisTable.players[reqSlot].PrepareNewGame()
		thisTable.players[reqSlot].ResetPlayerState()
		thisTable.players[reqSlot] = nil
		// Log player has been removed
		remaining = thisTable.CountSeatedPlayers()
		mylog.Debugf("Player has been removed from slot %d. Total players: %d\n",
			reqSlot, remaining)
	}
	return
}
