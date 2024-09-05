package ui

import (
	"fmt"
	"image"

	msgpb "go-pk-server/gen"
)

type PlayerPanel struct {
	Block
	player *msgpb.Player
}

func NewPlayerPanel() *PlayerPanel {
	return &PlayerPanel{
		Block: *NewBlock(),
	}
}

func (pp *PlayerPanel) Draw(buf *Buffer) {
	chipLine := 1
	line1Offset := 2
	statusLine := 3
	isEmpty := (pp.player == nil)

	if isEmpty {
		pp.Title = "Slot Empty"
		pp.Block.Draw(buf)
		return
	}

	if !isEmpty && pp.Title == "" {
		pp.Title = pp.player.Name
	}

	// Trim title to fit in the block
	if len(pp.Title) > pp.Inner.Dx() {
		pp.Title = pp.player.Name[:pp.Inner.Dx()]
	}
	pp.Block.Draw(buf)

	// Draw cells
	buf.SetCell(Cell{VERTICAL_RIGHT, pp.BorderStyle}, image.Pt(pp.Min.X, pp.Inner.Min.Y+line1Offset))
	// buf.SetCell(Cell{VERTICAL_RIGHT, pp.BorderStyle}, image.Pt(pp.Min.X, pp.Inner.Min.Y+line2Offset))
	buf.SetCell(Cell{VERTICAL_LEFT, pp.BorderStyle}, image.Pt(pp.Inner.Max.X, pp.Inner.Min.Y+line1Offset))
	// buf.SetCell(Cell{VERTICAL_LEFT, pp.BorderStyle}, image.Pt(pp.Inner.Max.X, pp.Inner.Min.Y+line2Offset))
	buf.Fill(Cell{HORIZONTAL_LINE, pp.BorderStyle}, image.Rect(pp.Inner.Min.X, pp.Inner.Min.Y+line1Offset, pp.Inner.Max.X, pp.Inner.Min.Y+line1Offset+1))
	// buf.Fill(Cell{HORIZONTAL_LINE, pp.BorderStyle}, image.Rect(pp.Inner.Min.X, pp.Inner.Min.Y+line2Offset, pp.Inner.Max.X, pp.Inner.Min.Y+line2Offset+1))

	// Draw player chips
	buf.SetCell(Cell{SHADED_BLOCKS[2], NewStyle(ColorGreen)}, image.Pt(pp.Inner.Min.X+1, pp.Inner.Min.Y+chipLine))
	chips := fmt.Sprintf("%d", pp.player.Chips)
	cells := ParseStyles(chips, pp.TitleStyle)
	for x, cell := range cells {
		if x+pp.Inner.Min.X >= pp.Inner.Max.X {
			break
		}
		buf.SetCell(cell, image.Pt(x+pp.Inner.Min.X+3, pp.Inner.Min.Y+chipLine))
	}

	// Draw player status
	status := pp.player.Status.String()
	cells = ParseStyles(status, pp.TitleStyle)
	for x, cell := range cells {
		if x+pp.Inner.Min.X >= pp.Inner.Max.X {
			break
		}
		buf.SetCell(cell, image.Pt(x+pp.Inner.Min.X+3, pp.Inner.Min.Y+statusLine))
	}
}

func (pp *PlayerPanel) SetPlayers(player *msgpb.Player) {
	pp.player = player
}

func (pp *PlayerPanel) SetCoodinate(x, y int) {
	pp.SetRect(x, y, x+14, y+7)
}

// ========================
// PlayersWidget
type PlayersGroup struct {
	Block
	PlayersUI    []*PlayerPanel
	OtherPlayers int
	RefLayout    Layout
}

func NewPlayersGroup() *PlayersGroup {
	return &PlayersGroup{
		Block: *NewBlock(),
	}
}

func (pg *PlayersGroup) SetMaxOtherPlayers(maxOtherPlayers int) {
	if maxOtherPlayers < 2 {
		panic("Minimum number of players is 2")
	}
	if maxOtherPlayers != pg.OtherPlayers {
		pg.OtherPlayers = maxOtherPlayers
		pg.RefLayout = OTHER_PLAYERS[pg.OtherPlayers]
		if pg.RefLayout == nil {
			panic(fmt.Sprintf("No layout found for %d other players", pg.OtherPlayers))
		}
	}

	// Create player panels
	pg.PlayersUI = make([]*PlayerPanel, pg.OtherPlayers)

	for i := 0; i < pg.OtherPlayers; i++ {
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
