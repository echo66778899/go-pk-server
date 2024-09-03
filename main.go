package main

import (
	core "go-pk-server/core"
	mylog "go-pk-server/log"
	snetwork "go-pk-server/network"
)

func main() {
	core.MyGame.StartEngine(true)

	mylog.Infof("Starting server on :%d", 8080)

	connections := snetwork.NewConnectionManager()
	connections.CreateRoom(2222)
	connections.CreateRoom(3333)

	err := connections.StartServer(":8080")
	if err != nil {
		mylog.Error("Failed to start server:", err)
	}
}
