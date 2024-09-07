package engine

// Define game rule machanism

import (
	"fmt"
	msgpb "go-pk-server/gen"
	mylog "go-pk-server/log"
)

type GameStatistcs struct {
	TotalHandsPlayed int
}

type Game struct {
	GameStatistcs
	setting *msgpb.GameSetting
	gs      GameState
	tm      *TableManager
	deck    *Deck

	funcReqEngineState func(EngineState)
}

// NewPokerGame creates a new PokerGame
func NewGame(setting *msgpb.GameSetting, tm *TableManager, d *Deck, reqEState func(EngineState)) *Game {
	tm.UpdateMaxNoOfSlot(int(setting.MaxPlayers))
	return &Game{
		setting: setting,
		gs: GameState{
			pot: NewPot(),
			cc:  CommunityCards{},
		},
		tm:                 tm,
		deck:               d,
		funcReqEngineState: reqEState,
	}
}

func (g *Game) Play() bool {
	// Check if the number of players is valid
	if g.tm.GetNumberOfPlayers() < 2 {
		// Log error when the number of players is less than 2
		fmt.Println("error: Number of players is less than 2")
		return false
	}

	g.gs.NumPlayingPlayer = g.tm.GetNumberOfPlayingPlayers()
	g.updateDealerPostion(g.TotalHandsPlayed == 0)

	mylog.Infof("Starting a NEW_GAME with %d players", g.tm.GetNumberOfPlayers())
	mylog.Infof("Number of playing players: %d\n", g.gs.NumPlayingPlayer)
	mylog.Debugf("Who is the dealer? => %s\n", g.tm.GetPlayer(g.gs.ButtonPosition).Name())

	// Shuffle the deck
	g.deck.Shuffle()
	g.deck.CutTheCard()

	// Start the first round
	g.handleEnterNewRoundLogic()

	return true
}

func (g *Game) HandleActions(action ActionIf) {
	if action == nil {
		return
	}

	// Check if the action is valid
	player := g.tm.GetPlayer(action.FromWho())
	if player == nil {
		return
	}

	if player.Status() != PlayerStatus_Wait4Act {
		// Log warning the player is not allowed to fold, the action is invalid
		fmt.Println("error: Player", player.Name(), "is not allowed to ", action.WhatAction(), ", the action is invalid")
		return
	}

	fmt.Printf("Handling %v action from player %s.\n", action.WhatAction(), player.Name())
	fmt.Printf("BEFORE Current bet: %d\n", g.gs.CurrentBet)

	switch action.WhatAction() {
	case Fold:
		// Execute fold action
		player.UpdateStatus(PlayerStatus_Fold)
		player.UpdateSuggestions([]PlayerActType{})
		// Decrease the number of playing players
		g.gs.NumPlayingPlayer--

	case Check:
		// Execute check action
		// Check if the player is allowed to check
		if player.CurrentBet() == g.gs.CurrentBet {
			player.UpdateStatus(PlayerStatus_Check)
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
			player.UpdateStatus(PlayerStatus_AllIn)
		} else {
			player.UpdateCurrentBet(callChip + player.CurrentBet())
			player.TakeChips(callChip)
			player.UpdateStatus(PlayerStatus_Call)
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
			player.UpdateStatus(PlayerStatus_Raise)
			g.gs.pot.AddToPot(player.Position(), raiseAmount)
		} else {
			raiseAmount := player.Chips()
			player.UpdateCurrentBet(raiseAmount + player.CurrentBet())
			player.TakeChips(raiseAmount)
			player.UpdateStatus(PlayerStatus_AllIn)
		}
		g.gs.CurrentBet = player.CurrentBet()
		// Update all player status to PlayerStatus_Playing and NextPlayer to PlayerStatus_Wait4Act
		for _, p := range g.tm.GetListOfOtherPlayers(action.FromWho(), PlayerStatus_Call, PlayerStatus_Raise, PlayerStatus_Check) {
			p.UpdateStatus(PlayerStatus_Playing)
		}
	case AllIn:
		// Execute all-in action
		allInAmount := player.Chips()
		if allInAmount > g.gs.CurrentBet {
			player.UpdateCurrentBet(allInAmount + player.CurrentBet())
			player.TakeChips(allInAmount)
			player.UpdateStatus(PlayerStatus_AllIn)
			g.gs.pot.AddToPot(player.Position(), allInAmount)
			g.gs.CurrentBet = player.CurrentBet()

			// Update all player status to PlayerStatus_Playing and NextPlayer to PlayerStatus_Wait4Act
			for _, p := range g.tm.GetListOfOtherPlayers(action.FromWho(), PlayerStatus_Call, PlayerStatus_Raise, PlayerStatus_Check) {
				p.UpdateStatus(PlayerStatus_Playing)
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

	fmt.Printf("AFTER  Current bet: %d, Number of PlayerStatus_Playing: %d\n", g.gs.CurrentBet, g.gs.NumPlayingPlayer)

	if np := g.tm.NextPlayer(action.FromWho(), PlayerStatus_Playing); np != nil {
		np.UpdateStatus(PlayerStatus_Wait4Act)
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
			g.gs.CurrentRound = msgpb.RoundStateType_SHOW_DOWN
		}
		g.handleEnterNewRoundLogic()
		return
	}
}

func (g *Game) prepareForIncomingGame() {
	// Reset the game state for a new game
	g.tm.ResetForNewGame()

	// Log statistics
	g.TotalHandsPlayed++
	fmt.Printf("Total hands played: %d\n", g.TotalHandsPlayed)

	// Move the dealer position
	g.updateDealerPostion(false)
	g.gs.pot.ResetPot()

	g.gs.cc.Reset()
	g.gs.CurrentRound = msgpb.RoundStateType_PREFLOP

	// Update engine state to wait for start plauing game
	g.funcReqEngineState(EngineState_WAIT_FOR_PLAYING)
}

func (g *Game) prepareForNewBettingRound() {
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

func (g *Game) updateDealerPostion(firstGame bool) {
	if firstGame {
		// Select the first dealer, choose the player next to the last player
		p := g.tm.NextPlayer(g.tm.GetMaxNoSlot()-1, PlayerStatus_Playing)
		if p == nil {
			// Log error when selecting the first dealer
			fmt.Println("error: Can not select the first dealer")
			return
		}
		g.gs.ButtonPosition = p.Position()
		fmt.Printf("Selecting the first dealer: %s\n", g.tm.GetPlayer(g.gs.ButtonPosition).Name())
		return
	}
	nextButton := g.tm.NextPlayer(g.gs.ButtonPosition, PlayerStatus_Playing).Position()
	// Log dealer position
	fmt.Printf("Moving dealer from player %v to player %v\n", g.tm.GetPlayer(g.gs.ButtonPosition).Name(), g.tm.GetPlayer(nextButton).Name())
	g.gs.ButtonPosition = nextButton
}

func (g *Game) takeBlinds() {
	sbPlayer := g.tm.NextPlayer(g.gs.ButtonPosition, PlayerStatus_Playing)
	bbPlayer := g.tm.NextPlayer(sbPlayer.Position(), PlayerStatus_Playing)

	if sbPlayer == nil || bbPlayer == nil {
		// Log error when taking blinds
		fmt.Println("error: Can not take blinds")
		return
	}

	sbPlayer.TakeChips(int(g.setting.SmallBlind))
	sbPlayer.UpdateCurrentBet(int(g.setting.SmallBlind))

	bbPlayer.TakeChips(int(g.setting.BigBlind))
	bbPlayer.UpdateCurrentBet(int(g.setting.BigBlind))

	// Update the current bet
	g.gs.CurrentBet = int(g.setting.BigBlind)

	// Add the blinds to the pot
	g.gs.pot.AddToPot(sbPlayer.Position(), sbPlayer.CurrentBet())
	g.gs.pot.AddToPot(bbPlayer.Position(), bbPlayer.CurrentBet())

	// log take blinds from players successfyully
	fmt.Printf("Small blind %s takes %d chips\n", sbPlayer.Name(), int(g.setting.SmallBlind))
	fmt.Printf("Big blind %s takes %d chips\n", bbPlayer.Name(), int(g.setting.BigBlind))

	// Update the next active player
	np := g.tm.NextPlayer(bbPlayer.Position(), PlayerStatus_Playing)
	np.UpdateStatus(PlayerStatus_Wait4Act)
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
	case msgpb.RoundStateType_FLOP:
		// Burn a card
		_ = g.deck.Draw()
		// Add 3 cards to the community cards
		g.gs.cc.AddCard(g.deck.Draw())
		g.gs.cc.AddCard(g.deck.Draw())
		g.gs.cc.AddCard(g.deck.Draw())

		// Print the community cards at flop
		fmt.Printf("============ BOARD at FLOP ===========\n%s\n======================================\n", g.gs.cc.String())
	case msgpb.RoundStateType_TURN:
		// Burn a card
		_ = g.deck.Draw()
		// Add a card to the community cards
		g.gs.cc.AddCard(g.deck.Draw())

		// Print the community cards at turn
		fmt.Printf("============ BOARD at TURN ===========\n%s\n======================================\n", g.gs.cc.String())
	case msgpb.RoundStateType_RIVER:
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
	np := g.tm.NextPlayer(g.gs.ButtonPosition, PlayerStatus_Playing)

	if np != nil {
		np.UpdateStatus(PlayerStatus_Wait4Act)
		np.UpdateSuggestions([]PlayerActType{Check, Raise, AllIn})
		return true
	}

	return false
}

func (g *Game) handleEnterNewRoundLogic() {
	switch g.gs.CurrentRound {
	case msgpb.RoundStateType_INITIAL:
		g.prepareForIncomingGame()
	case msgpb.RoundStateType_PREFLOP:
		g.prepareForNewBettingRound()
		g.takeBlinds()
		g.dealCardsToPlayers()
		// Next player to act is the player next to the big blind
		g.gs.CurrentRound = msgpb.RoundStateType_FLOP
	case msgpb.RoundStateType_FLOP:
		g.prepareForNewBettingRound()
		g.dealCommunityCards()
		// First player to act is the player next to the dealer
		if !g.firstPlayerActionInRound() {
			g.dealTheRestOfCommunityCards()
			g.evaluateHands()
			g.prepareForIncomingGame()
		} else {
			g.gs.CurrentRound = msgpb.RoundStateType_TURN
		}
	case msgpb.RoundStateType_TURN:
		g.prepareForNewBettingRound()
		g.dealCommunityCards()
		// First player to act is the player next to the dealer
		if !g.firstPlayerActionInRound() {
			g.dealTheRestOfCommunityCards()
			g.evaluateHands()
			g.prepareForIncomingGame()
		} else {
			g.gs.CurrentRound = msgpb.RoundStateType_RIVER
		}
	case msgpb.RoundStateType_RIVER:
		g.prepareForNewBettingRound()
		g.dealCommunityCards()

		if !g.firstPlayerActionInRound() {
			g.dealTheRestOfCommunityCards()
			g.evaluateHands()
			g.prepareForIncomingGame()
		}
		g.gs.CurrentRound = msgpb.RoundStateType_SHOW_DOWN
	case msgpb.RoundStateType_SHOW_DOWN:
		// Evaluate hands to find the winner for main pot and side pot
		g.evaluateHands()
		g.prepareForIncomingGame()
	}
}

func (g *Game) evaluateHands() {
	if g.gs.NumPlayingPlayer == 1 {
		// Log the winner
		onePlayer := g.tm.GetListOfPlayers(PlayerStatus_Playing, PlayerStatus_Call, PlayerStatus_Check, PlayerStatus_Raise, PlayerStatus_AllIn)
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
			if p != nil && p.Status() != PlayerStatus_Fold {
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
