package engine

import (
	"log"
)

type TableManager struct {
	numberOfSlots int
	players       map[int]Player
}

func NewTableManager() *TableManager {
	return &TableManager{
		numberOfSlots: 0,
		players:       make(map[int]Player),
	}
}

func (tm *TableManager) UpdateMaxNoOfSlot(sn int) {
	if sn > tm.numberOfSlots {
		newPlayers := make(map[int]Player, sn)
		for i := 0; i < tm.numberOfSlots; i++ {
			newPlayers[i] = tm.players[i]
		}
	} else if sn < tm.numberOfSlots && sn >= tm.GetNumberOfPlayers() {
		newPlayers := make(map[int]Player, sn)
		// Copy not nil players to new map
		j := 0
		for i := 0; i < tm.numberOfSlots; i++ {
			if tm.players[i] != nil {
				newPlayers[j] = tm.players[i]
				j++
			}
		}
		tm.players = newPlayers
	} else {
		// Log invalid number of slots
		log.Printf("Invalid number of slots=%d\n", sn)
		return

	}
	tm.numberOfSlots = sn
}

func (tm *TableManager) GetMaxNoSlot() int {
	return tm.numberOfSlots
}

func (tm *TableManager) AddPlayer(reqSlot int, p Player) {
	// If requested slot is in an empty slot
	if tm.players[reqSlot] == nil {
		tm.players[reqSlot] = p
		// Log player has been added
		log.Printf("Player %s has been added to slot %d. Total players: %d\n",
			p.Name(), reqSlot, tm.GetNumberOfPlayers())
	} else {
		// Log requets slot is not available
		log.Printf("Slot %d is not available\n", reqSlot)
	}
}

func (tm *TableManager) RemovePlayer(reqSlot int) {
	if tm.players[reqSlot] != nil {
		tm.players[reqSlot] = nil
		// Log player has been removed
		log.Printf("Player has been removed from slot %d. Total players: %d\n",
			reqSlot, tm.GetNumberOfPlayers())
	}
}

func (tm *TableManager) GetNumberOfPlayingPlayers() int {
	count := 0
	for _, p := range tm.players {
		if p != nil && p.Status() == Playing {
			log.Printf("Player %s status: %v, chips: %d\n", p.Name(), p.Status(), p.Chips())
			count++
		}
	}
	return count
}

func (tm *TableManager) GetNumberOfPlayers() int {
	count := 0
	for _, p := range tm.players {
		if p != nil {
			count++
		}
	}
	return count
}

func (tm *TableManager) GetPlayer(reqSlot int) Player {
	// Check if requested slot is valid < number of slots
	if reqSlot >= tm.numberOfSlots || reqSlot < 0 {
		panic("Invalid slot number")
	}

	if tm.players[reqSlot] != nil {
		return tm.players[reqSlot]
	}
	// Log player not found
	log.Printf("Player not found at slot=%d\n", reqSlot)
	return nil
}

// | 0 | 1 | 2 | 3 | 4 | 5 |
// number of slots = 6, from slot = 2 , it < 8
// 2, 3, 4, 5, 0, 1
func (tm *TableManager) NextPlayer(fromSlot int, status PlayerStatus) Player {
	// Log next from slot
	log.Println("----------------------------------------------")
	log.Printf("Find next player from slot=%d with status=%v\n", fromSlot, status)
	log.Println("----------------------------------------------")
	defer log.Println("----------------------------------------------")

	fromSlot += 1

	for it := fromSlot; it < tm.numberOfSlots+fromSlot; it++ {
		slot := it % tm.numberOfSlots
		if tm.players[slot] != nil {
			log.Printf("Found player %s is at slot=%d with status=%v\n",
				tm.players[slot].Name(), slot, tm.players[slot].Status())
			if tm.players[slot].Status() == status {
				return tm.players[slot]
			}
		}
	}
	// Log could not find any player
	log.Printf("Could not find any player at slot=%d with status=%v\n", fromSlot, status)
	return nil
}

func (tm *TableManager) GetListOfOtherPlayers(exceptSlot int, expect ...PlayerStatus) []Player {
	otherPlayers := make([]Player, 0)
	for slot, p := range tm.players {
		if p != nil && slot != exceptSlot {
			for _, s := range expect {
				if p.Status() == s {
					otherPlayers = append(otherPlayers, p)
				}
			}
		}
	}
	return otherPlayers
}

func (tm *TableManager) GetListOfPlayers(expect ...PlayerStatus) []Player {
	players := make([]Player, 0)
	for _, p := range tm.players {
		if p != nil {
			for _, s := range expect {
				if p.Status() == s {
					players = append(players, p)
				}
			}
		}
	}
	return players
}

func (tm *TableManager) IsAllPlayersActed() bool {
	// Check if all players have acted
	for _, p := range tm.players {
		if p != nil && p.Status() == Wait4Act || p.Status() == Playing {
			return false
		}
	}
	return true
}

func (tm *TableManager) IsAllOthersFold(ex int) bool {
	// Check if all other players have folded
	for _, p := range tm.players {
		if p != nil && p.Position() != ex && p.Status() != Folded {
			return false
		}
	}
	return true
}

func (tm *TableManager) GetOnlyOnePlayingPlayer() Player {
	// Get only one playing player
	var playingPlayer Player
	for _, p := range tm.players {
		if p != nil && p.Status() == Playing {
			if playingPlayer != nil {
				return nil
			}
			playingPlayer = p
		}
	}
	return playingPlayer
}

func (tm *TableManager) ResetForNewGame() {
	for _, p := range tm.players {
		if p != nil {
			p.ResetForNewGame()
		}
	}
}

func (tm *TableManager) ResetForNewRound() {
	for _, p := range tm.players {
		if p != nil {
			p.ResetForNewRound()
		}
	}
}
