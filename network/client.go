package network

import (
	"errors"
	core "go-pk-server/core"
	msgpb "go-pk-server/gen"
	mylog "go-pk-server/log"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	Username string
	GroupId  uint64
	player   *core.OnlinePlayer
	ws       *websocket.Conn
}

func newConnectedClient(username string, gId uint64, roomAgent core.Agent, ws *websocket.Conn) *Client {
	return &Client{
		Username: username,
		GroupId:  gId,
		player:   core.NewOnlinePlayer(username, roomAgent, gId),
		ws:       ws,
	}
}

func (c *Client) send(message *msgpb.ServerMessage) error {
	if message == nil {
		return errors.New("message is nil")
	}
	if c.ws == nil {
		return errors.New("websocket connection not found")
	}
	// Serialize (marshal) the protobuf message
	sendData, err := proto.Marshal(message)
	if err != nil {
		mylog.Fatalf("Failed to marshal proto: %v", err)
		return err
	}
	// Send the response
	if err := c.ws.WriteMessage(websocket.BinaryMessage, sendData); err != nil {
		mylog.Errorf("Failed to write message to client: %v", err)
		return err
	}
	return nil
}

func (c *Client) handleMessage(message *msgpb.ClientMessage) {
	if message == nil {
		return
	}

	switch x := message.GetMessage().(type) {
	case *msgpb.ClientMessage_PlayerAction:
		// Handle the player action
		mylog.Debugf("Player Action: %+v\n", x.PlayerAction)
		switch x.PlayerAction.ActionType {
		case "fold":
			mylog.Infof("Client player %s folded", c.Username)
			core.MyGame.PlayerAction(c.player.NewReAct(core.Fold, 0))
		case "call":
			mylog.Infof("Client player %s called", c.Username)
			core.MyGame.PlayerAction(c.player.NewReAct(core.Call, 0))
		case "check":
			mylog.Infof("Client player %s checked", c.Username)
			core.MyGame.PlayerAction(c.player.NewReAct(core.Check, 0))
		case "raise":
			mylog.Infof("Client player %s raised", c.Username)
			core.MyGame.PlayerAction(c.player.NewReAct(core.Raise, int(x.PlayerAction.RaiseAmount)))
		case "allin":
			mylog.Infof("Client player %s all-in", c.Username)
			core.MyGame.PlayerAction(c.player.NewReAct(core.AllIn, 0))
		default:
			mylog.Errorf("Server not support player action type: %v", x.PlayerAction.ActionType)
		}
	case *msgpb.ClientMessage_JoinGame:
		// Handle the join game request
		mylog.Debugf("Join Game: %+v\n", x.JoinGame)
		c.player.UpdatePosition(int(x.JoinGame.ChooseSlot))
		core.MyGame.PlayerJoin(c.player)
		mylog.Infof("Client player %s joined the game", c.Username)
	case *msgpb.ClientMessage_ControlMessage:
		// Handle the custom control message
		mylog.Debugf("Control Message: %+v\n", x.ControlMessage)
		switch x.ControlMessage {
		case "request_buyin":
			mylog.Infof("Client player %s requested 1 buyin", c.Username)
			c.player.AddChips(1 * 2000)
		case "start_game":
			mylog.Infof("Client player %s started the game", c.Username)
			core.MyGame.StartGame()
		case "next_game":
			mylog.Infof("Client player %s next the game", c.Username)
			core.MyGame.NextGame()
		case "sync_game_state":
			mylog.Infof("Client player %s requested to sync game state", c.Username)
			gsMsg := core.MyGame.SyncGameState()
			c.send(&msgpb.ServerMessage{
				Message: &msgpb.ServerMessage_GameState{
					GameState: gsMsg,
				},
			})
		default:
			mylog.Errorf("Server not support control message type: %v", x.ControlMessage)
		}

	default:
		mylog.Error("Unknown message type.")
	}
}

func (c *Client) handleDisconnect() {
	mylog.Infof("Player %s left the game", c.Username)
	core.MyGame.PlayerLeave(c.player)
}
