package ui

import (
	"fmt"
	msgpb "go-pk-server/gen"
	"image"
	"sync"
)

type CallChipText struct {
	image.Rectangle
	sync.Mutex

	Text      string
	TextStyle Style
	WrapText  bool

	IsVisible bool
}

func NewCallChipText() *CallChipText {
	return &CallChipText{
		TextStyle: Theme.Paragraph.Text,
		WrapText:  true,
		IsVisible: true,
	}
}

func (t *CallChipText) GetRect() image.Rectangle {
	return t.Rectangle
}

func (t *CallChipText) SetRect(x, y, x2, y2 int) {
	t.Rectangle = image.Rect(x, y, x2, y2)
}

func (t *CallChipText) Draw(buf *Buffer) {
	if !t.IsVisible {
		return
	}

	if UI_MODEL_DATA.YourPlayerState != nil && UI_MODEL_DATA.CurrentBet > 0 {
		if UI_MODEL_DATA.YourPlayerState.GetCurrentBet() < int32(UI_MODEL_DATA.CurrentBet) {
			toCallText := int32(UI_MODEL_DATA.CurrentBet) - UI_MODEL_DATA.YourPlayerState.GetCurrentBet()

			if toCallText > UI_MODEL_DATA.YourPlayerState.GetChips() {
				toCallText = UI_MODEL_DATA.YourPlayerState.GetChips()
			}

			t.Text = fmt.Sprintf("%d\n", toCallText)
		} else {
			t.Text = ""
		}
	}

	cells := ParseStyles(t.Text, t.TextStyle)
	if t.WrapText {
		cells = WrapCells(cells, uint(t.Dx()-2))
	}

	rows := SplitCells(cells, '\n')

	buf.SetCell(Cell{SHADED_BLOCKS[2], NewStyle(ColorWhite)}, image.Pt(t.Min.X, t.Min.Y+1))

	for y, row := range rows {
		if y+t.Min.Y >= t.Max.Y {
			break
		}
		row = TrimCells(row, t.Dx())
		for _, cx := range BuildCellWithXArray(row) {
			x, cell := cx.X+2, cx.Cell
			buf.SetCell(cell, image.Pt(x, y).Add(t.Min))
		}
	}
}

func (t *CallChipText) SetRefForText(pt *msgpb.PeerState) {
	t.IsVisible = true
}
