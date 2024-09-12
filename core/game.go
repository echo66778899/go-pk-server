package engine

// Define game rule machanism

import (
	msgpb "go-pk-server/gen"
	mylog "go-pk-server/log"
	"strings"
	"time"
)

// Interface to manage the table
type TableManager interface {
	GetMaxSlot() int
	UpdateMaxSeat(int)
	CountSeatedPlayers() int
	CountPlayablePlayers() int
	CheckPlayersReadiness(*msgpb.GameSetting) bool
	DetermineNextButtonPosition(int) int
	GetSmallBlindPosition() int
	GetBigBlindPosition() int
	DealCardsToPlayers(*Deck)
	FindNextPlayablePlayer(int, map[msgpb.PlayerStatusType]bool) (Player, bool)
	PrepareForNewRound()
	PrepareNewGame()
	GetPlayerAtPosition(int) (Player, bool)
	FindLastStayingPlayer() (Player, bool)
	UpdatePlayerStatusDueToCurrentBetIncrease(int)
	DoAttachedFunctionToAllPlayers(func(Player))
}

type GameStatistcs struct {
	TotalHandsPlayed int
}

type Game struct {
	GameStatistcs
	setting *msgpb.GameSetting
	gs      GameState
	tm      TableManager
	deck    *Deck
	auto    *AutoInputProducer

	funcReqEngineState func(EngineState, string)
}

// NewPokerGame creates a new PokerGame
func NewGame(setting *msgpb.GameSetting, tm TableManager, d *Deck, auto *AutoInputProducer, reqEState func(EngineState, string)) *Game {
	tm.UpdateMaxSeat(int(setting.MaxPlayers))
	return &Game{
		setting: setting,
		gs: GameState{
			pot: NewPot(),
			cc:  CommunityCards{},
		},
		tm:                 tm,
		deck:               d,
		auto:               auto,
		funcReqEngineState: reqEState,
	}
}

func (g *Game) Play() bool {
	mylog.Infof("Play when current round is %s \n", g.gs.CurrentRound)
	// Check if the number of players is valid
	if g.tm.CountSeatedPlayers() < 2 {
		// Log error when the number of players is less than 2
		mylog.Error("Number of players is less than 2")
		return false
	}

	mylog.Infof("Checking if can start a new game\n")
	if !g.isPlayersReadyToPlay() {
		// Log error when the player is not ready
		mylog.Error("Some players are not ready to Play")
		return false
	}

	// How many players are sitting in
	mylog.Infof("Starting a NEW_GAME with %d players in the table", g.tm.CountSeatedPlayers())
	// How many players are playing
	g.gs.readyPlayersCount = g.tm.CountPlayablePlayers()
	mylog.Infof("[New game] The number of playing players: %d\n", g.gs.readyPlayersCount)

	// Start the first round
	g.handleCurrentRoundIsOver()

	return true
}

func (g *Game) HandleEndGame() bool {
	mylog.Infof("Handling end game when current round is %s \n", g.gs.CurrentRound)
	// Check if the game is over
	if g.gs.CurrentRound != msgpb.RoundStateType_SHOWDOWN {
		// Log error when the game is not over
		mylog.Error("Game is not over")
		return false
	}

	mylog.Infof("Checking if can continue to play the next game\n")
	// First check if can play
	if !g.isPlayersReadyToPlay() {
		// Log error when the player is not ready
		mylog.Error("Some players are not ready to Play. Change to INITIAL state")
		g.ResetGame(false)
		return false
	}

	g.gs.readyPlayersCount = g.tm.CountPlayablePlayers()
	if g.gs.readyPlayersCount < 2 {
		mylog.Warnf("Number of players is less than 2. Change to INITIAL state\n")
		g.ResetGame(false)
		return false
	}

	mylog.Infof("[Next game] The number of playing players: %d\n", g.gs.readyPlayersCount)
	g.handleCurrentRoundIsOver()
	return true
}

func (g *Game) HandlePlayerLeaveDuringTheGame(leftPos int) {
	next, ok := g.tm.FindNextPlayablePlayer(leftPos, map[msgpb.PlayerStatusType]bool{
		msgpb.PlayerStatusType_Playing: true,
		msgpb.PlayerStatusType_SB:      true,
		msgpb.PlayerStatusType_BB:      true,
	})

	g.gs.readyPlayersCount--

	if !ok {
		mylog.Warn("Can not find the next player, the round is over")
		g.handleCurrentRoundIsOver()
		return
	}

	// Treat the left player as folded
	g.determineNextStepAfterProcessingPlayerAction(next.Position(), msgpb.PlayerGameActionType_FOLD)
}

// Call when the game can not be continued
func (g *Game) ResetGame(updateDealer bool) {
	mylog.Infof("Resetting the game when current round is %s \n", g.gs.CurrentRound)

	// Check if the pot is empty
	if g.gs.pot.Total() > 0 {
		// Log warning the pot is not empty
		mylog.Warnf("Pot is not empty: %d\n", g.gs.pot.Total())
		// Return the chips to the players
		number := g.tm.CountSeatedPlayers()
		g.tm.DoAttachedFunctionToAllPlayers(func(p Player) {
			p.AddChips(g.gs.pot.Total() / number)
			p.UpdateCurrentBet(0)
		})
	}

	// Statistics
	g.TotalHandsPlayed++

	mylog.Info("Resetting the game")
	g.tm.PrepareNewGame()
	g.gs.pot.ResetPot()
	g.gs.cc.Reset()
	g.gs.CurrentRound = msgpb.RoundStateType_INITIAL
	g.gs.FinalResult = nil

	g.funcReqEngineState(EngineState_WAIT_FOR_PLAYING, "Game was reset")
}

func (g *Game) ShowPlayerHand(pos int) bool {
	if g.gs.CurrentRound != msgpb.RoundStateType_SHOWDOWN {
		// Log error when the game is not over
		mylog.Errorf("Request show hand when Game is not over. Current round is %s \n", g.gs.CurrentRound)
		return false
	}

	// find the player at
	p, ok := g.tm.GetPlayerAtPosition(pos)

	if !ok {
		// Log error when the player is not found
		mylog.Errorf("Requested show hand player not found at position %d\n", pos)
		return false
	}

	if !p.HasPocketCards() {
		// Log error when the player does not have pocket cards
		mylog.Errorf("Player %s does not have pocket cards\n", p.Name())
		return false
	}

	// Log the player hand
	mylog.Infof("Player %s want to show his/her hand: %s\n", p.Name(), p.ShowHand().String())

	if g.gs.FinalResult == nil {
		g.gs.FinalResult = &msgpb.Result{
			ShowingCards: []*msgpb.PeerState{
				{
					TablePos:      int32(p.Position()),
					PlayerCards:   p.ShowHand().Cards(),
					HandRanking:   p.ShowHand().GetPlayerHandRanking(),
					EvaluatedHand: p.ShowHand().BestHand(),
				},
			},
		}
	} else {
		g.gs.FinalResult.ShowingCards = append(g.gs.FinalResult.ShowingCards, &msgpb.PeerState{
			TablePos:      int32(p.Position()),
			PlayerCards:   p.ShowHand().Cards(),
			HandRanking:   p.ShowHand().GetPlayerHandRanking(),
			EvaluatedHand: p.ShowHand().BestHand(),
		})
	}

	return true
}

func (g *Game) HandleActions(action ActionIf) {
	mylog.Infof("Handling action when current round is %s \n", g.gs.CurrentRound)
	if action == nil {
		return
	}

	// Check if the action is valid
	player, ok := g.tm.GetPlayerAtPosition(action.AtPosition())
	if !ok {
		mylog.Errorf("Action from a not found player (%d)\n", action.AtPosition())
		return
	}

	if player.Status() != msgpb.PlayerStatusType_Wait4Act {
		// Log warning the player is not allowed to fold, the action is invalid
		mylog.Errorf("Player %s is not allowed to %d due to unexpected action", player.Name(), action.WhatAction())
		return
	}

	mylog.Debugf("Handling %v action from player %s.\n", action.WhatAction(), player.Name())
	mylog.Debugf("BEFORE Current bet: %d\n", g.gs.CurrentBet)

	switch action.WhatAction() {
	case msgpb.PlayerGameActionType_FOLD:
		// Execute fold action
		player.UpdateStatus(msgpb.PlayerStatusType_Fold)
		player.UpdateCurrentBet(0)
		// Do not drop pocket cards because we want all players can show the hand at the end
		//player.DropPocketCards()
		// Decrease the number of playing players
		g.gs.readyPlayersCount--

	case msgpb.PlayerGameActionType_CHECK:
		// Execute check action
		// Check if the player is allowed to check
		if player.CurrentBet() == g.gs.CurrentBet {
			player.UpdateStatus(msgpb.PlayerStatusType_Check)
		} else {
			// Log info the player name is not allowed to check, the action is invalid
			mylog.Errorf("Player %s is not allowed to check, the action is invalid", player.Name())
			return
		}
	case msgpb.PlayerGameActionType_CALL:
		callChip := g.gs.CurrentBet - player.CurrentBet()
		// If the player chip is less than the current bet, the player is all-in
		if player.Chips() <= callChip {
			player.UpdateCurrentBet(player.Chips() + player.CurrentBet())
			player.GetChipForBet(player.Chips())
			player.UpdateStatus(msgpb.PlayerStatusType_AllIn)
		} else {
			player.UpdateCurrentBet(callChip + player.CurrentBet())
			player.GetChipForBet(callChip)
			player.UpdateStatus(msgpb.PlayerStatusType_Call)
		}
		g.gs.pot.AddToPot(player.Position(), callChip)
		// call does not change the current bet
	case msgpb.PlayerGameActionType_RAISE:
		// Execute raise action
		raiseAmount := action.HowMuch()
		// log the raise amount
		mylog.Infof("Player %s raises %d\n", player.Name(), raiseAmount)
		callAmount := g.gs.CurrentBet - player.CurrentBet()

		if raiseAmount+callAmount < player.Chips() {
			raiseAmount += callAmount
			player.UpdateCurrentBet(raiseAmount + player.CurrentBet())
			player.GetChipForBet(raiseAmount)
			player.UpdateStatus(msgpb.PlayerStatusType_Raise)
			g.gs.pot.AddToPot(player.Position(), raiseAmount)
		} else {
			raiseAmount := player.Chips()
			player.UpdateCurrentBet(raiseAmount + player.CurrentBet())
			player.GetChipForBet(raiseAmount)
			player.UpdateStatus(msgpb.PlayerStatusType_AllIn)
		}
		g.gs.CurrentBet = player.CurrentBet()
		// Update all player status to msgpb.PlayerStatusType_Playing
		g.tm.UpdatePlayerStatusDueToCurrentBetIncrease(player.Position())
	case msgpb.PlayerGameActionType_ALLIN:
		// Execute all-in action
		allInAmount := player.Chips()
		mylog.Infof("Player %s goes all-in with %d\n", player.Name(), allInAmount)
		player.UpdateCurrentBet(allInAmount + player.CurrentBet())
		player.GetChipForBet(allInAmount)
		player.UpdateStatus(msgpb.PlayerStatusType_AllIn)
		g.gs.pot.AddToPot(player.Position(), allInAmount)
		if player.CurrentBet() > g.gs.CurrentBet {
			g.gs.CurrentBet = player.CurrentBet()
		}

		// Update all player status to msgpb.PlayerStatusType_Playing
		g.tm.UpdatePlayerStatusDueToCurrentBetIncrease(player.Position())
	default:
		// Log invalid action
		mylog.Errorf("Invalid player action: %s\n", action.WhatAction())
		return
	}

	mylog.Debugf("AFTER Current bet: %d, Number of player in-game: (%d)\n", g.gs.CurrentBet, g.gs.readyPlayersCount)

	g.determineNextStepAfterProcessingPlayerAction(player.Position(), action.WhatAction())
}

func (g *Game) determineNextStepAfterProcessingPlayerAction(pPos int, actType msgpb.PlayerGameActionType) {
	if g.gs.readyPlayersCount <= 1 {
		// Log only 1
		mylog.Infof("Only one last player in the game! Gane is going to be over")
		g.evaluateHandsAndUpdateResult()
		g.gs.CurrentRound = msgpb.RoundStateType_SHOWDOWN
		return
	} else if np, _ := g.tm.FindNextPlayablePlayer(pPos, map[msgpb.PlayerStatusType]bool{ // and first found to msgpb.PlayerStatusType_Wait4Act
		msgpb.PlayerStatusType_Playing: true,
		msgpb.PlayerStatusType_SB:      true,
		msgpb.PlayerStatusType_BB:      true}); np != nil {
		// Update the next player status
		np.UpdateStatus(msgpb.PlayerStatusType_Wait4Act)
		switch actType {
		case msgpb.PlayerGameActionType_CHECK:
			if g.gs.CurrentBet == 0 || g.gs.CurrentBet == np.CurrentBet() {
				// Todo: Maybe add FOLD to the list of invalid actions
				np.UpdateInvalidAction([]msgpb.PlayerGameActionType{msgpb.PlayerGameActionType_CALL})
			}
		case msgpb.PlayerGameActionType_FOLD:
			if g.gs.CurrentBet > np.CurrentBet() {
				np.UpdateInvalidAction([]msgpb.PlayerGameActionType{msgpb.PlayerGameActionType_CHECK})
			}
		case msgpb.PlayerGameActionType_CALL, msgpb.PlayerGameActionType_RAISE, msgpb.PlayerGameActionType_ALLIN:
			if np.Chips() < g.gs.CurrentBet*2 {
				np.UpdateInvalidAction([]msgpb.PlayerGameActionType{msgpb.PlayerGameActionType_CHECK, msgpb.PlayerGameActionType_RAISE})
			} else {
				if g.gs.CurrentBet > np.CurrentBet() {
					np.UpdateInvalidAction([]msgpb.PlayerGameActionType{msgpb.PlayerGameActionType_CHECK})
				} else if g.gs.CurrentBet == np.CurrentBet() {
					np.UpdateInvalidAction([]msgpb.PlayerGameActionType{msgpb.PlayerGameActionType_CALL})
				}
			}
		default:
			mylog.Debug("error: Can not suggest action for player", np.Name())
		}
	} else {
		mylog.Warn("Can not find the next player, the round is over")
		g.handleCurrentRoundIsOver()
		return
	}
}

func (g *Game) isPlayersReadyToPlay() (canPlay bool) {
	// Check minimum chips and update player status regarding to the minimum stack size
	return g.tm.CheckPlayersReadiness(g.setting)
}

func (g *Game) prepareForIncomingGame() {
	// Reset the game state for a new game
	g.tm.PrepareNewGame()

	// Shuffle the deck
	g.deck.Shuffle()

	// Reset the game state for a new game
	g.gs.pot.ResetPot()
	g.gs.cc.Reset()
	g.gs.CurrentRound = msgpb.RoundStateType_PREFLOP
	g.gs.FinalResult = nil

	// Choose the dealer position
	g.gs.ButtonPosition = g.tm.DetermineNextButtonPosition(g.gs.ButtonPosition)
	if g.gs.ButtonPosition >= 0 {
		if b, _ := g.tm.GetPlayerAtPosition(g.gs.ButtonPosition); b != nil {
			mylog.Debugf("Game number [%d]. Who is the dealer? -> %s\n", g.TotalHandsPlayed, b.Name())
		}
	} else {
		g.funcReqEngineState(EngineState_WAIT_FOR_PLAYING, "Can not determine the dealer position")
	}
}

func (g *Game) prepareForNewBettingRound() {
	// Reset player state
	g.tm.PrepareForNewRound()
	// Reset new round state
	g.resetGameStateForNewRound()

	// Log current pot when entering new round
	mylog.Infof("Current pot: %d\n", g.gs.pot.Total())
}

func (g *Game) resetGameStateForNewRound() {
	g.gs.CurrentBet = 0
}

func (g *Game) takeBlinds() {
	sbPlayer, _ := g.tm.FindNextPlayablePlayer(g.gs.ButtonPosition,
		map[msgpb.PlayerStatusType]bool{msgpb.PlayerStatusType_SB: true})
	if sbPlayer == nil {
		mylog.Errorf("Can not determine the small blind player\n")
		return
	}
	bbPlayer, _ := g.tm.FindNextPlayablePlayer(sbPlayer.Position(),
		map[msgpb.PlayerStatusType]bool{msgpb.PlayerStatusType_BB: true})

	if bbPlayer == nil {
		mylog.Errorf("Can not determine the big blind player\n")
		return
	}

	sbPlayer.GetChipForBet(int(g.setting.SmallBlind))
	sbPlayer.UpdateCurrentBet(int(g.setting.SmallBlind))

	bbPlayer.GetChipForBet(int(g.setting.BigBlind))
	bbPlayer.UpdateCurrentBet(int(g.setting.BigBlind))

	// Update the current bet
	g.gs.CurrentBet = int(g.setting.BigBlind)

	// Add the blinds to the pot
	g.gs.pot.AddToPot(sbPlayer.Position(), sbPlayer.CurrentBet())
	g.gs.pot.AddToPot(bbPlayer.Position(), bbPlayer.CurrentBet())

	// log take blinds from players successfyully
	mylog.Infof("Small blind %s takes %d chips\n", sbPlayer.Name(), int(g.setting.SmallBlind))
	mylog.Infof("Big blind %s takes %d chips\n", bbPlayer.Name(), int(g.setting.BigBlind))

	// Update the next active player
	np, ok := g.tm.FindNextPlayablePlayer(bbPlayer.Position(),
		map[msgpb.PlayerStatusType]bool{
			msgpb.PlayerStatusType_Playing: true,
			msgpb.PlayerStatusType_SB:      true})
	if ok {
		np.UpdateStatus(msgpb.PlayerStatusType_Wait4Act)
		np.UpdateInvalidAction([]msgpb.PlayerGameActionType{msgpb.PlayerGameActionType_CHECK})
	}
	// np.UpdateSuggestions([]msgpb.PlayerGameActionType{Fold, Call, Raise, AllIn})
}

// Deal the community cards when preflop is over,
// the turn card when when the flop is over,
// the river card when the turn is over.
func (g *Game) dealCommunityCards() {
	switch g.gs.CurrentRound {
	case msgpb.RoundStateType_PREFLOP:
		// Burn a card
		_ = g.deck.Draw()
		// Add 3 cards to the community cards
		g.gs.cc.AddCard(g.deck.Draw())
		g.gs.cc.AddCard(g.deck.Draw())
		g.gs.cc.AddCard(g.deck.Draw())

		// Print the community cards at flop
		mylog.Info("=========================== BOARD at FLOP =============================")
		mylog.Infof("%s\n", g.gs.cc.String())
		mylog.Info("=======================================================================")
	case msgpb.RoundStateType_FLOP:
		// Burn a card
		_ = g.deck.Draw()
		// Add a card to the community cards
		g.gs.cc.AddCard(g.deck.Draw())

		// Print the community cards at turn
		mylog.Info("======================================= BOARD at TURN ==========================================")
		mylog.Infof("%s\n", g.gs.cc.String())
		mylog.Info("================================================================================================")
	case msgpb.RoundStateType_TURN:
		// Burn a card
		_ = g.deck.Draw()
		// Add a card to the community cards
		g.gs.cc.AddCard(g.deck.Draw())

		// Print the community cards at river
		mylog.Info("==================================================== BOARD at RIVER ======================================================")
		mylog.Infof("%s\n", g.gs.cc.String())
		mylog.Info("==========================================================================================================================")
	default:
		// Log error when dealing community cards at wrong round
		mylog.Infof("error: dealing community cards at wrong round\n")
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
		mylog.Debug("Enough cards in the commuuity cards")
	}

	// Print the community cards at river
	mylog.Info("================================================== BOARD at REST =========================================================")
	mylog.Infof("%s\n", g.gs.cc.String())
	mylog.Info("==========================================================================================================================")
}

func (g *Game) firstPlayerActionInRound() bool {
	// If there is only one player, the player wins the pot
	if g.tm.CountPlayablePlayers() == 1 {
		// Log only one player in the game
		mylog.Debug("Only one player in the game!")
		return false
	}

	// First player to act is the player next to the dealer
	np, found := g.tm.FindNextPlayablePlayer(g.gs.ButtonPosition,
		map[msgpb.PlayerStatusType]bool{
			msgpb.PlayerStatusType_Playing: true,
			msgpb.PlayerStatusType_SB:      true,
			msgpb.PlayerStatusType_BB:      true,
		})

	if np != nil {
		np.UpdateStatus(msgpb.PlayerStatusType_Wait4Act)
		np.UpdateInvalidAction([]msgpb.PlayerGameActionType{msgpb.PlayerGameActionType_FOLD, msgpb.PlayerGameActionType_CALL})
		// np.UpdateSuggestions([]msgpb.PlayerGameActionType{Check, Raise, AllIn})
		return found
	}
	return found
}

func (g *Game) handleCurrentRoundIsOver() {
	switch g.gs.CurrentRound {
	case msgpb.RoundStateType_INITIAL:
		g.prepareForIncomingGame()
		// Deal private cards to players
		g.tm.DealCardsToPlayers(g.deck)
		// reset betting state for next preflop round
		g.prepareForNewBettingRound()
		// Take blinds should be done after reset betting state
		g.takeBlinds()
		// Now state is PREFLOP
		g.gs.CurrentRound = msgpb.RoundStateType_PREFLOP
	case msgpb.RoundStateType_PREFLOP:
		// This executes when the preflop round is over
		// reset betting state for next flop round and deal 3 community cards
		g.prepareForNewBettingRound()
		g.dealCommunityCards()

		// If can not find the first player to act, deal the rest of community cards
		// And evaluate hands to find the winner at preflop
		if !g.firstPlayerActionInRound() {
			g.dealTheRestOfCommunityCards()
			g.evaluateHandsAndUpdateResult()
			// Now is in show down
			g.gs.CurrentRound = msgpb.RoundStateType_SHOWDOWN
		} else {
			// Now state is being at FLOP
			g.gs.CurrentRound = msgpb.RoundStateType_FLOP
		}
	case msgpb.RoundStateType_FLOP:
		// This executes when the flop round is over
		// reset betting state for next turn round and deal 1 turn card
		g.prepareForNewBettingRound()
		g.dealCommunityCards()

		// If can not find the first player to act, deal the rest of community cards
		// And evaluate hands to find the winner at flop
		if !g.firstPlayerActionInRound() {
			g.dealTheRestOfCommunityCards()
			g.evaluateHandsAndUpdateResult()
			// Now is in show down
			g.gs.CurrentRound = msgpb.RoundStateType_SHOWDOWN
		} else {
			// Now state is being at TURN
			g.gs.CurrentRound = msgpb.RoundStateType_TURN
		}
	case msgpb.RoundStateType_TURN:
		// This executes when the turn round is over
		// reset betting state for next river round and deal 1 river card
		g.prepareForNewBettingRound()
		g.dealCommunityCards()

		// If can not find the first player to act, deal the rest of community cards
		// And evaluate hands to find the winner at turn
		if !g.firstPlayerActionInRound() {
			g.dealTheRestOfCommunityCards()
			g.evaluateHandsAndUpdateResult()
			// Now is in show down
			g.gs.CurrentRound = msgpb.RoundStateType_SHOWDOWN
		} else {
			// Now state is being at RIVER
			g.gs.CurrentRound = msgpb.RoundStateType_RIVER
		}
	case msgpb.RoundStateType_RIVER:
		// This executes when the turn round is over
		// Evaluate hands to find the winner at end of river round
		g.evaluateHandsAndUpdateResult()
		g.gs.CurrentRound = msgpb.RoundStateType_SHOWDOWN
	case msgpb.RoundStateType_SHOWDOWN:
		// This state indice the game is over
		// We will prepare for the next game
		// Log statistics
		g.TotalHandsPlayed++
		mylog.Infof("Statistics: Total hands played: %d\n", g.TotalHandsPlayed)

		mylog.Info("Continue to play next game from SHOWDOWN state")
		// In case the game is over, continue to the next game
		g.prepareForIncomingGame()
		// Deal private cards to players
		g.tm.DealCardsToPlayers(g.deck)
		// reset betting state for next preflop round
		g.prepareForNewBettingRound()
		// Take blinds should be done after reset betting state
		g.takeBlinds()
		// Now state is PREFLOP
		g.gs.CurrentRound = msgpb.RoundStateType_PREFLOP
	}
}

func (g *Game) evaluateHandsAndUpdateResult() {
	if g.gs.readyPlayersCount == 1 {
		mylog.Info("Only one player in the game! The player wins the pot")
		last, ok := g.tm.FindLastStayingPlayer()
		if ok {
			// Just evaluate the hand of the last player in case they want to show their hand
			last.ShowHand().Evaluate(&g.gs.cc)
			// Log error when finding the last player
			mylog.Infof("Player [%s] wins the pot (%d) WITHOUT SHOWING HAND\n", last.Name(), g.gs.pot.Total())
			last.UpdateStatus(msgpb.PlayerStatusType_WINNER)
			sidePots := g.gs.pot.CalculateSidePots()
			for _, sidePot := range sidePots {
				last.AddWonChips(sidePot.Amount)
			}
		} else {
			panic("Game logic error: Can not find the last player")
		}
	} else {
		// Evaluate hands to find the winner for main pot and side pot
		mylog.Debug("Evaluating hands to determine the winner")
		allEveluateedHands := []*msgpb.PeerState{}

		// First evaluate the player's hands,
		// including the fold players to enable them can show their hands at showdown
		g.tm.DoAttachedFunctionToAllPlayers(func(p Player) {
			if p.HasPocketCards() && p.ShowHand() != nil {
				// Start evaluating the player's hand
				mylog.Debugf("Evaluating player %s: [%s]\n", p.Name(), p.ShowHand().String())
				p.ShowHand().Evaluate(&g.gs.cc)

				// Print its rank
				mylog.Debugf("Player %s's best hand: [%s] >> tiebreaker: %+v >> (%s)\n",
					p.Name(),
					p.ShowHand().BestHandString(),
					p.ShowHand().SortedDecendingRankValue(),
					p.ShowHand().GetPlayerHandRanking())
			}
		})

		// Find the winner for the main pot and side pots
		sidePots := g.gs.pot.CalculateSidePots()
		winners := []Player{}
		for i, sidePot := range sidePots {
			// log info
			mylog.Infof("Find who will win the POT[%d]: %d chips\n", i, sidePot.Amount)
			joinedPlayers := []string{}
			// Find the winner for the side pot
			for _, pos := range sidePot.Players {
				p, ok := g.tm.GetPlayerAtPosition(pos)
				// Add name to the list of players that shared the side pot
				if ok && p.HasPocketCards() && (p.Status() != msgpb.PlayerStatusType_Fold) {
					joinedPlayers = append(joinedPlayers, p.Name())
				} else if ok && (p.Status() == msgpb.PlayerStatusType_Fold) {
					joinedPlayers = append(joinedPlayers, p.Name()+" (Folded)")
				} else {
					joinedPlayers = append(joinedPlayers, " (Exited)")
				}

				if p != nil && p.HasPocketCards() && (p.Status() != msgpb.PlayerStatusType_Fold) {
					// Reset the player current bet for UI display correctly
					p.UpdateCurrentBet(0)

					// If the player is the first player to be evaluated, set it as the winner
					if len(winners) == 0 {
						winners = []Player{p}
						continue
					}

					if p.ShowHand().Compare(winners[0].ShowHand()) > 0 {
						// New winner, clear the list of winners
						winners = []Player{p}
						// Log the new winner with the hand
						mylog.Infof("Player [%s] wins compare with a hand >> [ %s ]\n",
							p.Name(), p.ShowHand().GetPlayerHandRanking())
					} else if p.ShowHand().Compare(winners[0].ShowHand()) == 0 {
						// Add the player to the list of winners
						winners = append(winners, p)
						// Compare the kicker
						mylog.Debugf("Player [%s] and player [%s] equal hand ranking and tiebreaks [ %s == %s ]\n",
							p.Name(), winners[0].Name(),
							p.ShowHand().GetPlayerHandRanking(), winners[0].ShowHand().GetPlayerHandRanking())
					} else {
						// Log this edge case
						mylog.Debugf("Pre winner [%s] still wins the tiebreakers", winners[0].Name())
					}
				}
			}

			// Log the side pot
			if i == 0 {
				mylog.Infof("MAIN pot: %d chips (Players: [%s]). Pot's winner: [%s]\n",
					sidePot.Amount, strings.Join(joinedPlayers, ", "), getWinnersName(winners))
			} else {
				mylog.Infof("Side pot[%d]: %d chips (Players: [%s]). Pot's winner: [%s]\n",
					i, sidePot.Amount, strings.Join(joinedPlayers, ", "), getWinnersName(winners))
			}

			// Distribute the side pot to the winner
			for _, winner := range winners {
				winner.AddWonChips(sidePot.Amount / len(winners))
			}
			// Clear the list of winners
			joinedPlayers, winners = nil, nil
		}

		// Update the player status and add the player to the list of evaluated hands
		g.tm.DoAttachedFunctionToAllPlayers(func(p Player) {
			if p.HasPocketCards() && (p.Status() != msgpb.PlayerStatusType_Fold) {
				if p.ChipChange() >= 0 {
					p.UpdateStatus(msgpb.PlayerStatusType_WINNER)
					mylog.Infof("Player %s WIN the pot (+%d) with a hand >> [ %s ]\n",
						p.Name(), p.ChipChange(), p.ShowHand().GetPlayerHandRanking())
				} else {
					p.UpdateStatus(msgpb.PlayerStatusType_LOSER)
					mylog.Warnf("Player %s LOSE the pot (%d) with a hand >> [ %s ]\n",
						p.Name(), p.ChipChange(), p.ShowHand().GetPlayerHandRanking())
				}
				// Add the player to the list of evaluated hands
				allEveluateedHands = append(allEveluateedHands, &msgpb.PeerState{
					TablePos:      int32(p.Position()),
					PlayerCards:   p.ShowHand().Cards(),
					HandRanking:   p.ShowHand().GetPlayerHandRanking(),
					EvaluatedHand: p.ShowHand().BestHand(),
				})
			}
		})

		// Update the result and showing hands
		g.gs.FinalResult = &msgpb.Result{
			ShowingCards: allEveluateedHands,
		}
	}

	// Auto create the next game input after X setting seconds
	if g.setting.AutoNextGame {
		g.auto.CreatGameInputAfter(GameEnded, time.Duration(g.setting.AutoNextTime)*time.Second)
	}
}
