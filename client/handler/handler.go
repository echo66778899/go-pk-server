package handler

import (
	"go-pk-server/client/ui"
	msgpb "go-pk-server/gen"

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

	// Check if username is in the list of players
	for _, player := range gs.Players {
		if player.Name == ui.UI_MODEL_DATA.YourUsernameID {
			ui.UI_MODEL_DATA.YourTablePosition = int(player.TablePosition)
			break
		}
	}
	if gs.CommunityCards != nil {
		ui.UI_MODEL_DATA.CommunityCards = gs.CommunityCards
	}
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

func handleErrorMessage(em string) {
	log.Printf("Error Message: %+v\n", em)
}
