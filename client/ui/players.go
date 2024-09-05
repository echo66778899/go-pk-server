package ui

import (
	"fmt"
	"sync"

	msgpb "go-pk-server/gen"
)

type MainPlayer struct {
	CurrentPlayerPossition int
	MainCard               []msgpb.Card
	sync.Mutex
}

var CurrentPlayer = MainPlayer{CurrentPlayerPossition: 0}

// ========================
// PlayersWidget
type PlayersGroup struct {
	Block
	PlayersUI    []*PlayerPanel
	TotalPlayers int
	ShiftStep    int
	RefLayout    Layout
}

func NewPlayersGroup() *PlayersGroup {
	return &PlayersGroup{
		Block: *NewBlock(),
	}
}

func (pg *PlayersGroup) UpdateState(players []*msgpb.Player, force bool) {
	if force {
		for _, p := range pg.PlayersUI {
			p.SetPlayers(nil)
		}
		for _, p := range players {
			if p == nil {
				continue
			}
			if int(p.TablePosition) == CurrentPlayer.CurrentPlayerPossition {
				pg.ShiftStep = int(pg.TotalPlayers - int(p.TablePosition))
				break
			}
		}
	}

	for i, p := range players {
		if i >= pg.TotalPlayers {
			break
		}
		// Calculate the index of the player in the UI
		table_idx := int(p.TablePosition)
		ui_idx := (table_idx + pg.ShiftStep) % pg.TotalPlayers // 1 + 5 % 6 = 0
		pg.PlayersUI[ui_idx].SetPlayers(p)
	}
}

func (pg *PlayersGroup) SetMaxOtherPlayers(maxOtherPlayers int) {
	if maxOtherPlayers < 2 {
		panic("Minimum number of players is 2")
	}
	if maxOtherPlayers != pg.TotalPlayers {
		pg.TotalPlayers = maxOtherPlayers
		pg.RefLayout = OTHER_PLAYERS[pg.TotalPlayers]
		if pg.RefLayout == nil {
			panic(fmt.Sprintf("No layout found for %d other players", pg.TotalPlayers))
		}
	} else {
		return
	}

	// Create player panels
	pg.PlayersUI = make([]*PlayerPanel, pg.TotalPlayers)

	for i := 0; i < pg.TotalPlayers; i++ {
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
