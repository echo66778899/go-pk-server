package ui

import (
	"log"
	"strconv"

	gpbmessage "go-pk-server/gen"
)

type UIButtonMenuType int

const (
	ButtonMenuType_PLAYING_BTN UIButtonMenuType = 0
	ButtonMenuType_CTRL_BTN    UIButtonMenuType = 1
	ButtonMenuType_SLOTS_BTN   UIButtonMenuType = 2
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

	// call chips text
	toCallChipTxt *CallChipText

	// For control buttons
	pauseGameBtn   *Button
	startGameBtn   *Button
	leaveGameBtn   *Button
	requestChipBtn *Button
	paybackChipBtn *Button

	// ButtonCtrl is the controller for the buttons.
	availableSlots   int
	selectingBtnMenu UIButtonMenuType
	IsVisible        bool
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
		foldbtn:  NewButton("[F]old", BNT_FoldButton),
		checkBtn: NewButton("[C]heck", BNT_CheckButton),
		callBtn:  NewButton("[C]all", BNT_CallButton),
		raiseBtn: NewButton("[R]aise", BNT_RaiseButton),
		allInBtn: NewButton("[A]ll-In", BNT_AllInButton),

		raisePannel:   NewRaiseWidget(),
		toCallChipTxt: NewCallChipText(),

		startGameBtn:   NewButton("[S]tart Game", BNT_StartGameButton),
		pauseGameBtn:   NewButton("[P]ause Game", BNT_PauseGameButton),
		leaveGameBtn:   NewButton("[L]eave Game", BNT_LeaveGameButton),
		requestChipBtn: NewButton("[T]ake Buy-In", BNT_RequestBuyinButton),
		paybackChipBtn: NewButton("[G]ive Buy-In", BNT_PaybackBuyinButton),

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

	b.toCallChipTxt.SetRect((b.X0*2+step*5)/2-2, b.Y-2, (b.X0*2+step*5)/2+6, b.Y-1)
	// Set the center of the raise pannel
	b.raisePannel.SetCoordinator((b.X0*2+step*7)/2+5, b.Y-8)
	b.raisePannel.MaxVal = 10000

	numberOfButton = 5
	step = (CONTROL_PANEL_X_RIGHT - b.X0) / numberOfButton
	// Set the center of the button
	b.startGameBtn.SetCenterCoordinator((b.X0*2+step)/2, b.Y)
	b.pauseGameBtn.SetCenterCoordinator((b.X0*2+step*3)/2, b.Y)
	b.leaveGameBtn.SetCenterCoordinator((b.X0*2+step*5)/2, b.Y)
	b.requestChipBtn.SetCenterCoordinator((b.X0*2+step*7)/2, b.Y)
	b.paybackChipBtn.SetCenterCoordinator((b.X0*2+step*9)/2, b.Y)
	// Add the control buttons to the list
	b.ControlButtons = []*Button{b.startGameBtn, b.pauseGameBtn, b.leaveGameBtn, b.requestChipBtn, b.paybackChipBtn}
}

func (b *ButtonCtrlCenter) GetDisplayingButton() []Drawable {
	if !b.IsVisible {
		return nil
	}
	// Add shold be in order of the buttons
	switch b.selectingBtnMenu {
	case ButtonMenuType_PLAYING_BTN:
		isRaiseSelected, isCallVisible := false, false
		drawables := make([]Drawable, len(b.PlayingButtons))
		for i, btn := range b.PlayingButtons {
			if btn.IsSelected() && btn.Type == BNT_RaiseButton {
				isRaiseSelected = true
			}
			if btn.Type == BNT_CallButton && btn.IsEnabled() {
				isCallVisible = true
			}
			drawables[i] = btn
		}
		if isRaiseSelected {
			drawables = append(drawables, b.raisePannel)
		}
		if isCallVisible {
			drawables = append(drawables, b.toCallChipTxt)
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
				btn := NewButton("Slot "+strconv.Itoa(slot), BNT_JoinSlotButton)
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

func (b *ButtonCtrlCenter) SetMenu(sel UIButtonMenuType) {
	b.selectingBtnMenu = sel
}

func (b *ButtonCtrlCenter) DisableListButton(p *gpbmessage.PlayerState) {
	if p == nil {
		return
	}
	notAllowAction := p.GetNoActions()
	if len(notAllowAction) <= 0 {
		return
	}
	// Create a set (map) of types to disable for faster lookup
	disableSet := make(map[gpbmessage.PlayerGameActionType]struct{}, len(notAllowAction))
	for _, t := range notAllowAction {
		disableSet[t] = struct{}{}
	}

	// Iterate over buttons and disable only if in the disableSet
	for _, btn := range b.GetSelectingButtons() {
		if _, found := disableSet[gpbmessage.PlayerGameActionType(btn.Type)]; found {
			btn.Enable(false)
		}
	}
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
			if btn.Type == BNT_JoinSlotButton {
				convertToSlotNo := func(s string) int {
					// [space][space]Slot[Space][Number][space][space]
					slot, _ := strconv.Atoi(s[7 : len(s)-2])
					return slot
				}
				b.enterEventHandler(int(btn.Type), convertToSlotNo(btn.Text))
			} else if btn.Type == BNT_RaiseButton {
				b.enterEventHandler(int(btn.Type), b.raisePannel.Data)
				b.raisePannel.Data = 0
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
			b.raisePannel.Increase(UI_MODEL_DATA.CurrentBet)
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
			b.raisePannel.Decrease(UI_MODEL_DATA.CurrentBet)
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

	// Find the currently selected button
	var currentIndex int = -1
	for i, btn := range buttons {
		if btn.IsSelected() {
			currentIndex = i
			btn.Unselect()
			break
		}
	}

	// If no button was selected, start from the last button
	if currentIndex == -1 {
		currentIndex = btnMax
	}

	// Move to the previous enabled button
	for j := 1; j <= btnMax; j++ {
		prevIndex := (currentIndex + btnMax - j) % btnMax
		if buttons[prevIndex].IsEnabled() {
			buttons[prevIndex].Select()
			return
		}
	}

	for i, btn := range buttons {
		if btn.IsSelected() {
			btn.Unselect()
			for j := 1; j < btnMax; j++ {
				prevIndex := (i + btnMax - j) % btnMax
				if buttons[prevIndex].IsEnabled() {
					buttons[prevIndex].Select()
					return
				}
			}
		}
	}
	// If no enabled buttons were found, select the last one
	buttons[btnMax-1].Select()
}

func (b *ButtonCtrlCenter) MoveRight() {
	if !b.CtrlEnabled {
		log.Printf("MoveRight - is disabled")
		return
	}

	buttons := b.GetSelectingButtons()
	btnMax := len(buttons)

	// Find the currently selected button
	var currentIndex int = -1
	for i, btn := range buttons {
		if btn.IsSelected() {
			currentIndex = i
			btn.Unselect()
			break
		}
	}

	// If no button was selected, start from the first button
	if currentIndex == -1 {
		currentIndex = -1
	}

	// Move to the next enabled button
	for j := 1; j <= btnMax; j++ {
		nextIndex := (currentIndex + j) % btnMax
		if buttons[nextIndex].IsEnabled() {
			buttons[nextIndex].Select()
			return
		}
	}

	// If no enabled buttons were found, select the first one
	buttons[0].Select()
}
