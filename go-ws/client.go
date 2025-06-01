package main

import "github.com/gorilla/websocket"

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
