package handler

import (
	msgpb "go-pk-server/gen"
	"log"
)

func HandleServerMessage(message *msgpb.ServerMessage) {
	if message == nil {
		return
	}
	switch x := message.GetMessage().(type) {
	case *msgpb.ServerMessage_GameState:
		handleGameState(message.GetGameState())
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
}

func handlePeerState(ps *msgpb.PeerState) {
	log.Printf("Peer State: %+v\n", ps)
}

func handleErrorMessage(em string) {
	log.Printf("Error Message: %+v\n", em)
}
