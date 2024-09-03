package network

import (
	core "go-pk-server/core"
	mylog "go-pk-server/log"
	"go-pk-server/msg"
	"strconv"

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

func (c *Client) handleMessage(message msg.CommunicationMessage) {
	if message.Payload == nil {
		return
	}
	msgPayload, ok := message.Payload.(map[string]interface{})
	if !ok {
		return
	}

	switch message.Type {
	case msg.PlayerActMsgType:
		defautl := core.Fold
		value := 0
		actionName := msgPayload["action_name"].(string)
		switch actionName {
		case "call":
			defautl = core.Call
		case "check":
			defautl = core.Check
		case "fold":
			defautl = core.Fold
		case "raise":
			defautl = core.Raise
			value = int(msgPayload["value"].(float64))
		case "allin":
			defautl = core.AllIn
		default:
			mylog.Errorf("Server not support player action type: %v", actionName)
			return
		}
		mylog.Infof("Client player %s sent player action: %v", c.Username, actionName)
		core.MyGame.PlayerAction(c.player.NewReAct(defautl, value))

	case msg.CtrlMsgType:
		ctrlType := msgPayload["control_type"].(string)
		switch ctrlType {
		case "join_slot":
			slot := msgPayload["data"].(string)
			// convert to int
			slotNo, err := strconv.Atoi(slot)
			if err != nil {
				mylog.Errorf("Failed to convert slot number to int: %v", err)
				return
			}
			mylog.Infof("Client player %s joined the game", c.Username)
			c.player.UpdatePosition(slotNo)
			core.MyGame.PlayerJoin(c.player)
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
			mylog.Errorf("Server not support control message type: %v", ctrlType)
		}
	default:
		mylog.Errorf("Server not support message type: %v", message)
	}
}

func (c *Client) handleDisconnect() {
	mylog.Infof("Client player %s left the game", c.Username)
	//core.MyGame.PlayerLeave(c.player)
}
