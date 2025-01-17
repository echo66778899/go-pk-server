package main

import (
	core "go-pk-server/core"
	msgpb "go-pk-server/gen"
	mylog "go-pk-server/log"
	snetwork "go-pk-server/network"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		mylog.Fatal("Usage: go run main.go <host:port>")
	}

	// Get the first argument as address
	address := os.Args[1]

	core.MyGame.StartEngine(true)
	mylog.Infof("Starting server on :%d", 8080)

	room2222 := snetwork.NewRoom(2, "1", "Hai Phan")
	room2222.Serve()
	room2222.SetSettingGetter(
		func() *msgpb.GameSetting {
			return core.MyGame.GetGameSetting()
		},
		func() *msgpb.GameState {
			return core.MyGame.GetGameState()
		},
	)

	core.MyGame.SetRoomAgent(room2222)

	connections := snetwork.NewConnectionManager()
	connections.AddRoom(2, room2222)

	err := connections.StartServer(address)
	if err != nil {
		mylog.Error("Failed to start server:", err)
	}
}
