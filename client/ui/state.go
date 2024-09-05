package ui

const (
	ROOM_LOGIN = "room_login"
	IN_GAME    = "in_game"
)

type State struct {
	currentScreen string
	currentTab    string
}

func NewState() *State {
	return &State{}
}

func (s *State) SetScreen(screen string) {
	s.currentScreen = screen
}

func (s *State) SetTab(tab string) {
	s.currentTab = tab
}

func (s *State) GetCurrentScreen() string {
	return s.currentScreen
}

func (s *State) GetCurrentTab() string {
	return s.currentTab
}
