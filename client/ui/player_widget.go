package ui

import (
	"fmt"
	"image"

	msgpb "go-pk-server/gen"
)

type PlayerPanel struct {
	Block
	player *msgpb.PlayerState
	ppInfo *msgpb.PeerState

	InfoTextStyle Style
	CurBetStyle   Style
	Slot          int
	CurBet        int
}

func NewPlayerPanel() *PlayerPanel {
	return &PlayerPanel{
		Block:         *NewBlock(),
		InfoTextStyle: NewStyle(ColorWhite, ColorBlack, ModifierBold),
		CurBetStyle:   NewStyle(ColorWhite, ColorBlack, ModifierBold),
	}
}

func (pp *PlayerPanel) Draw(buf *Buffer) {
	chipLine := 1
	line1Offset := 2
	statusLine := 3

	if pp.player == nil {
		emptyStyle := NewStyle(ColorDarkGray, ColorBlack, ModifierBold)
		pp.Title = "Empty Slot"
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

	status := ""
	if pp.player != nil {
		playerName := pp.player.Name
		// Add space to title to center it
		if len(playerName) > 8 {
			playerName = TrimString(playerName, 8)
		}
		pp.Title = " " + playerName + " "
		// Update style based on player status
		switch pp.player.Status {
		case msgpb.PlayerStatusType_Wait4Act:
			pp.TitleStyle = NewStyle(ColorGreen, ColorBlack, ModifierBold)
			pp.BorderStyle = NewStyle(ColorGreen, ColorBlack, ModifierBold)
			pp.InfoTextStyle = NewStyle(ColorWhite, ColorBlack, ModifierBold)
			pp.CurBetStyle = NewStyle(ColorWhite, ColorBlack, ModifierBold)
			status = ""
		case msgpb.PlayerStatusType_Fold:
			pp.TitleStyle = NewStyle(ColorDarkGray, ColorBlack, ModifierBold)
			pp.BorderStyle = NewStyle(ColorDarkGray, ColorBlack, ModifierBold)
			pp.InfoTextStyle = NewStyle(ColorDarkGray, ColorBlack, ModifierBold)
			status = msgpb.PlayerStatusType_name[int32(msgpb.PlayerStatusType_Fold)]
			pp.CurBet = 0
		case msgpb.PlayerStatusType_Playing:
			pp.TitleStyle = NewStyle(ColorWhite, ColorBlack, ModifierBold)
			pp.BorderStyle = NewStyle(ColorWhite, ColorBlack, ModifierBold)
			pp.InfoTextStyle = NewStyle(ColorWhite, ColorBlack, ModifierBold)
			pp.CurBetStyle = NewStyle(ColorWhite, ColorBlack, ModifierBold)
			status = "--||--"
		case msgpb.PlayerStatusType_Unplayable:
			pp.TitleStyle = NewStyle(ColorDarkGray, ColorBlack, ModifierBold)
			pp.BorderStyle = NewStyle(ColorDarkGray, ColorBlack, ModifierBold)
			pp.InfoTextStyle = NewStyle(ColorBlack, ColorBlack, ModifierBold)
			status = msgpb.PlayerStatusType_name[int32(pp.player.Status)]
			pp.CurBet = 0
		case msgpb.PlayerStatusType_LOSER:
			pp.TitleStyle = NewStyle(ColorWhite, ColorBlack, ModifierBold)
			pp.BorderStyle = NewStyle(ColorWhite, ColorBlack, ModifierBold)
			pp.InfoTextStyle = NewStyle(ColorLightRed, ColorBlack, ModifierBold)
			pp.CurBetStyle = NewStyle(ColorLightRed, ColorBlack, ModifierBold)
			status = msgpb.PlayerStatusType_name[int32(pp.player.Status)]
			pp.CurBet = int(pp.player.ChangeAmount) * -1
		case msgpb.PlayerStatusType_WINNER:
			pp.TitleStyle = NewStyle(ColorWhite, ColorBlack, ModifierBold)
			pp.BorderStyle = NewStyle(ColorWhite, ColorBlack, ModifierBold)
			pp.InfoTextStyle = NewStyle(ColorLightGreen, ColorBlack, ModifierBold)
			pp.CurBetStyle = NewStyle(ColorLightGreen, ColorBlack, ModifierBold)
			status = msgpb.PlayerStatusType_name[int32(pp.player.Status)]
			pp.CurBet = int(pp.player.ChangeAmount)
		default:
			pp.TitleStyle = NewStyle(ColorWhite, ColorBlack, ModifierBold)
			pp.BorderStyle = NewStyle(ColorWhite, ColorBlack, ModifierBold)
			pp.InfoTextStyle = NewStyle(ColorWhite, ColorBlack, ModifierBold)
			pp.CurBetStyle = NewStyle(ColorWhite, ColorBlack, ModifierBold)
			status = msgpb.PlayerStatusType_name[int32(pp.player.Status)]
		}
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
	buf.SetCell(Cell{SHADED_BLOCKS[2], NewStyle(ColorGreen)}, image.Pt(pp.Inner.Min.X+3, pp.Inner.Min.Y+chipLine))
	chips := fmt.Sprintf("%d", pp.player.Chips)
	cells := ParseStyles(chips, pp.InfoTextStyle)
	for x, cell := range cells {
		if x+pp.Inner.Min.X+4 >= pp.Inner.Max.X {
			break
		}
		buf.SetCell(cell, image.Pt(x+pp.Inner.Min.X+4, pp.Inner.Min.Y+chipLine))
	}

	// Draw player status
	cells = ParseStyles(status, pp.InfoTextStyle)
	// Center the status text
	statusLen := len(status)
	xStart := (pp.Inner.Dx() - statusLen) / 2
	for x, cell := range cells {
		if x+pp.Inner.Min.X+xStart >= pp.Inner.Max.X {
			break
		}
		buf.SetCell(cell, image.Pt(x+pp.Inner.Min.X+xStart, pp.Inner.Min.Y+statusLine))
	}

	// Draw the current bet amount > 0
	if pp.CurBet > 0 || pp.CurBet < 0 ||
		status == msgpb.PlayerStatusType_name[int32(msgpb.PlayerStatusType_WINNER)] {
		if status == msgpb.PlayerStatusType_name[int32(msgpb.PlayerStatusType_WINNER)] {
			// Draw plus sign
			buf.SetCell(Cell{COLLAPSED, NewStyle(ColorLightGreen, ColorBlack, ModifierBold)},
				image.Pt(pp.Inner.Min.X+2, pp.Inner.Min.Y+statusLine+1))
			buf.SetCell(Cell{SHADED_BLOCKS[2], NewStyle(ColorGreen)}, image.Pt(pp.Inner.Min.X+4, pp.Inner.Min.Y+statusLine+1))
		} else if status == msgpb.PlayerStatusType_name[int32(msgpb.PlayerStatusType_LOSER)] {
			// Draw minus sign
			buf.SetCell(Cell{EXPANDED, NewStyle(ColorLightRed, ColorBlack, ModifierBold)},
				image.Pt(pp.Inner.Min.X+2, pp.Inner.Min.Y+statusLine+1))
			buf.SetCell(Cell{SHADED_BLOCKS[2], NewStyle(ColorRed)}, image.Pt(pp.Inner.Min.X+4, pp.Inner.Min.Y+statusLine+1))
		}
		// Draw cells chip icon
		// Draw current bet amount
		curBetChars := fmt.Sprintf("%d\n", pp.CurBet)
		cells = ParseStyles(curBetChars, pp.CurBetStyle)
		for x, cell := range cells {
			if x+pp.Inner.Min.X+5 >= pp.Inner.Max.X {
				break
			}
			buf.SetCell(cell, image.Pt(x+pp.Inner.Min.X+5, pp.Inner.Min.Y+statusLine+1))
		}
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

func (pp *PlayerPanel) SetPlayers(player *msgpb.PlayerState) {
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

func (pp *PlayerPanel) SetCurBet(bet int) {
	pp.CurBet = bet
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
