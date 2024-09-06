package ui

import (
	"fmt"
	"image"

	msgpb "go-pk-server/gen"
)

type PlayerPanel struct {
	Block
	player *msgpb.Player
	ppInfo *msgpb.PeerState

	Slot int
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

	if pp.player == nil {
		emptyStyle := NewStyle(ColorDarkGray, ColorBlack, ModifierBold)
		pp.Title = "Slot Empty"
		pp.TitleStyle = emptyStyle
		pp.BorderStyle = emptyStyle
		pp.Block.Draw(buf)

		slotNo := fmt.Sprintf("[ %d ]\n", pp.Slot)
		// Draw cells at center indicating empty slot able to join
		cells := ParseStyles(slotNo, emptyStyle)
		for x, cell := range cells {
			if x+pp.Inner.Min.X+4 >= pp.Inner.Max.X {
				break
			}
			buf.SetCell(cell, image.Pt(x+pp.Inner.Min.X+4, pp.Inner.Min.Y+line1Offset))
		}

		return
	}

	if pp.player != nil {
		pp.Title = pp.player.Name
	}

	// Trim title to fit in the block
	if len(pp.Title) > pp.Inner.Dx() {
		pp.Title = pp.player.Name[:pp.Inner.Dx()]
	}
	pp.Block.Draw(buf)

	// Draw cells
	buf.SetCell(Cell{VERTICAL_RIGHT, pp.BorderStyle}, image.Pt(pp.Min.X, pp.Inner.Min.Y+line1Offset))
	buf.SetCell(Cell{VERTICAL_LEFT, pp.BorderStyle}, image.Pt(pp.Inner.Max.X, pp.Inner.Min.Y+line1Offset))
	buf.Fill(Cell{HORIZONTAL_LINE, pp.BorderStyle}, image.Rect(pp.Inner.Min.X, pp.Inner.Min.Y+line1Offset, pp.Inner.Max.X, pp.Inner.Min.Y+line1Offset+1))

	// Draw player chips
	buf.SetCell(Cell{SHADED_BLOCKS[2], NewStyle(ColorGreen)}, image.Pt(pp.Inner.Min.X+2, pp.Inner.Min.Y+chipLine))
	chips := fmt.Sprintf("%d", pp.player.Chips)
	cells := ParseStyles(chips, pp.TitleStyle)
	for x, cell := range cells {
		if x+pp.Inner.Min.X+4 >= pp.Inner.Max.X {
			break
		}
		buf.SetCell(cell, image.Pt(x+pp.Inner.Min.X+4, pp.Inner.Min.Y+chipLine))
	}

	// Draw player status
	status := pp.player.Status
	cells = ParseStyles(status, pp.TitleStyle)
	for x, cell := range cells {
		if x+pp.Inner.Min.X+3 >= pp.Inner.Max.X {
			break
		}
		buf.SetCell(cell, image.Pt(x+pp.Inner.Min.X+3, pp.Inner.Min.Y+statusLine))
	}

	if pp.ppInfo != nil {
		// Draw pocket pair
		for i, card := range pp.ppInfo.PlayerCards {
			x := pp.Inner.Min.X + (i * 7) + 14
			y := pp.Inner.Min.Y
			style := cardStyle[card.GetSuit()]
			buf.SetCell(Cell{TOP_LEFT, style}, image.Pt(x, y))
			buf.Fill(Cell{HORIZONTAL_LINE, style}, image.Rect(x+1, y, x+6, y+1))
			buf.SetCell(Cell{TOP_RIGHT, style}, image.Pt(x+6, y))
			buf.Fill(Cell{VERTICAL_LINE, style}, image.Rect(x, y+1, x+1, y+4))

			buf.SetCell(Cell{ranksIcon[card.GetRank()], style}, image.Pt(x+1, y+1))
			buf.SetCell(Cell{suitsIcon[card.GetSuit()], style}, image.Pt(x+3, y+2))
			buf.SetCell(Cell{ranksIcon[card.GetRank()], style}, image.Pt(x+5, y+3))

			buf.SetCell(Cell{BOTTOM_LEFT, style}, image.Pt(x, y+4))
			buf.Fill(Cell{HORIZONTAL_LINE, style}, image.Rect(x+1, y+4, x+6, y+5))
			buf.SetCell(Cell{BOTTOM_RIGHT, style}, image.Pt(x+6, y+4))
			buf.Fill(Cell{VERTICAL_LINE, style}, image.Rect(x+6, y+1, x+7, y+4))
		}
	}
}

func (pp *PlayerPanel) SetPlayers(player *msgpb.Player) {
	if player == nil {
		pp.ppInfo = nil
	}
	pp.player = player
}

func (pp *PlayerPanel) SetSlot(slot int) {
	pp.Slot = slot
}

func (pp *PlayerPanel) SetPocketPair(pb *msgpb.PeerState) {
	pp.ppInfo = pb
	pp.SetRect(pp.Min.X, pp.Min.Y, pp.Min.X+14, pp.Min.Y+7)
}

func (pp *PlayerPanel) SetCoodinate(x, y int) {
	pp.SetRect(x, y, x+14, y+7)
}

// Overide the SetRect method of the Block
func (pp *PlayerPanel) GetRect() image.Rectangle {
	if pp.ppInfo != nil {
		return image.Rect(pp.Min.X, pp.Min.Y, pp.Min.X+14+16, pp.Min.Y+7)
	} else {
		return image.Rect(pp.Min.X, pp.Min.Y, pp.Min.X+14, pp.Min.Y+7)
	}
}
