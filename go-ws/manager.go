package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		CheckOrigin:     checkOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Manager Struct
type Manager struct {
	clients ClientList // Manager's client list
	sync.RWMutex
	handlers map[string]EventHandler // Event Handlers
}

// Initialize new Manager
func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

// Setup Event Handlers to manager's handlers
func (m *Manager) setupEventHandlers() {
	m.handlers[EventSendMesasge] = SendMesasge
}

// Route Socket Event
func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("there is no such event type")
	}
}

// Send Message Event Handler
func SendMesasge(event Event, c *Client) error {
	fmt.Println("✅ [SendMesasge] event ==> ", event)
	return nil
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
	go client.writeMessages()
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

// check domain origin
func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	switch origin {
	case "http://localhost:8080":
		return true
	default:
		return false
	}
}
