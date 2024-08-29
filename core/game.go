package engine

import "fmt"

// Define game rule machanism

type GameSetting struct {
	NumPlayers   int
	MaxStackSize int
	MinStackSize int
	SmallBlind   int
	BigBlind     int
}

type GameStatistcs struct {
	TotalHandsPlayed int
}

type Game struct {
	GameSetting
	GameStatistcs
	gs   GameState
	pm   *PlayerManager
	deck *Deck
}

// NewPokerGame creates a new PokerGame
func NewGame(setting GameSetting, pm *PlayerManager, d *Deck) *Game {
	return &Game{
		GameSetting: setting,
		gs: GameState{
			pot: NewPot(),
			cc:  CommunityCards{},
		},
		pm:   pm,
		deck: d,
	}
}

func (g *Game) Play() {
	fmt.Printf("\n\n\n\nStarting a new game with %d players\n", g.pm.GetNumberOfPlayers())
	// Start the game
	g.pm.ResetForNewGame()
	g.deck.Shuffle()
	g.deck.CutTheCard()
	g.updateDealerPostion()
	g.handleRoundLogic()
	// Log game state
	fmt.Printf("Who is the dealer? => %s\n", g.pm.GetPlayer(g.gs.ButtonPosition).Name())
}

func (g *Game) HandleActions(action ActionIf) {
	if action == nil {
		return
	}
	action.Execute(&g.gs, g.pm)

	// Check if the round is over
	if g.pm.IsAllPlayersActed() {
		g.handleRoundLogic()
	}
}

func (g *Game) resetForNewBettingRound() {
	// Reset player state
	g.pm.ResetForNewRound()
	// Reset new round state
	g.resetGameStateForNewRound()

	// Log current pot when entering new round
	fmt.Printf("Current pot: %d\n", g.gs.pot.Total())
}

func (g *Game) resetGameStateForNewRound() {
	g.gs.CurrentBet = 0
}

func (g *Game) updateDealerPostion() {
	if g.TotalHandsPlayed == 0 {
		// Select the first dealer, choose the player next to the last player
		g.gs.ButtonPosition = g.pm.NextPlayer(g.pm.GetMaxNoSlot()-1, Active).Position()
		fmt.Printf("Selecting the first dealer: %s\n", g.pm.GetPlayer(g.gs.ButtonPosition).Name())
		return
	}
	nextButton := g.pm.NextPlayer(g.gs.ButtonPosition, Active).Position()
	g.gs.ButtonPosition = nextButton
	// Log dealer position
	fmt.Printf("Moving dealer from %d to %d\n", g.gs.ButtonPosition, nextButton)
}

func (g *Game) takeBlinds() {
	sbPlayer := g.pm.NextPlayer(g.gs.ButtonPosition, Active)
	bbPlayer := g.pm.NextPlayer(sbPlayer.Position(), Active)

	sbPlayer.TakeChips(g.SmallBlind)
	sbPlayer.UpdateCurrentBet(g.SmallBlind)

	bbPlayer.TakeChips(g.BigBlind)
	bbPlayer.UpdateCurrentBet(g.BigBlind)

	// Update the current bet
	g.gs.CurrentBet = g.BigBlind

	// Add the blinds to the pot
	g.gs.pot.AddToPot(sbPlayer.Position(), sbPlayer.CurrentBet())
	g.gs.pot.AddToPot(bbPlayer.Position(), bbPlayer.CurrentBet())

	// log take blinds from players successfyully
	fmt.Printf("Small blind %s takes %d chips\n", sbPlayer.Name(), g.SmallBlind)
	fmt.Printf("Big blind %s takes %d chips\n", bbPlayer.Name(), g.BigBlind)

	// Update the next active player
	np := g.pm.NextPlayer(bbPlayer.Position(), Active)
	np.UpdateStatus(WaitForAct)
}

func (g *Game) dealCardsToPlayers() {
	for i := 0; i < 2; i++ {
		for _, p := range g.pm.players {
			if p != nil {
				p.DealCard(g.deck.Draw(), i)
			}
		}
	}

	// If debug mode is on, log all the players' hands
	if DebugMode {
		for _, p := range g.pm.players {
			if p != nil {
				fmt.Printf("%s's hand: [%s]\n", p.Name(), p.ShowHand().String())
			}
		}
	}
}

func (g *Game) dealCommunityCards() {
	switch g.gs.CurrentRound {
	case Flop:
		// Burn a card
		_ = g.deck.Draw()
		// Add 3 cards to the community cards
		g.gs.cc.AddCard(g.deck.Draw())
		g.gs.cc.AddCard(g.deck.Draw())
		g.gs.cc.AddCard(g.deck.Draw())

		// Print the community cards at flop
		fmt.Printf("============ BOARD at FLOP ===========\n%s\n======================================\n", g.gs.cc.String())
	case Turn:
		// Burn a card
		_ = g.deck.Draw()
		// Add a card to the community cards
		g.gs.cc.AddCard(g.deck.Draw())

		// Print the community cards at turn
		fmt.Printf("============ BOARD at TURN ===========\n%s\n======================================\n", g.gs.cc.String())
	case River:
		// Burn a card
		_ = g.deck.Draw()
		// Add a card to the community cards
		g.gs.cc.AddCard(g.deck.Draw())

		// Print the community cards at river
		fmt.Printf("============ BOARD at RIVER ===========\n%s\n======================================\n", g.gs.cc.String())
	default:
		// Log error when dealing community cards at wrong round
		fmt.Printf("Error: dealing community cards at wrong round\n")
	}
}

func (g *Game) handleRoundLogic() {
	switch g.gs.CurrentRound {
	case PreFlop:
		g.resetForNewBettingRound()
		g.takeBlinds()
		g.dealCardsToPlayers()
		// Next player to act is the player next to the big blind
		g.gs.CurrentRound = Flop
	case Flop:
		g.resetForNewBettingRound()
		g.dealCommunityCards()
		// First player to act is the player next to the dealer
		g.pm.NextPlayer(g.gs.ButtonPosition, Active).UpdateStatus(WaitForAct)
		g.gs.CurrentRound = Turn
	case Turn:
		g.resetForNewBettingRound()
		g.dealCommunityCards()
		// First player to act is the player next to the dealer
		g.pm.NextPlayer(g.gs.ButtonPosition, Active).UpdateStatus(WaitForAct)
		g.gs.CurrentRound = River
	case River:
		g.resetForNewBettingRound()
		g.dealCommunityCards()
		// First player to act is the player next to the dealer
		g.pm.NextPlayer(g.gs.ButtonPosition, Active).UpdateStatus(WaitForAct)
		g.gs.CurrentRound = Showdown
	case Showdown:
		// Evaluate hands to find the winner for main pot and side pot
		g.evaluateHands()
	}
}

func (g *Game) evaluateHands() {
	fmt.Println("Evaluating hands to determine the winner")
	var winner Player
	// Printf all players that will be evaluated
	for _, p := range g.pm.players {
		if p != nil && p.Status() != Folded {
			fmt.Printf("Evaluating player %s: [%s]\n", p.Name(), p.ShowHand().String())
			p.ShowHand().Evaluate(&g.gs.cc)
			// Print its rank
			fmt.Printf("Player %s's best hand: [%s] (%s)\n", p.Name(), p.ShowHand().BestHand(), p.ShowHand().HandRankingString())
			// If the player is the first player to be evaluated, set it as the winner
			if winner == nil {
				winner = p
				continue
			}

			if p.ShowHand().Compare(winner.ShowHand()) > 0 {
				winner = p
			} else if p.ShowHand().Compare(winner.ShowHand()) == 0 {
				// Compare the kicker
				if r := compareTiebreakers(p.ShowHand().Kicker(), winner.ShowHand().Kicker()); r > 0 {
					winner = p
				} else if r == 0 {
					fmt.Printf("Player %s and player %s have the same hand\n", p.Name(), winner.Name())
				}
			}
		}
	}

	// Print the winner
	fmt.Printf("The winner is %s with %s\n", winner.Name(), winner.ShowHand().HandRankingString())
}
