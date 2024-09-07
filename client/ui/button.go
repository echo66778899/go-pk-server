package ui

import "image"

type UIButtonType int

const (
	// Define values for UIButtonType must correspond to the order of PlayerGameActionType
	BNT_FoldButton  UIButtonType = 0
	BNT_CheckButton UIButtonType = 1
	BNT_CallButton  UIButtonType = 2
	BNT_RaiseButton UIButtonType = 3
	BNT_AllInButton UIButtonType = 4

	BNT_PauseGameButton    UIButtonType = 5
	BNT_ResumeGameButton   UIButtonType = 6
	BNT_StartGameButton    UIButtonType = 7
	BNT_LeaveGameButton    UIButtonType = 8
	BNT_RequestBuyinButton UIButtonType = 9
	BNT_PaybackBuyinButton UIButtonType = 10

	BNT_JoinSlotButton UIButtonType = 20
)

type Button struct {
	Block
	Text         string
	ActiveStyle  Style
	DisableStyle Style
	DefaultStyle Style

	// Controller is the controller that handles the button's actions.
	Type       UIButtonType
	isSelected bool
	isDisable  bool
}

func NewButton(text string, t UIButtonType) *Button {
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
