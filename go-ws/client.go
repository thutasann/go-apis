package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client struct
type Client struct {
	connection *websocket.Conn // Client connection
	manager    *Manager        // Client manager
}

// Connected Client List
type ClientList map[*Client]bool

// Initialize New Client
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
	}
}

// Read Client Mesasges
func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		messageType, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
				log.Printf("🔴 error reading message: %v", err)
			}
			break
		}

		log.Println("messageType ==> ", messageType)
		log.Println("payload ==> ", string(payload))
	}
}
