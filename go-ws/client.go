package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client struct
type Client struct {
	connection *websocket.Conn // Client connection
	manager    *Manager        // Client manager
	egress     chan []byte     // egress is used to avoid concurrent writes on the websocket connection
}

// Connected Client List
type ClientList map[*Client]bool

// Initialize New Client
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan []byte),
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

		for wsClient := range c.manager.clients {
			wsClient.egress <- payload
		}

		log.Println("messageType ==> ", messageType)
		log.Println("payload ==> ", string(payload))
	}
}

// Write Messages
func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("🔴 connection closed: ", err)
				}
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("🔴 failed to send mesasge: %v", err)
				return
			}

			log.Println("✅ message sent")
		}
	}
}
