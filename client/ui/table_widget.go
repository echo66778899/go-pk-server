package ui

import (
	"image"
	"math"
	"sync"
)

type Table struct {
	centerX, centerY int
	radiusX, radiusY int
	Style            Style
	sync.Mutex
}

func NewTable() *Table {
	return &Table{
		Style: Style{
			Fg: ColorWhite,
			Bg: ColorBlack,
		},
		radiusX: 55,
		radiusY: 17,
	}
}

func (t *Table) SetRect(x, y, width, height int) {
	t.centerX = x + width/2
	t.centerY = y + height/2
	t.radiusX = width / 2
	t.radiusY = height / 2
}

func (t *Table) GetRect() image.Rectangle {
	return image.Rect(t.centerX-t.radiusX, t.centerY-t.radiusY, t.centerX+t.radiusX, t.centerY+t.radiusY)
}

func (t *Table) Draw(buf *Buffer) {

	for angle := 0.0; angle < 2*math.Pi; angle += 0.01 {
		// Parametric equation for an ellipse
		x := int(float64(t.radiusX) * math.Cos(angle))
		y := int(float64(t.radiusY) * math.Sin(angle))

		// Plot the point on the ellipse outline
		buf.SetCell(NewCell(DOT, t.Style), image.Pt(t.centerX+x, t.centerY+y))
	}
}
