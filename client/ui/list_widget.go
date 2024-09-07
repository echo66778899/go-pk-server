package ui

import (
	"image"

	rw "github.com/mattn/go-runewidth"
)

type List struct {
	Block
	Rows             []string
	WrapText         bool
	TextStyle        Style
	SelectedRow      int
	topRow           int
	SelectedRowStyle Style
}

func NewList() *List {
	return &List{
		Block:            *NewBlock(),
		TextStyle:        Theme.List.Text,
		SelectedRowStyle: Theme.List.Text,
	}
}

func (l *List) Draw(buf *Buffer) {
	l.BorderStyle = NewStyle(ColorDarkGray, ColorBlack, ModifierBold)
	l.Block.Draw(buf)

	// Custom title
	buf.SetString(
		" Player Balance ",
		l.TitleStyle,
		image.Pt(l.Min.X+7, l.Min.Y),
	)

	point := l.Inner.Min

	// adjusts view into widget
	if l.SelectedRow >= l.Inner.Dy()+l.topRow {
		l.topRow = l.SelectedRow - l.Inner.Dy() + 1
	} else if l.SelectedRow < l.topRow {
		l.topRow = l.SelectedRow
	}

	// draw rows
	for row := l.topRow; row < len(l.Rows) && point.Y < l.Inner.Max.Y; row++ {
		cells := ParseStyles(l.Rows[row], l.TextStyle)
		if l.WrapText {
			cells = WrapCells(cells, uint(l.Inner.Dx()))
		}
		for j := 0; j < len(cells) && point.Y < l.Inner.Max.Y; j++ {
			style := cells[j].Style
			if row == l.SelectedRow {
				style = l.SelectedRowStyle
			}
			if cells[j].Rune == '\n' {
				point = image.Pt(l.Inner.Min.X, point.Y+1)
			} else {
				if point.X+1 == l.Inner.Max.X+1 && len(cells) > l.Inner.Dx() {
					buf.SetCell(NewCell(ELLIPSES, style), point.Add(image.Pt(-1, 0)))
					break
				} else {
					buf.SetCell(NewCell(cells[j].Rune, style), point)
					point = point.Add(image.Pt(rw.RuneWidth(cells[j].Rune), 0))
				}
			}
		}
		point = image.Pt(l.Inner.Min.X, point.Y+1)
	}

	// draw UP_ARROW if needed
	if l.topRow > 0 {
		buf.SetCell(
			NewCell(UP_ARROW, NewStyle(ColorWhite)),
			image.Pt(l.Inner.Max.X-1, l.Inner.Min.Y),
		)
	}

	// draw DOWN_ARROW if needed
	if len(l.Rows) > int(l.topRow)+l.Inner.Dy() {
		buf.SetCell(
			NewCell(DOWN_ARROW, NewStyle(ColorWhite)),
			image.Pt(l.Inner.Max.X-1, l.Inner.Max.Y-1),
		)
	}
}

// ScrollAmount scrolls by amount given. If amount is < 0, then scroll up.
// There is no need to set l.topRow, as this will be set automatically when drawn,
// since if the selected item is off screen then the topRow variable will change accordingly.
func (l *List) ScrollAmount(amount int) {
	if len(l.Rows)-int(l.SelectedRow) <= amount {
		l.SelectedRow = len(l.Rows) - 1
	} else if int(l.SelectedRow)+amount < 0 {
		l.SelectedRow = 0
	} else {
		l.SelectedRow += amount
	}
}

func (l *List) ScrollUp() {
	l.ScrollAmount(-1)
}

func (l *List) ScrollDown() {
	l.ScrollAmount(1)
}

func (l *List) ScrollPageUp() {
	// If an item is selected below top row, then go to the top row.
	if l.SelectedRow > l.topRow {
		l.SelectedRow = l.topRow
	} else {
		l.ScrollAmount(-l.Inner.Dy())
	}
}

func (l *List) ScrollPageDown() {
	l.ScrollAmount(l.Inner.Dy())
}

func (l *List) ScrollHalfPageUp() {
	l.ScrollAmount(-int(FloorFloat64(float64(l.Inner.Dy()) / 2)))
}

func (l *List) ScrollHalfPageDown() {
	l.ScrollAmount(int(FloorFloat64(float64(l.Inner.Dy()) / 2)))
}

func (l *List) ScrollTop() {
	l.SelectedRow = 0
}

func (l *List) ScrollBottom() {
	l.SelectedRow = len(l.Rows) - 1
}
