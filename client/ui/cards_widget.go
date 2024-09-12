package ui

import (
	"fmt"
	"image"

	msgpb "go-pk-server/gen"
)

type Gauge struct {
	Block
	Percent    int
	BarColor   Color
	Label      string
	LabelStyle Style
}

func NewGauge() *Gauge {
	return &Gauge{
		Block:      *NewBlock(),
		BarColor:   Theme.Gauge.Bar,
		LabelStyle: Theme.Gauge.Label,
	}
}

func (g *Gauge) Draw(buf *Buffer) {
	g.Block.Draw(buf)

	label := g.Label
	if label == "" {
		label = fmt.Sprintf("%d%%", g.Percent)
	}

	// plot bar
	barWidth := int((float64(g.Percent) / 100) * float64(g.Inner.Dx()))
	buf.Fill(
		NewCell(' ', NewStyle(ColorClear, g.BarColor)),
		image.Rect(g.Inner.Min.X, g.Inner.Min.Y, g.Inner.Min.X+barWidth, g.Inner.Max.Y),
	)

	// plot label
	labelXCoordinate := g.Inner.Min.X + (g.Inner.Dx() / 2) - int(float64(len(label))/2)
	labelYCoordinate := g.Inner.Min.Y + ((g.Inner.Dy() - 1) / 2)
	if labelYCoordinate < g.Inner.Max.Y {
		for i, char := range label {
			style := g.LabelStyle
			if labelXCoordinate+i+1 <= g.Inner.Min.X+barWidth {
				style = NewStyle(g.BarColor, ColorClear, ModifierReverse)
			}
			buf.SetCell(NewCell(char, style), image.Pt(labelXCoordinate+i, labelYCoordinate))
		}
	}
}

type Cards struct {
	Block
	x0, y0 int
	Cards  []*msgpb.Card
}

func NewCards() *Cards {
	return &Cards{
		Block: *NewBlock(),
	}
}

func (c *Cards) SetTitle(title string) {
	c.Title = title
}

func (c *Cards) SetCoodinate(x0, y0 int) {
	c.x0 = x0
	c.y0 = y0
}

func (c *Cards) SetCards(cards []*msgpb.Card) {
	c.Cards = cards
	//y := c.x0 + len(c.Cards)*8
	x1 := c.x0 + 5*8
	c.SetRect(c.x0, c.y0, x1+3, c.y0+7)
}

func (c *Cards) Draw(buf *Buffer) {
	// Plot the box
	c.BorderStyle = NewStyle(ColorDarkGray, ColorBlack, ModifierBold)
	c.Block.Draw(buf)

	// plot cards
	//  ┌─────┐
	//  │.....│
	//  │.....│
	//  │.....│
	//  └─────┘
	for i := range c.Cards {
		card := c.Cards[i]
		x := c.Inner.Min.X + 1 + (i * 8)
		y := c.Inner.Min.Y
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
