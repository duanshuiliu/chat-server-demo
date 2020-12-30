package controllers

import (
	"encoding/json"
	"sync"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	logger   "chat/pkg/logger"
	protocol "chat/internal/conn/websocket/protocol"
)

var WebsocketClients *ClientsPool

// 客户端连接池
type ClientsPool struct {
	Pool         sync.Map
	LeaveChan    chan *Client
	RegisterChan chan *Client
}

// 连接池的注册与销毁
func (this *ClientsPool) Run() {
	for {
		select {
		case client := <- this.RegisterChan:
			logger.NewLogger().Info("Websocket客户端 - 注册", zap.Int("client", client.Id))

			this.Pool.Store(client.Id, client)

			friendsMsg := &protocol.Message{
				Type: protocol.MsgTypeRegister,
				Kind: protocol.MsgKindText,
				Data: make(map[string]interface{}),
			}

			friendsOfSendSelf := make([]map[string]interface{}, 0)
			userOfSelf        := protocol.PackRegisterMsg(client.Id, protocol.ChatTypeSingle, client.Name, "")

			// 设置一个公用群
			userOfCommon := protocol.PackRegisterMsg(1, protocol.ChatTypeGroup, "群聊", "")
			friendsOfSendSelf = append(friendsOfSendSelf, userOfCommon)

			this.Pool.Range(func(k, v interface{}) bool {
				cli, ok := v.(*Client)

				if ok && cli.Id != client.Id {

					sendMsg := &protocol.Message{
						Type: protocol.MsgTypeRegister,
						Kind: protocol.MsgKindText,
						Data: make(map[string]interface{}),
					}

					friendsOfSendOther := make([]map[string]interface{}, 0)
					userOfOther        := protocol.PackRegisterMsg(cli.Id, protocol.ChatTypeSingle, cli.Name, "")

					friendsOfSendOther = append(friendsOfSendOther, userOfSelf)
					friendsOfSendSelf  = append(friendsOfSendSelf, userOfOther)

					// 通知其它小伙伴
					sendMsg.Data["users"] = friendsOfSendOther
					cli.ReceiveChan <- sendMsg
				}

				return true
			})

			// 获取当前存在的小伙伴列表
			friendsMsg.Data["users"] = friendsOfSendSelf
			client.ReceiveChan <- friendsMsg

		case client := <- this.LeaveChan:
			logger.NewLogger().Info("Websocket客户端 - 销毁", zap.Int("client", client.Id))

			_, loaded := this.Pool.LoadAndDelete(client.Id)
			if loaded {
				logger.NewLogger().Info("Websocket客户端 - 销毁成功", zap.Int("client", client.Id))

				// 关闭客户端资源
				client.Status = 0
				close(client.ReceiveChan)
				err := client.Conn.Close()
				if err != nil {
					logger.NewLogger().Error("Websocket客户端 - 关闭连接失败", zap.Int("client", client.Id), zap.String("info", err.Error()))
				}

				userOfSelf := protocol.PackRegisterMsg(client.Id, protocol.ChatTypeSingle, client.Name, "")

				this.Pool.Range(func(k, v interface{}) bool {
					cli, ok := v.(*Client)

					if ok && cli.Id != client.Id {

						sendMsg := &protocol.Message{
							Type: protocol.MsgTypeLeave,
							Kind: protocol.MsgKindText,
							Data: make(map[string]interface{}),
						}

						friendsOfSendOther := make([]map[string]interface{}, 0)
						friendsOfSendOther = append(friendsOfSendOther, userOfSelf)

						// 通知其它小伙伴
						sendMsg.Data["users"] = friendsOfSendOther
						cli.ReceiveChan <- sendMsg
					}

					return true
				})
			} else {
				logger.NewLogger().Info("Websocket客户端 - 已销毁", zap.Int("client", client.Id))
			}
		}
	}
}

// 客户端
type Client struct {
	Id          int
	Type        int
	Name        string
	Conn        *websocket.Conn
	ReceiveChan chan *protocol.Message
	Status      int
}

// 接收的消息
type ReceiveMsg struct {
	Type int    `json:"type"`
	Kind int    `json:"kind"`
	Data string `json:"data"`
}

// 接收消息
func (this *Client) ReadMessage() (msg []byte, err error) {
	_, msg, err = this.Conn.ReadMessage()

	if err != nil {
		logger.NewLogger().Error("Websocket客户端 - Receive error", zap.String("info", err.Error()))
		WebsocketClients.LeaveChan <- this
	} else {
		logger.NewLogger().Info("Websocket客户端 - Receive", zap.String("msg", string(msg)), zap.Int("clientId", this.Id))
	}

	return
}

// 发送消息
func (this *Client) SendMessage(msg []byte) (err error) {
	err = this.Conn.WriteMessage(websocket.TextMessage, msg)

	if err != nil {
		logger.NewLogger().Error("Websocket客户端 - Send error", zap.String("info", err.Error()))
		WebsocketClients.LeaveChan <- this
	}

	return
}

// 处理消息
func (this *Client) Run() {
	for msg := range this.ReceiveChan {
		sendMsg, err := json.Marshal(msg)
		if err != nil {
			logger.NewLogger().Error("组装发送消息数据错误", zap.String("info", err.Error()), zap.Any("msg", sendMsg))
			break
		}

		err = this.SendMessage(sendMsg)
		if err != nil {
			break
		}
	}
}