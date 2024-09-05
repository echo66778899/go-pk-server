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

func (d *DealerWidget) IndexUI(idx, maxSlot int) {
	if idx > maxSlot {
		panic("index out of range")
	}

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
