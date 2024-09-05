package ui

import (
	"image"

	rw "github.com/mattn/go-runewidth"
)

// Cell represents a viewable terminal cell
type Cell struct {
	Rune  rune
	Style Style
}

var CellClear = Cell{
	Rune:  ' ',
	Style: StyleClear,
}

// NewCell takes 1 to 2 arguments
// 1st argument = rune
// 2nd argument = optional style
func NewCell(rune rune, args ...interface{}) Cell {
	style := StyleClear
	if len(args) == 1 {
		style = args[0].(Style)
	}
	return Cell{
		Rune:  rune,
		Style: style,
	}
}

// Buffer represents a section of a terminal and is a renderable rectangle of cells.
type Buffer struct {
	image.Rectangle
	CellMap map[image.Point]Cell
}

func NewBuffer(r image.Rectangle) *Buffer {
	buf := &Buffer{
		Rectangle: r,
		CellMap:   make(map[image.Point]Cell),
	}
	buf.Fill(CellClear, r) // clears out area
	return buf
}

func (bf *Buffer) GetCell(p image.Point) Cell {
	return bf.CellMap[p]
}

func (bf *Buffer) SetCell(c Cell, p image.Point) {
	bf.CellMap[p] = c
}

func (bf *Buffer) Fill(c Cell, rect image.Rectangle) {
	for x := rect.Min.X; x < rect.Max.X; x++ {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			bf.SetCell(c, image.Pt(x, y))
		}
	}
}

func (bf *Buffer) SetString(s string, style Style, p image.Point) {
	x := 0
	for _, char := range s {
		bf.SetCell(Cell{char, style}, image.Pt(p.X+x, p.Y))
		x += rw.RuneWidth(char)
	}
}
