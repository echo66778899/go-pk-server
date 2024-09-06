package ui

import (
	"log"
	"strconv"
)

type ButtonMenuType int

const (
	ButtonMenuType_PLAYING_BTN ButtonMenuType = 0
	ButtonMenuType_CTRL_BTN    ButtonMenuType = 1
	ButtonMenuType_SLOTS_BTN   ButtonMenuType = 2
)

type ButtonCtrlCenter struct {
	X0, X1, Y int
	// For playing buttons
	foldbtn  *Button
	checkBtn *Button
	callBtn  *Button
	raiseBtn *Button
	allInBtn *Button

	// additional pannel
	raisePannel *RaiseWidget

	// For control buttons
	joinTableBtn   *Button
	startGameBtn   *Button
	leaveGameBtn   *Button
	requestChipBtn *Button

	// ButtonCtrl is the controller for the buttons.
	availableSlots   int
	selectingBtnMenu ButtonMenuType
	CtrlEnabled      bool
	PlayingButtons   []*Button
	ControlButtons   []*Button
	SlotsButtons     []*Button

	// callout handler when button is selected and enter is pressed
	// callout
	enterEventHandler func(...int)
}

func NewButtonCtrlCenter() *ButtonCtrlCenter {
	b := &ButtonCtrlCenter{
		X0:       CONTROL_PANEL_X_LEFT,
		X1:       CONTROL_PANEL_X_RIGHT,
		Y:        CONTROL_PANEL_Y,
		foldbtn:  NewButton("Fold", BNT_FoldButton),
		checkBtn: NewButton("Check", BNT_CheckButton),
		callBtn:  NewButton("Call", BNT_CallButton),
		raiseBtn: NewButton("Raise", BNT_RaiseButton),
		allInBtn: NewButton("All-In", BNT_AllInButton),

		raisePannel: NewRaiseWidget(),

		joinTableBtn:   NewButton("Join Table", BNT_JoinTableButton),
		startGameBtn:   NewButton("Start Game", BNT_StartGameButton),
		leaveGameBtn:   NewButton("Leave Game", BNT_LeaveGameButton),
		requestChipBtn: NewButton("Request Chip", BNT_RequestChipButton),

		CtrlEnabled: true,
	}
	return b
}

func (b *ButtonCtrlCenter) SetUserButtonInteractionHandler(handler func(...int)) {
	b.enterEventHandler = handler
}

func (b *ButtonCtrlCenter) InitButtonPosition() {
	numberOfButton := 5
	step := (b.X1 - b.X0) / numberOfButton
	// Set the center of the button
	b.foldbtn.SetCenterCoordinator((b.X0*2+step)/2, b.Y)
	b.checkBtn.SetCenterCoordinator((b.X0*2+step*3)/2, b.Y)
	b.callBtn.SetCenterCoordinator((b.X0*2+step*5)/2, b.Y)
	b.raiseBtn.SetCenterCoordinator((b.X0*2+step*7)/2, b.Y)
	b.allInBtn.SetCenterCoordinator((b.X0*2+step*9)/2, b.Y)
	// Add the playing buttons to the list
	b.PlayingButtons = []*Button{b.foldbtn, b.checkBtn, b.callBtn, b.raiseBtn, b.allInBtn}

	b.raisePannel.SetCoordinator((b.X0*2+step*7)/2+5, b.Y-8)
	b.raisePannel.MaxVal = 5000

	numberOfButton = 4
	step = (CONTROL_PANEL_X_RIGHT - b.X0) / numberOfButton
	// Set the center of the button
	b.joinTableBtn.SetCenterCoordinator((b.X0*2+step)/2, b.Y)
	b.startGameBtn.SetCenterCoordinator((b.X0*2+step*3)/2, b.Y)
	b.leaveGameBtn.SetCenterCoordinator((b.X0*2+step*5)/2, b.Y)
	b.requestChipBtn.SetCenterCoordinator((b.X0*2+step*7)/2, b.Y)
	// Add the control buttons to the list
	b.ControlButtons = []*Button{b.joinTableBtn, b.startGameBtn, b.leaveGameBtn, b.requestChipBtn}
}

func (b *ButtonCtrlCenter) GetDisplayingButton() []Drawable {
	// Add shold be in order of the buttons
	switch b.selectingBtnMenu {
	case ButtonMenuType_PLAYING_BTN:
		isRaiseSelected := false
		drawables := make([]Drawable, len(b.PlayingButtons))
		for i, btn := range b.PlayingButtons {
			if btn.IsSelected() && btn.Type == BNT_RaiseButton {
				isRaiseSelected = true
			}
			drawables[i] = btn
		}
		if isRaiseSelected {
			drawables = append(drawables, b.raisePannel)
		}
		return drawables
	case ButtonMenuType_CTRL_BTN:
		drawables := make([]Drawable, len(b.ControlButtons))
		for i, btn := range b.ControlButtons {
			drawables[i] = btn
		}
		return drawables
	case ButtonMenuType_SLOTS_BTN:
		drawables := make([]Drawable, len(b.SlotsButtons))
		for i, btn := range b.SlotsButtons {
			drawables[i] = btn
		}
		return drawables
	}
	return nil
}

func (b *ButtonCtrlCenter) UpdateState() {
	players := UI_MODEL_DATA.Players
	log.Printf("Players: %v, MaxPlayers: %v", players, UI_MODEL_DATA.MaxPlayers)

	// Calculate the number of available slots from the player list
	available := UI_MODEL_DATA.MaxPlayers
	emptySlots := make(map[int]bool)
	for i := 0; i < UI_MODEL_DATA.MaxPlayers; i++ {
		emptySlots[i] = true
	}
	for _, p := range players {
		if p != nil {
			available--
			emptySlots[int(p.TablePosition)] = false
		}
	}

	// Only need to new number of buttons if the number of available slots change
	if available != b.availableSlots {
		b.availableSlots = available
		i, step := 0, (CONTROL_PANEL_X_RIGHT-b.X0)/b.availableSlots
		b.SlotsButtons = make([]*Button, 0)
		for slot, empty := range emptySlots {
			if empty {
				btn := NewButton("Slot "+strconv.Itoa(slot), BNT_SlotButton)
				btn.SetCenterCoordinator((b.X0*2+(step*((2*i)+1)))/2, b.Y)
				b.SlotsButtons = append(b.SlotsButtons, btn)
				i++
			}
		}
	}
}

func (b *ButtonCtrlCenter) GoNextMenu() {
	switch b.selectingBtnMenu {
	case ButtonMenuType_PLAYING_BTN:
		b.selectingBtnMenu = ButtonMenuType_CTRL_BTN
	case ButtonMenuType_CTRL_BTN:
		b.selectingBtnMenu = ButtonMenuType_SLOTS_BTN
	case ButtonMenuType_SLOTS_BTN:
		b.selectingBtnMenu = ButtonMenuType_PLAYING_BTN
	default:
	}
}

func (b *ButtonCtrlCenter) GetSelectingButtons() []*Button {
	switch b.selectingBtnMenu {
	case ButtonMenuType_PLAYING_BTN:
		return b.PlayingButtons
	case ButtonMenuType_CTRL_BTN:
		return b.ControlButtons
	case ButtonMenuType_SLOTS_BTN:
		return b.SlotsButtons
	default:
	}
	return make([]*Button, 0)
}

func (b *ButtonCtrlCenter) SetMenu(sel ButtonMenuType) {
	b.selectingBtnMenu = sel
}

func (b *ButtonCtrlCenter) EnableButtonCtrl(enabled bool) {
	b.CtrlEnabled = enabled
	log.Printf("ButtonCtrlCenter - EnableButtonCtrl: %v", enabled)
	for _, btn := range b.GetSelectingButtons() {
		btn.Enable(enabled)
	}
}

func (b *ButtonCtrlCenter) Enter() {
	if !b.CtrlEnabled {
		log.Printf("Enter - is disabled")
		return
	}
	for _, btn := range b.GetSelectingButtons() {
		if btn.IsSelected() {
			if btn.Type == BNT_SlotButton {
				convertToSlotNo := func(s string) int {
					// [space][space]Slot[Space][Number][space][space]
					slot, _ := strconv.Atoi(s[7 : len(s)-2])
					return slot
				}
				b.enterEventHandler(int(btn.Type), convertToSlotNo(btn.Text))
			} else if btn.Type == BNT_RaiseButton {
				b.enterEventHandler(int(btn.Type), b.raisePannel.Data)
			} else {
				b.enterEventHandler(int(btn.Type))
			}
		}
	}
	log.Printf("Pressed Enter - disabling button control")
	b.EnableButtonCtrl(false)
}

func (b *ButtonCtrlCenter) MoveUp() {
	if !b.CtrlEnabled {
		log.Printf("MoveUp - is disabled")
		return
	}

	switch b.selectingBtnMenu {
	case ButtonMenuType_PLAYING_BTN:
		if b.raiseBtn != nil && b.raiseBtn.IsSelected() {
			b.raisePannel.Increase(20)
		}
	case ButtonMenuType_CTRL_BTN:
	default:
	}
}

func (b *ButtonCtrlCenter) MoveDown() {
	if !b.CtrlEnabled {
		log.Printf("MoveDown - is disabled")
		return
	}

	switch b.selectingBtnMenu {
	case ButtonMenuType_PLAYING_BTN:
		if b.raiseBtn != nil && b.raiseBtn.IsSelected() {
			b.raisePannel.Decrease(20)
		}
	case ButtonMenuType_CTRL_BTN:
	default:
	}
}

func (b *ButtonCtrlCenter) MoveLeft() {
	if !b.CtrlEnabled {
		log.Printf("MoveLeft - is disabled")
		return
	}

	buttons := b.GetSelectingButtons()
	btnMax := len(buttons)

	for i, btn := range buttons {
		if btn.IsSelected() {
			btn.Unselect()
			i = (i + btnMax - 1) % btnMax
			buttons[i].Select()
			return
		}
	}
	if btnMax > 0 {
		buttons[btnMax-1].Select()
	}
}

func (b *ButtonCtrlCenter) MoveRight() {
	if !b.CtrlEnabled {
		log.Printf("MoveRight - is disabled")
		return
	}

	buttons := b.GetSelectingButtons()
	btnMax := len(buttons)

	for i, btn := range buttons {
		if btn.IsSelected() {
			btn.Unselect()
			i = (i + 1) % btnMax
			buttons[i].Select()
			return
		}
	}
	if btnMax > 0 {
		buttons[0].Select()
	}
}
