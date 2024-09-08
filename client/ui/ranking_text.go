package ui

import (
	msgpb "go-pk-server/gen"
	"image"
	"sync"
)

type Text struct {
	image.Rectangle
	sync.Mutex

	Text      string
	TextStyle Style
	WrapText  bool

	ptInfo *msgpb.PeerState
}

func NewText() *Text {
	return &Text{
		TextStyle: Theme.Paragraph.Text,
		WrapText:  true,
	}
}

func (t *Text) GetRect() image.Rectangle {
	return t.Rectangle
}

func (t *Text) SetRect(x, y, x2, y2 int) {
	t.Rectangle = image.Rect(x, y, x2, y2)
}

func (t *Text) Draw(buf *Buffer) {
	// If there is no peer state info, then return
	if t.ptInfo != nil {
		t.Text = t.ptInfo.GetHandRanking()
	} else {
		t.Text = ""
		return
	}

	cells := ParseStyles(t.Text, t.TextStyle)
	if t.WrapText {
		cells = WrapCells(cells, uint(t.Dx()-2))
	}

	rows := SplitCells(cells, '\n')

	for y, row := range rows {
		if y+t.Min.Y >= t.Max.Y {
			break
		}
		row = TrimCells(row, t.Dx())
		for _, cx := range BuildCellWithXArray(row) {
			x, cell := cx.X+1, cx.Cell
			buf.SetCell(cell, image.Pt(x, y).Add(t.Min))
		}
	}
}

func (t *Text) SetRefForText(pt *msgpb.PeerState) {
	t.ptInfo = pt
}

// ================================================================
// To draw the ranking display, we need to create a new struct called RankingText.
type RankingText struct {
	HandRankingTexts []*Text
	shiftStep        int
}

func NewRankingText() *RankingText {
	return &RankingText{
		HandRankingTexts: make([]*Text, 0),
	}
}

func (rt *RankingText) ClearAllTexts() {
	for _, text := range rt.HandRankingTexts {
		text.ptInfo = nil
	}
}

func (rt *RankingText) UpdateTextAtPosition(table_idx int, pb *msgpb.PeerState) {
	ui_idx := (table_idx + rt.shiftStep) % UI_MODEL_DATA.MaxPlayers
	rt.HandRankingTexts[ui_idx].SetRefForText(pb)
}

func (rt *RankingText) UpdateTextsBasedPlayers() {
	// Calculate the index of the player in the UI
	if len(UI_MODEL_DATA.Players) > 0 {
		for _, p := range UI_MODEL_DATA.Players {
			if p == nil {
				continue
			}
			if int(p.TablePosition) == UI_MODEL_DATA.YourTablePosition {
				rt.shiftStep = int(UI_MODEL_DATA.MaxPlayers - int(p.TablePosition))
				break
			}
		}
	}

	// Base on number of players, update the ranking display
	available := UI_MODEL_DATA.MaxPlayers
	if len(rt.HandRankingTexts) == available {
		return
	}

	rt.HandRankingTexts = make([]*Text, available)

	refLayout := PLAYER_LAYOUT[available]

	for i := 0; i < available; i++ {
		rt.HandRankingTexts[i] = NewText()
		rt.HandRankingTexts[i].TextStyle = Theme.Paragraph.Text
		rt.HandRankingTexts[i].WrapText = true
		// Init layout position
		rt.HandRankingTexts[i].SetRect(refLayout[i].X, refLayout[i].Y-2, refLayout[i].X+30, refLayout[i].Y)
	}
}

func (rt *RankingText) GetDisplayingTexts() []Drawable {
	texts := make([]Drawable, 0)
	for _, text := range rt.HandRankingTexts {
		if text == nil {
			continue
		}
		if text.ptInfo != nil {
			texts = append(texts, text)
		}
	}
	return texts
}
