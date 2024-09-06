package ui

import (
	"image"
	"sync"
)

type DealerWidget struct {
	X, Y      int
	Icon      rune
	IconStyle Style

	sync.Mutex
}

func NewDealerWidget() *DealerWidget {
	return &DealerWidget{
		Icon: 'D',
		IconStyle: Style{
			Fg:       ColorBlack,
			Bg:       ColorWhite,
			Modifier: ModifierBold,
		},
	}
}

func (d *DealerWidget) GetRect() image.Rectangle {
	return image.Rect(d.X, d.Y, d.X+3, d.Y+1)
}

func (d *DealerWidget) SetRect(x, y, w, h int) {
	d.X, d.Y = x, y
}

func (d *DealerWidget) Draw(buf *Buffer) {
	buf.SetCell(NewCell(' ', d.IconStyle), image.Pt(d.X, d.Y))
	buf.SetCell(NewCell(d.Icon, d.IconStyle), image.Pt(d.X+1, d.Y))
	buf.SetCell(NewCell(' ', d.IconStyle), image.Pt(d.X+2, d.Y))
}

func (d *DealerWidget) IndexUI(dealerIdx, maxSlot int) {
	if dealerIdx > maxSlot {
		panic("index out of range")
	}

	// Convert dealer index to player index where max players is 6
	// your position is always 2, UI index is 0. So, we need to convert dealer index to UI index
	// If button is at 2, then dealer index is 0
	// If button is at 4, then dealer index is 2
	// If button is at 0, then dealer index is 4
	idx := (dealerIdx + UI_MODEL_DATA.MaxPlayers - UI_MODEL_DATA.YourTablePosition) % UI_MODEL_DATA.MaxPlayers

	refLayout := OTHER_PLAYERS[maxSlot][idx]
	if refLayout.X < TABLE_CENTER_X-10 {
		d.X = refLayout.X + 10
	} else if refLayout.X > TABLE_CENTER_X+10 {
		d.X = refLayout.X - 2
	} else {
		d.X = refLayout.X + 2
	}

	if refLayout.Y < TABLE_CENTER_Y-5 {
		d.Y = refLayout.Y + 8
	} else if refLayout.Y > TABLE_CENTER_Y+5 {
		d.Y = refLayout.Y - 2
	} else {
		d.Y = refLayout.Y
		if d.X < TABLE_CENTER_X {
			d.X += 6
		} else {
			d.X -= 4
		}
	}
}
