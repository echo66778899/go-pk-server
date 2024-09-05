package ui

import "image"

type ButtonType int

const (
	BNT_FoldButton  ButtonType = 0
	BNT_CheckButton ButtonType = 1
	BNT_CallButton  ButtonType = 2
	BNT_RaiseButton ButtonType = 3
	BNT_AllInButton ButtonType = 4

	BNT_JoinTableButton   ButtonType = 5
	BNT_StartGameButton   ButtonType = 6
	BNT_LeaveGameButton   ButtonType = 7
	BNT_RequestChipButton ButtonType = 8
)

type Button struct {
	Block
	Text         string
	ActiveStyle  Style
	DisableStyle Style
	DefaultStyle Style

	// Controller is the controller that handles the button's actions.
	Type       ButtonType
	isSelected bool
	isDisable  bool
}

func NewButton(text string, t ButtonType) *Button {
	return &Button{
		Block:        *NewBlock(),
		Text:         "  " + text + "  ",
		Type:         t,
		ActiveStyle:  Style{ColorBlack, ColorWhite, ModifierBold},
		DisableStyle: Style{ColorDarkGray, ColorBlack, ModifierBold},
		DefaultStyle: Style{ColorWhite, ColorBlack, ModifierBold},
	}
}

func (b *Button) Draw(buf *Buffer) {
	style := b.DefaultStyle

	if b.isDisable {
		style = b.DisableStyle
	} else if b.isSelected {
		style = b.ActiveStyle
	}
	b.BorderStyle = style

	// Draw the border and title
	b.Block.Draw(buf)
	cells := ParseStyles(b.Text, style)

	for x, cx := range BuildCellWithXArray(cells) {
		buf.SetCell(cx.Cell, image.Pt(x, 0).Add(b.Inner.Min))
	}
}

func (b *Button) SetCenterCoordinator(x, y int) {
	x -= len(b.Text) / 2
	y -= 1
	b.SetRect(x, y, x+len(b.Text)+2, y+3)
}

func (b *Button) Select() {
	b.isSelected = true
}

func (b *Button) IsSelected() bool {
	return b.isSelected
}

func (b *Button) Unselect() {
	b.isSelected = false
}

func (b *Button) Enable(e bool) {
	b.isDisable = !e
	if !e {
		b.isSelected = false
	}
}

func (b *Button) IsEnabled() bool {
	return !b.isDisable
}
