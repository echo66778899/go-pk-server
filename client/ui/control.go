package ui

type ButtonMenuType int

const (
	PlayingButton ButtonMenuType = 0
	ControlButton ButtonMenuType = 1
)

type ButtonCtrlCenter struct {
	// State is the state of the UI.
	State
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
	selectedMenu  ButtonMenuType
	CtrlEnabled   bool
	PlayingButton []*Button
	ControlButton []*Button

	// callout handler when button is selected and enter is pressed
	// callout
	enterEventHandler func(ButtonType)
}

func NewButtonCtrlCenter() *ButtonCtrlCenter {
	b := &ButtonCtrlCenter{
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

func (b *ButtonCtrlCenter) SetOutsideHandler(handler func(ButtonType)) {
	b.enterEventHandler = handler
}

func (b *ButtonCtrlCenter) InitButtonPosition() {
	numberOfButton := 5
	step := (CONTROL_PANEL_X_RIGHT - CONTROL_PANEL_X_LEFT) / numberOfButton
	b.checkBtn.SetCenterCoordinator((CONTROL_PANEL_X_LEFT*2+step*3)/2, CONTROL_PANEL_Y)
	b.foldbtn.SetCenterCoordinator((CONTROL_PANEL_X_LEFT*2+step)/2, CONTROL_PANEL_Y)
	b.callBtn.SetCenterCoordinator((CONTROL_PANEL_X_LEFT*2+step*5)/2, CONTROL_PANEL_Y)
	b.raiseBtn.SetCenterCoordinator((CONTROL_PANEL_X_LEFT*2+step*7)/2, CONTROL_PANEL_Y)
	b.allInBtn.SetCenterCoordinator((CONTROL_PANEL_X_LEFT*2+step*9)/2, CONTROL_PANEL_Y)
	b.PlayingButton = []*Button{b.foldbtn, b.checkBtn, b.callBtn, b.raiseBtn, b.allInBtn}
	b.raisePannel.SetCoordinator((CONTROL_PANEL_X_LEFT*2+step*7)/2+5, CONTROL_PANEL_Y-7)
	b.raisePannel.MaxVal = 20000

	numberOfButton = 4
	step = (CONTROL_PANEL_X_RIGHT - CONTROL_PANEL_X_LEFT) / numberOfButton
	b.joinTableBtn.SetCenterCoordinator((CONTROL_PANEL_X_LEFT*2+step)/2, CONTROL_PANEL_Y)
	b.startGameBtn.SetCenterCoordinator((CONTROL_PANEL_X_LEFT*2+step*3)/2, CONTROL_PANEL_Y)
	b.leaveGameBtn.SetCenterCoordinator((CONTROL_PANEL_X_LEFT*2+step*5)/2, CONTROL_PANEL_Y)
	b.requestChipBtn.SetCenterCoordinator((CONTROL_PANEL_X_LEFT*2+step*7)/2, CONTROL_PANEL_Y)
	b.ControlButton = []*Button{b.joinTableBtn, b.startGameBtn, b.leaveGameBtn, b.requestChipBtn}
}

func (b *ButtonCtrlCenter) GetDisplayingButton() []Drawable {
	// Add shold be in order of the buttons
	switch b.selectedMenu {
	case PlayingButton:
		isRaiseSelected := false
		drawables := make([]Drawable, len(b.PlayingButton))
		for i, btn := range b.PlayingButton {
			if btn.IsSelected() && btn.Type == BNT_RaiseButton {
				isRaiseSelected = true
			}
			drawables[i] = btn
		}
		if isRaiseSelected {
			drawables = append(drawables, b.raisePannel)
		}
		return drawables
	case ControlButton:
		drawables := make([]Drawable, len(b.ControlButton))
		for i, btn := range b.ControlButton {
			drawables[i] = btn
		}
		return drawables
	}
	return nil
}

func (b *ButtonCtrlCenter) SetMenu(sel ButtonMenuType) {
	b.selectedMenu = sel
}

func (b *ButtonCtrlCenter) EnableButtonCtrl(enabled bool) {
	b.CtrlEnabled = enabled
	switch b.selectedMenu {
	case PlayingButton:
		for _, btn := range b.PlayingButton {
			btn.Enable(enabled)
		}
	case ControlButton:
		for _, btn := range b.ControlButton {
			btn.Enable(enabled)
		}
	default:
	}
}

func (b *ButtonCtrlCenter) ToggleMenu() {
	switch b.selectedMenu {
	case PlayingButton:
		b.selectedMenu = ControlButton
	case ControlButton:
		b.selectedMenu = PlayingButton
	default:
	}
}

func (b *ButtonCtrlCenter) Enter() {
	if !b.CtrlEnabled {
		return
	}
	switch b.selectedMenu {
	case PlayingButton:
		for _, btn := range b.PlayingButton {
			if btn.IsSelected() {
				b.enterEventHandler(btn.Type)
			}
		}
	case ControlButton:
		for _, btn := range b.ControlButton {
			if btn.IsSelected() {
				b.enterEventHandler(btn.Type)
			}
		}
	default:
	}

	b.EnableButtonCtrl(false)
}

func (b *ButtonCtrlCenter) MoveUp() {
	if !b.CtrlEnabled {
		return
	}

	switch b.selectedMenu {
	case PlayingButton:
		if b.raiseBtn != nil && b.raiseBtn.IsSelected() {
			b.raisePannel.Increase(20)
		}
	case ControlButton:
	default:
	}
}

func (b *ButtonCtrlCenter) MoveDown() {
	if !b.CtrlEnabled {
		return
	}

	switch b.selectedMenu {
	case PlayingButton:
		if b.raiseBtn != nil && b.raiseBtn.IsSelected() {
			b.raisePannel.Decrease(20)
		}
	case ControlButton:
	default:
	}
}

func (b *ButtonCtrlCenter) MoveLeft() {
	if !b.CtrlEnabled {
		return
	}

	switch b.selectedMenu {
	case PlayingButton:
		btnMax := len(b.PlayingButton)
		for i, btn := range b.PlayingButton {
			if btn.IsSelected() {
				btn.Unselect()
				i = (i + btnMax - 1) % btnMax
				b.PlayingButton[i].Select()
				return
			}
		}
		if btnMax > 0 {
			b.PlayingButton[btnMax-1].Select()
		}
	case ControlButton:
		btnMax := len(b.ControlButton)
		for i, btn := range b.ControlButton {
			if btn.IsSelected() {
				btn.Unselect()
				i = (i + btnMax - 1) % btnMax
				b.ControlButton[i].Select()
				return
			}
		}
		if btnMax > 0 {
			b.ControlButton[btnMax-1].Select()
		}
	default:
	}
}

func (b *ButtonCtrlCenter) MoveRight() {
	if !b.CtrlEnabled {
		return
	}
	switch b.selectedMenu {
	case PlayingButton:
		btnMax := len(b.PlayingButton)
		for i, btn := range b.PlayingButton {
			if btn.IsSelected() {
				btn.Unselect()
				i = (i + 1) % btnMax
				b.PlayingButton[i].Select()
				return
			}
		}
		if btnMax > 0 {
			b.PlayingButton[0].Select()
		}
	case ControlButton:
		btnMax := len(b.ControlButton)
		for i, btn := range b.ControlButton {
			if btn.IsSelected() {
				btn.Unselect()
				i = (i + 1) % btnMax
				b.ControlButton[i].Select()
				return
			}
		}
		if btnMax > 0 {
			b.ControlButton[0].Select()
		}
	default:
	}
}
