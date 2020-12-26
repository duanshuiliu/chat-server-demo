package protocol

import "time"

var (
	MsgTypeRegister int = 10000  // 注册
	MsgTypePrivate  int = 10001  // 私聊

	MsgKindText     int = 1      // 纯文字

	MsgTypeGroup    int = 1      // 群聊
	MsgTypeSingle   int = 2      // 单聊
)

type Message struct {
	Type int `json:"type"`
	Kind int `json:"kind"`
	Data map[string]interface{} `json:"data"`
}

type PrivateMsg struct {
	ReceiverId int    `mapstructure:"receiver_id"`
	Content    string `json:"content"`
}

func PackRegisterMsg(chatId int, chatType int, chatName string, avatar string) map[string]interface{} {
	msg := make(map[string]interface{})

	msg["chat_id"]   = chatId
	msg["chat_type"] = chatType
	msg["chat_name"] = chatName
	msg["avatar"]    = avatar

	return msg
}

func PackPrivateMsg(chatId int, chatType int, senderId int, senderName string, senderAvatar string, content string) map[string]interface{} {
	msg := make(map[string]interface{})

	msg["chat_id"]       = chatId
	msg["chat_type"]     = chatType
	msg["sender_id"]     = senderId
	msg["sender_name"]   = senderName
	msg["sender_avatar"] = senderAvatar
	msg["content"]       = content
	msg["time"]          = time.Now().Format("2006-01-02 15:03:04")

	return msg
}