package handler

import (
	"go-pk-server/client/ui"
	msgpb "go-pk-server/gen"
	"strconv"

	"log"
)

type Handler struct {
	UI_Model *ui.Model
}

func (h *Handler) HandleServerMessage(message *msgpb.ServerMessage) {
	if message == nil {
		return
	}

	switch x := message.GetMessage().(type) {
	case *msgpb.ServerMessage_GameState:
		handleGameState(message.GetGameState())
	case *msgpb.ServerMessage_GameSetting:
		handleGameSetting(message.GetGameSetting())
	case *msgpb.ServerMessage_PeerState:
		handlePeerState(message.GetPeerState())
	case *msgpb.ServerMessage_BalanceInfo:
		handleBalanceInfo(message.GetBalanceInfo())
	case *msgpb.ServerMessage_ErrorMessage:
		handleErrorMessage(message.GetErrorMessage())
	default:
		log.Fatalf("Unknown message type: %v", x)
	}
}

func handleGameState(gs *msgpb.GameState) {
	log.Printf("Game State: %+v\n", gs)
	ui.UI_MODEL_DATA.Players = gs.Players
	ui.UI_MODEL_DATA.DealerPosition = int(gs.DealerId)
	ui.UI_MODEL_DATA.CurrentRound = gs.CurrentRound
	ui.UI_MODEL_DATA.CurrentBet = int(gs.CurrentBet)

	ui.UI_MODEL_DATA.YourPlayerState = nil
	// Check if username is in the list of players
	for _, player := range gs.Players {
		if player.Name == ui.UI_MODEL_DATA.YourLoginUsernameID {
			ui.UI_MODEL_DATA.YourTablePosition = int(player.TablePosition)
			ui.UI_MODEL_DATA.YourPlayerState = player
			break
		}
	}
	// Enable Slot buttons when the username is not in the list of players
	if ui.UI_MODEL_DATA.YourPlayerState == nil || len(ui.UI_MODEL_DATA.Players) == 0 {
		ui.UI_MODEL_DATA.ActiveButtonMenu = ui.ButtonMenuType_SLOTS_BTN
		ui.UI_MODEL_DATA.IsButtonEnabled = true
		ui.UI_MODEL_DATA.IsButtonsVisible = true
		// Reset the YourTablePosition
		ui.UI_MODEL_DATA.YourTablePosition = 0
	}

	// Enable control buttons when the username is in the list of players
	if ui.UI_MODEL_DATA.YourPlayerState != nil {
		ui.UI_MODEL_DATA.IsButtonsVisible = true
		switch gs.CurrentRound {
		case msgpb.RoundStateType_INITIAL:
			ui.UI_MODEL_DATA.ActiveButtonMenu = ui.ButtonMenuType_CTRL_BTN
			ui.UI_MODEL_DATA.IsButtonEnabled = true
		default:
			ui.UI_MODEL_DATA.ActiveButtonMenu = ui.ButtonMenuType_PLAYING_BTN
			if ui.UI_MODEL_DATA.YourPlayerState.Status == msgpb.PlayerStatusType_Wait4Act {
				ui.UI_MODEL_DATA.IsButtonEnabled = true
			} else if ui.UI_MODEL_DATA.YourPlayerState.Status == msgpb.PlayerStatusType_Spectating {
				ui.UI_MODEL_DATA.IsButtonsVisible = false
			} else {
				ui.UI_MODEL_DATA.IsButtonEnabled = false
			}
		}
	}

	if gs.CommunityCards != nil {
		ui.UI_MODEL_DATA.CommunityCards = gs.CommunityCards
	}

	// When current round is initial, clear your private state
	switch gs.CurrentRound {
	case msgpb.RoundStateType_INITIAL:
		ui.UI_MODEL_DATA.IsDealerVisible = true
		ui.UI_MODEL_DATA.YourPrivateState = nil
		ui.UI_MODEL_DATA.Pot = 0
		ui.UI_MODEL_DATA.CommunityCards = nil
		ui.UI_MODEL_DATA.Result = nil
	case msgpb.RoundStateType_PREFLOP:
		ui.UI_MODEL_DATA.CommunityCards = nil
		ui.UI_MODEL_DATA.Result = nil
	case msgpb.RoundStateType_SHOWDOWN:
		ui.UI_MODEL_DATA.IsDealerVisible = false
		ui.UI_MODEL_DATA.Result = gs.FinalResult
		ui.UI_MODEL_DATA.Pot = int(gs.PotSize)
	default:
		ui.UI_MODEL_DATA.IsDealerVisible = true
		ui.UI_MODEL_DATA.Pot = int(gs.PotSize)
	}

	// Update UI state based on Your Player status
	if ui.UI_MODEL_DATA.YourPlayerState != nil {
		switch ui.UI_MODEL_DATA.YourPlayerState.Status {
		case msgpb.PlayerStatusType_Fold:
			ui.UI_MODEL_DATA.YourTurn = true
			ui.UI_MODEL_DATA.YourPrivateState = nil
		default:
		}
	}

	// Show dealer icon when there are players
	ui.UI_MODEL_DATA.IsDealerVisible = (len(ui.UI_MODEL_DATA.Players) > 0)
}

func handleGameSetting(gs *msgpb.GameSetting) {
	log.Printf("Game Setting: %+v\n", gs)
	ui.UI_MODEL_DATA.MaxPlayers = int(gs.MaxPlayers)
}

func handlePeerState(ps *msgpb.PeerState) {
	log.Printf("Peer State: %+v\n", ps)
	if ps.TablePos == int32(ui.UI_MODEL_DATA.YourTablePosition) {
		ui.UI_MODEL_DATA.YourPrivateState = ps
	}
}

func handleBalanceInfo(bi *msgpb.BalanceInfo) {
	log.Printf("Balance Info: %+v\n", bi)

	ui.UI_MODEL_DATA.PlayersBalance = make([]string, len(bi.PlayerBalances))

	for i, balance := range bi.PlayerBalances {
		if balance != nil {
			ui.UI_MODEL_DATA.PlayersBalance[i] = "  [ " +
				ui.FixStringLength(balance.PlayerName, 10, ' ') +
				" ] : " +
				strconv.Itoa(int(balance.Balance))
		}
	}
}

func handleErrorMessage(em string) {
	log.Printf("Error Message: %+v\n", em)
}
