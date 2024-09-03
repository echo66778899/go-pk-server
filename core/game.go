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
	tm   *TableManager
	deck *Deck
}

// NewPokerGame creates a new PokerGame
func NewGame(setting GameSetting, tm *TableManager, d *Deck) *Game {
	return &Game{
		GameSetting: setting,
		gs: GameState{
			pot: NewPot(),
			cc:  CommunityCards{},
		},
		tm:   tm,
		deck: d,
	}
}

func (g *Game) Play() {
	// Check if the number of players is valid
	if g.tm.GetNumberOfPlayers() < 2 {
		// Log error when the number of players is less than 2
		fmt.Println("error: Number of players is less than 2")
		return
	}

	fmt.Printf("\n\n\n\nStarting a new game with %d players\n", g.tm.GetNumberOfPlayers())
	// Start the game
	g.tm.ResetForNewGame()
	g.gs.NumPlayingPlayer = g.tm.GetNumberOfPlayingPlayers()

	// Shuffle the deck
	g.deck.Shuffle()
	g.deck.CutTheCard()

	// Update the dealer position
	g.updateDealerPostion()
	g.gs.pot.ResetPot()

	// Start the first round
	g.handleEnterNewRoundLogic()

	// Log game state
	fmt.Printf("Who is the dealer? => %s\n", g.tm.GetPlayer(g.gs.ButtonPosition).Name())
	fmt.Printf("Number of playing players: %d\n", g.gs.NumPlayingPlayer)

}

func (g *Game) NextGame() {
	// Reset the game state for a new game
	if g.gs.CurrentRound == PreFlop {
		// Start the first round
		g.handleEnterNewRoundLogic()
		return
	}

	// Log error when the game is not in the PreFlop round
	fmt.Println("error: Can not start a new game when the game is not in the PreFlop round")
}

func (g *Game) HandleActions(action ActionIf) {
	if action == nil {
		return
	}

	// Check if the action is valid
	player := g.tm.GetPlayer(action.FromWho())
	if player.Status() != WaitForAct {
		// Log warning the player is not allowed to fold, the action is invalid
		fmt.Println("error: Player", player.Name(), "is not allowed to ", action.WhatAction(), ", the action is invalid")
		return
	}

	fmt.Printf("Handling %v action from player %s.\n", action.WhatAction(), player.Name())
	fmt.Printf("BEFORE Current bet: %d\n", g.gs.CurrentBet)

	switch action.WhatAction() {
	case Fold:
		// Execute fold action
		player.UpdateStatus(Folded)
		player.UpdateSuggestions([]PlayerActType{})
		// Decrease the number of playing players
		g.gs.NumPlayingPlayer--

	case Check:
		// Execute check action
		// Check if the player is allowed to check
		if player.CurrentBet() == g.gs.CurrentBet {
			player.UpdateStatus(Checked)
		} else {
			// Log info the player name is not allowed to check, the action is invalid
			fmt.Println("error: Player", player.Name(), "is not allowed to check, the action is invalid")
		}
	case Call:
		// Execute call action
		// If the player chip is less than the current bet, the player is all-in
		callChip := g.gs.CurrentBet - player.CurrentBet()
		if player.Chips() <= callChip {
			player.UpdateCurrentBet(player.Chips() + player.CurrentBet())
			player.TakeChips(player.Chips())
			player.UpdateStatus(AlledIn)
		} else {
			player.UpdateCurrentBet(callChip + player.CurrentBet())
			player.TakeChips(callChip)
			player.UpdateStatus(Called)
		}
		g.gs.pot.AddToPot(player.Position(), callChip)
		g.gs.CurrentBet = player.CurrentBet()
	case Raise:
		// Execute raise action
		raiseAmount := action.HowMuch()
		// log the raise amount
		fmt.Printf("Player %s raises %d\n", player.Name(), raiseAmount)
		callAmount := g.gs.CurrentBet - player.CurrentBet()

		if raiseAmount+callAmount < player.Chips() {
			raiseAmount += callAmount
			player.UpdateCurrentBet(raiseAmount + player.CurrentBet())
			player.TakeChips(raiseAmount)
			player.UpdateStatus(Raised)
			g.gs.pot.AddToPot(player.Position(), raiseAmount)
		} else {
			raiseAmount := player.Chips()
			player.UpdateCurrentBet(raiseAmount + player.CurrentBet())
			player.TakeChips(raiseAmount)
			player.UpdateStatus(AlledIn)
		}
		g.gs.CurrentBet = player.CurrentBet()
		// Update all player status to Playing and NextPlayer to WaitForAct
		for _, p := range g.tm.GetListOfOtherPlayers(action.FromWho(), Called, Raised, Checked) {
			p.UpdateStatus(Playing)
		}
	case AllIn:
		// Execute all-in action
		allInAmount := player.Chips()
		if allInAmount > g.gs.CurrentBet {
			player.UpdateCurrentBet(allInAmount + player.CurrentBet())
			player.TakeChips(allInAmount)
			player.UpdateStatus(AlledIn)
			g.gs.pot.AddToPot(player.Position(), allInAmount)
			g.gs.CurrentBet = player.CurrentBet()

			// Update all player status to Playing and NextPlayer to WaitForAct
			for _, p := range g.tm.GetListOfOtherPlayers(action.FromWho(), Called, Raised, Checked) {
				p.UpdateStatus(Playing)
			}
		} else {
			// Log warning the player should go all-in
			fmt.Println("error: Player", player.Name(), "should go all-in rather than Raise")
		}
	default:
		// Log invalid action
		fmt.Printf("error: Invalid action: %s\n", action.WhatAction())
		return
	}

	fmt.Printf("AFTER  Current bet: %d, Number of Playing: %d\n", g.gs.CurrentBet, g.gs.NumPlayingPlayer)

	if np := g.tm.NextPlayer(action.FromWho(), Playing); np != nil {
		np.UpdateStatus(WaitForAct)
		switch action.WhatAction() {
		case Check, Fold:
			if g.gs.CurrentBet == 0 {
				np.UpdateSuggestions([]PlayerActType{Check, Raise, AllIn})
			} else {
				np.UpdateSuggestions([]PlayerActType{Fold, Call, Raise, AllIn})
			}
		case Call, Raise, AllIn:
			if np.Chips() < g.gs.CurrentBet*2 {
				np.UpdateSuggestions([]PlayerActType{Fold, Call, AllIn})
			} else {
				np.UpdateSuggestions([]PlayerActType{Fold, Call, Raise, AllIn})
			}
		default:
			fmt.Println("error: Can not suggest action for player", np.Name())
		}
	} else {
		// Can not find the next player, the round is over
		fmt.Println("Warning: Can not find the next player, the round is over")
		if g.gs.NumPlayingPlayer <= 1 {
			g.gs.CurrentRound = Showdown
		}
		g.handleEnterNewRoundLogic()
		return
	}
}

func (g *Game) resetForNewBettingRound() {
	// Reset player state
	g.tm.ResetForNewRound()
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
		p := g.tm.NextPlayer(g.tm.GetMaxNoSlot()-1, Playing)
		if p == nil {
			// Log error when selecting the first dealer
			fmt.Println("error: Can not select the first dealer")
			return
		}
		g.gs.ButtonPosition = p.Position()
		fmt.Printf("Selecting the first dealer: %s\n", g.tm.GetPlayer(g.gs.ButtonPosition).Name())
		return
	}
	nextButton := g.tm.NextPlayer(g.gs.ButtonPosition, Playing).Position()
	// Log dealer position
	fmt.Printf("Moving dealer from player %v to player %v\n", g.tm.GetPlayer(g.gs.ButtonPosition).Name(), g.tm.GetPlayer(nextButton).Name())
	g.gs.ButtonPosition = nextButton
}

func (g *Game) takeBlinds() {
	sbPlayer := g.tm.NextPlayer(g.gs.ButtonPosition, Playing)
	bbPlayer := g.tm.NextPlayer(sbPlayer.Position(), Playing)

	if sbPlayer == nil || bbPlayer == nil {
		// Log error when taking blinds
		fmt.Println("error: Can not take blinds")
		return
	}

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
	np := g.tm.NextPlayer(bbPlayer.Position(), Playing)
	np.UpdateStatus(WaitForAct)
	np.UpdateSuggestions([]PlayerActType{Fold, Call, Raise, AllIn})
}

func (g *Game) dealCardsToPlayers() {
	for i := 0; i < 2; i++ {
		for _, p := range g.tm.players {
			if p != nil {
				p.DealCard(g.deck.Draw(), i)
			}
		}
	}

	// If debug mode is on, log all the players' hands
	if DebugMode {
		for _, p := range g.tm.players {
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
		fmt.Printf("error: dealing community cards at wrong round\n")
	}
}

func (g *Game) dealTheRestOfCommunityCards() {
	communityCardsLen := len(g.gs.cc.Cards)
	needDraw := 5 - communityCardsLen
	if needDraw > 3 {
		// Burn a card
		_ = g.deck.Draw()
		// Add 3 cards to the community cards
		g.gs.cc.AddCard(g.deck.Draw())
		g.gs.cc.AddCard(g.deck.Draw())
		g.gs.cc.AddCard(g.deck.Draw())
		// Burn a card
		_ = g.deck.Draw()
		// Add a card to the community cards
		g.gs.cc.AddCard(g.deck.Draw())
		// Burn a card
		_ = g.deck.Draw()
		// Add a card to the community cards
		g.gs.cc.AddCard(g.deck.Draw())
	} else if needDraw > 1 {
		// Burn a card
		_ = g.deck.Draw()
		// Add a card to the community cards
		g.gs.cc.AddCard(g.deck.Draw())
		// Burn a card
		_ = g.deck.Draw()
		// Add a card to the community cards
		g.gs.cc.AddCard(g.deck.Draw())
	} else if needDraw > 0 {
		// Burn a card
		_ = g.deck.Draw()
		// Add a card to the community cards
		g.gs.cc.AddCard(g.deck.Draw())
	} else {
		fmt.Println("Enough cards in the commuuity cards")
	}

	fmt.Printf("============ BOARD at ENDED ===========\n%s\n======================================\n", g.gs.cc.String())
}

func (g *Game) firstPlayerActionInRound() bool {
	// If there is only one player, the player wins the pot
	if g.tm.GetNumberOfPlayingPlayers() == 1 {
		// Log only one player in the game
		fmt.Println("Only one player in the game!")
		return false
	}

	// First player to act is the player next to the dealer
	np := g.tm.NextPlayer(g.gs.ButtonPosition, Playing)

	if np != nil {
		np.UpdateStatus(WaitForAct)
		np.UpdateSuggestions([]PlayerActType{Check, Raise, AllIn})
		return true
	}

	return false
}

func (g *Game) handleEnterNewRoundLogic() {
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
		if !g.firstPlayerActionInRound() {
			g.dealTheRestOfCommunityCards()
			g.evaluateHands()
			g.resetForNewGame()
		} else {
			g.gs.CurrentRound = Turn
		}
	case Turn:
		g.resetForNewBettingRound()
		g.dealCommunityCards()
		// First player to act is the player next to the dealer
		if !g.firstPlayerActionInRound() {
			g.dealTheRestOfCommunityCards()
			g.evaluateHands()
			g.resetForNewGame()
		} else {
			g.gs.CurrentRound = River
		}
	case River:
		g.resetForNewBettingRound()
		g.dealCommunityCards()

		if !g.firstPlayerActionInRound() {
			g.dealTheRestOfCommunityCards()
			g.evaluateHands()
			g.resetForNewGame()
		}
		g.gs.CurrentRound = Showdown
	case Showdown:
		// Evaluate hands to find the winner for main pot and side pot
		g.evaluateHands()
		g.resetForNewGame()
	}
}

func (g *Game) resetForNewGame() {
	// Reset the game state for a new game
	g.tm.ResetForNewGame()
	g.gs.NumPlayingPlayer = g.tm.GetNumberOfPlayingPlayers()

	// Shuffle the deck
	g.deck.Shuffle()
	g.deck.CutTheCard()

	// Log statistics
	g.TotalHandsPlayed++
	fmt.Printf("Total hands played: %d\n", g.TotalHandsPlayed)

	// Move the dealer position
	g.updateDealerPostion()
	g.gs.pot.ResetPot()
	g.gs.cc.Reset()
	g.gs.CurrentRound = PreFlop
}

func (g *Game) evaluateHands() {
	if g.gs.NumPlayingPlayer == 1 {
		// Log the winner
		onePlayer := g.tm.GetListOfPlayers(Playing, Called, Checked, Raised, AlledIn)
		if len(onePlayer) != 1 {
			panic("error: more than one player in the game")
		}
		fmt.Printf("Player %s wins the pot (%d) with a hand down\n", onePlayer[0].Name(), g.gs.pot.Total())
		onePlayer[0].AddChips(g.gs.pot.Total())
	} else {
		// Evaluate hands to find the winner for main pot and side pot

		fmt.Println("Evaluating hands to determine the winner")
		var winner Player
		// Printf all players that will be evaluated
		for _, p := range g.tm.players {
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
}
