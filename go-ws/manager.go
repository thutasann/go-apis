package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Manager Struct
type Manager struct {
	clients ClientList // Manager's client list
	sync.RWMutex
}

// Initialize new Manager
func NewManager() *Manager {
	return &Manager{
		clients: make(ClientList),
	}
}

// Serve Websocket
func (m *Manager) serverWS(w http.ResponseWriter, r *http.Request) {
	log.Println("::: WS New Connection :::")

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := NewClient(conn, m)
	m.addClient(client)

	// Start client processes
	go client.readMessages()
}

// Add Client
func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[client] = true
}

// Remove Client
func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}
