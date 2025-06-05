package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second    // pong wait time
	pingInterval = (pongWait * 9) / 10 // ping interval
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

	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}

	c.connection.SetPongHandler(c.pongHandler)

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

	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("🔴 connection closed: ", err)
				}
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println("🔴 marshal error fro message: ", err)
				return
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("🔴 failed to send mesasge: %v", err)
				return
			}

			log.Println("✅ message sent")

		case <-ticker.C:
			log.Println("PING...")

			// send PING to the client
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte(`PONG`)); err != nil {
				log.Println("🔴 send PING to the client writeMessages err: ", err)
				return
			}
		}
	}
}

// Private: pong handler
func (c *Client) pongHandler(msg string) error {
	log.Println("pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}
