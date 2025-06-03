package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// Client struct
type Client struct {
	connection *websocket.Conn // Client connection
	manager    *Manager        // Client manager
	egress     chan Event      // egress is used to avoid concurrent writes on the websocket connection
}

// Connected Client List
type ClientList map[*Client]bool

// Initialize New Client
func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
	}
}

// Read Client Mesasges
func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		_, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
				log.Printf("🔴 error reading message: %v", err)
			}
			break
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("🔴 error marshalling event: %v", err)
			break
		}

		if err := c.manager.routeEvent(request, c); err != nil {
			log.Println("🔴 error handling message: ", err)
		}

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
