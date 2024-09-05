package ui

import (
	"image"

	msgpb "go-pk-server/gen"
)

type PlayerUI struct {
	Block
	x0, y0 int
	Player []msgpb.Card
}

func NewPlayer() *PlayerUI {
	return &PlayerUI{
		Block: *NewBlock(),
	}
}

func (p *PlayerUI) SetTitle(title string) {
	p.Title = title
}

func (p *PlayerUI) SetCoodinate(x0, y0 int) {
	p.x0 = x0
	p.y0 = y0
}

func (p *PlayerUI) SetCards(cards []msgpb.Card) {
	p.Player = cards
	p.SetRect(p.x0, p.y0, p.x0+16, p.y0+12)
}

func (p *PlayerUI) Draw(buf *Buffer) {
	// Plot the box
	p.Block.Draw(buf)

	// plot cards
	//  ┌─────┐
	//  │.....│
	//  │.....│
	//  │.....│
	//  └─────┘
	for i := range p.Player {
		card := &p.Player[i]
		x := p.Inner.Min.X + (i * 8)
		y := p.Inner.Min.Y
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
