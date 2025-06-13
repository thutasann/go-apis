package ws

import "github.com/gorilla/websocket"

type Client struct {
	Conn     *websocket.Conn
	Mesasge  chan *Messasge
	ID       string `json:"id"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}

type Messasge struct {
	Content  string `json:"content"`
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
}
