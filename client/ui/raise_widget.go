package ui

import (
	"fmt"
	"image"
)

type RaiseWidget struct {
	Block
	// The label text
	Text       string
	LabelStyle Style
	ValueStyle Style
	Value      string

	Data   int
	MaxVal int
	Labels string
}

func NewRaiseWidget() *RaiseWidget {
	return &RaiseWidget{
		Block:      *NewBlock(),
		Text:       "Amount",
		LabelStyle: NewStyle(ColorWhite),
		ValueStyle: NewStyle(ColorGreen),
		Data:       20,
	}
}

func (r *RaiseWidget) Draw(buf *Buffer) {
	r.BorderTop = false

	r.Block.Draw(buf)

	// Draw player chips
	buf.SetCell(Cell{SHADED_BLOCKS[2], NewStyle(ColorGreen)}, image.Pt(r.Inner.Min.X, r.Min.Y))

	// Draw player chips
	chips := fmt.Sprintf("%d", r.Data)
	cells := ParseStyles(chips, r.TitleStyle)
	for x, cell := range cells {
		if x+r.Inner.Min.X >= r.Inner.Max.X {
			break
		}
		buf.SetCell(cell, image.Pt(x+r.Inner.Min.X+1, r.Min.Y))
	}

	// empty cell
	buf.SetCell(Cell{' ', Theme.Default}, image.Pt(r.Min.X, r.Min.Y))
	buf.SetCell(Cell{' ', Theme.Default}, image.Pt(r.Inner.Max.X, r.Min.Y))

	// draw line
	buf.Fill(Cell{VERTICAL_LINE, r.BorderStyle}, image.Rect(r.Min.X+4, r.Min.Y+1, r.Inner.Min.X+4, r.Inner.Max.Y))
}

func (r *RaiseWidget) SetCoordinator(x, y int) {
	r.SetRect(x, y, x+9, y+10)
}

func (r *RaiseWidget) Increase(base int) {
	if r.Data+base > r.MaxVal {
		r.Data = r.MaxVal
		return
	}
	r.Data += base
}

func (r *RaiseWidget) Decrease(base int) {
	if r.Data-base < 0 {
		r.Data = 20
		return
	}
	r.Data -= base
}
