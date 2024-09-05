package network

import (
	core "go-pk-server/core"
	msgpb "go-pk-server/gen"
	mylog "go-pk-server/log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Username string
	GroupId  uint64
	player   *core.OnlinePlayer
	ws       *websocket.Conn
}

func newConnectedClient(username string, gId uint64, agent core.Agent, ws *websocket.Conn) *Client {
	return &Client{
		Username: username,
		GroupId:  gId,
		player:   core.NewOnlinePlayer(username, agent, gId),
		ws:       ws,
	}
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
		default:
			mylog.Errorf("Server not support control message type: %v", x.ControlMessage)
		}

	default:
		mylog.Error("Unknown message type.")
	}
}

func (c *Client) handleDisconnect() {
	mylog.Infof("Client player %s left the game", c.Username)
	//core.MyGame.PlayerLeave(c.player)
}
