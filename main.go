package main

import (
	core "go-pk-server/core"
	msgpb "go-pk-server/gen"
	mylog "go-pk-server/log"
	snetwork "go-pk-server/network"
)

func main() {
	core.MyGame.StartEngine(true)

	mylog.Infof("Starting server on :%d", 8080)

	room2222 := snetwork.NewRoom(2, "1")
	room2222.Serve()
	room2222.SetSettingGetter(func() *msgpb.GameSetting {
		return core.MyGame.GetGameSetting()
	})

	core.MyGame.SetRoomAgent(room2222)

	connections := snetwork.NewConnectionManager()
	connections.AddRoom(2, room2222)

	err := connections.StartServer(":8080")
	if err != nil {
		mylog.Error("Failed to start server:", err)
	}
}
