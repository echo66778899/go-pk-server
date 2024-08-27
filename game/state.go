package engine

type GameState struct {
	Players          []Player
	Pots             Pot
	CommunityCards   CommunityCards
	CurrentButtonIdx int
	CurrentPlayerIdx int
	CurrentBet       int
	NumPlayingPlayer int
	// Add any other necessary fields here
}
