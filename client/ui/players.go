package ui

import (
	"fmt"

	msgpb "go-pk-server/gen"
)

// ========================
// PlayersWidget
type PlayersGroup struct {
	Block
	PlayersUI  []*PlayerPanel
	ShiftStep  int
	RefLayout  Layout
	MaxPlayers int
}

func NewPlayersGroup() *PlayersGroup {
	return &PlayersGroup{
		Block: *NewBlock(),
	}
}

func (pg *PlayersGroup) UpdateState(force bool) {
	players := UI_MODEL_DATA.Players

	if force {
		for _, p := range pg.PlayersUI {
			p.SetPlayers(nil)
		}
		for _, p := range players {
			if p == nil {
				continue
			}
			if int(p.TablePosition) == UI_MODEL_DATA.YourTablePosition {
				pg.ShiftStep = int(UI_MODEL_DATA.MaxPlayers - int(p.TablePosition))
				break
			}
		}
	}

	for i, p := range players {
		if p == nil {
			continue
		}
		if i >= UI_MODEL_DATA.MaxPlayers {
			break
		}
		// Calculate the index of the player in the UI
		table_idx := int(p.TablePosition)
		ui_idx := (table_idx + pg.ShiftStep) % UI_MODEL_DATA.MaxPlayers // 1 + 5 % 6 = 0
		pg.PlayersUI[ui_idx].SetPlayers(p)
	}

	for i, p := range pg.PlayersUI {
		if p.player == nil {
			// Set logic index for empty slot
			p.SetSlot((i + UI_MODEL_DATA.YourTablePosition) % UI_MODEL_DATA.MaxPlayers)
		}
	}
}

func (pg *PlayersGroup) UpdateGroupPlayers(maxOtherPlayers int) {
	if maxOtherPlayers < 2 || pg.MaxPlayers == maxOtherPlayers {
		return
	}
	pg.MaxPlayers = maxOtherPlayers

	pg.RefLayout = OTHER_PLAYERS[UI_MODEL_DATA.MaxPlayers]
	if pg.RefLayout == nil {
		panic(fmt.Sprintf("No layout found for %d other players", UI_MODEL_DATA.MaxPlayers))
	}

	// Create player panels
	pg.PlayersUI = make([]*PlayerPanel, UI_MODEL_DATA.MaxPlayers)

	for i := 0; i < UI_MODEL_DATA.MaxPlayers; i++ {
		pg.PlayersUI[i] = NewPlayerPanel()
		pg.PlayersUI[i].SetCoodinate(pg.RefLayout[i].X, pg.RefLayout[i].Y)
	}
}

func (pg *PlayersGroup) GetAllItems() []Drawable {
	items := make([]Drawable, 0)
	for _, p := range pg.PlayersUI {
		items = append(items, p)
	}
	return items
}

func (pg *PlayersGroup) UpdatePocketPair(pb *msgpb.PeerState) {
	pp := pg.PlayersUI[0]
	if pp == nil {
		return
	}

	pp.SetPocketPair(pb)
}
