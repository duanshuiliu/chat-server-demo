package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"github.com/mitchellh/mapstructure"

	logger   "chat/pkg/logger"
	protocol "chat/internal/conn/websocket/protocol"
)

type ChatController struct {
	BaseController
}

func (this *ChatController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO 登录验证

	conn, err := this.Conn(w, r)
	if err != nil {
		logger.NewLogger().Error("websocket建立连接失败", zap.String("info", err.Error()))
		return
	}
	//defer conn.Close()

	token  := r.Header.Get("Sec-WebSocket-Protocol")
	tokenS := strings.Split(token, "#")
	if len(tokenS) != 2 { return }

	uid, _ := strconv.Atoi(tokenS[0])
	name   := tokenS[1]

	client := &Client{
		Id         : uid,
		Type       : 1,
		Name       : name,
		Conn       : conn,
		ReceiveChan: make(chan *protocol.Message, 10000),
		Status     : 1,
	}

	// 注册客户端
	WebsocketClients.RegisterChan <- client
	defer func() {
		WebsocketClients.LeaveChan <- client
	}()

	// 处理消息
	go client.Run()

	for {
		// 读取消息
		msg, err := client.ReadMessage()
		if err != nil {
			break
		}

		// 解析数据
		var m protocol.Message
		err = json.Unmarshal(msg, &m)
		if err != nil {
			logger.NewLogger().Info("Websocket客户端 - 获取消息解析错误", zap.String("info", err.Error()))
			continue
		}

		switch m.Type {
		// 私聊信息
		case 10001:
			var privateMsg protocol.PrivateMsg
			err := mapstructure.Decode(m.Data, &privateMsg)
			if err != nil {
				logger.NewLogger().Error("Websocket客户端 - 获取消息解析错误", zap.String("info", err.Error()))
				continue
			}

			receiver, loaded := WebsocketClients.Pool.Load(privateMsg.ReceiverId)
			if loaded {
				receiverClient, ok := receiver.(*Client)

				if ok {
					sendMsg := &protocol.Message{
						Type: m.Type,
						Kind: m.Kind,
						Data: make(map[string]interface{}),
					}

					sendMsg.Data = protocol.PackPrivateMsg(client.Id, protocol.MsgTypeSingle, client.Id, client.Name, "", privateMsg.Content)
					receiverClient.ReceiveChan <- sendMsg
				}
			}
		}
	}
}
