package ui

import "image"

type TextBox struct {
	Block
	Text      string
	TextStyle Style
	WrapText  bool
}

func NewParagraph() *TextBox {
	return &TextBox{
		Block:     *NewBlock(),
		TextStyle: Theme.Paragraph.Text,
		WrapText:  true,
	}
}

func (tb *TextBox) Draw(buf *Buffer) {
	tb.BorderStyle = NewStyle(ColorDarkGray, ColorBlack, ModifierBold)
	tb.Block.Draw(buf)

	cells := ParseStyles(tb.Text, tb.TextStyle)
	if tb.WrapText {
		cells = WrapCells(cells, uint(tb.Inner.Dx()))
	}

	rows := SplitCells(cells, '\n')

	for y, row := range rows {
		if y+tb.Inner.Min.Y >= tb.Inner.Max.Y {
			break
		}
		row = TrimCells(row, tb.Inner.Dx())
		for _, cx := range BuildCellWithXArray(row) {
			x, cell := cx.X+1, cx.Cell
			buf.SetCell(cell, image.Pt(x, y).Add(tb.Inner.Min))
		}
	}
}
