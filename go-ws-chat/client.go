package main

import (
	"github.com/gorilla/websocket"
)

// client is a single chatting user in a room
type client struct {

	// a websocket for this user
	socket *websocket.Conn

	// receive is a channel to receive messages from other clients
	receive chan []byte

	// chat room
	room *room
}

// send messaegs function
func (c *client) read() {

	// close the connection when we are done
	defer c.socket.Close()

	for {
		_, msg, err := c.socket.ReadMessage()

		if err != nil {
			return
		}

		c.room.forward <- msg
	}
}

// used to received messages
func (c *client) write() {
	defer c.socket.Close()

	for msg := range c.receive {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
