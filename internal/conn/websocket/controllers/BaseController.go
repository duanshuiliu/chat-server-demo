package controllers

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type BaseController struct {
	//
}

func (this *BaseController) Conn(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	responseHeader := make(http.Header)
	responseHeader.Add("Sec-WebSocket-Protocol", r.Header.Get("Sec-WebSocket-Protocol"))
	return upgrader.Upgrade(w, r, responseHeader)
}

