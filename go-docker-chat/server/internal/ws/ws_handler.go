package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WS Handler
type Handler struct {
	hub *Hub
}

// Create Room Request
type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Get Rooms Response
type GetRoomsRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Get Client Response
type GetClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// Websocket Upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Initialize a new WS Handler
func NewHandler(h *Hub) *Handler {
	return &Handler{
		hub: h,
	}
}

// Handle Create Room
func (h *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	c.JSON(http.StatusOK, req)
}

// Handle Join Room
func (h *Handler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID := c.Param("roomId")
	clientID := c.Query("userId")
	username := c.Query("username")

	cl := &Client{
		Conn:     conn,
		Mesasge:  make(chan *Messasge, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	m := &Messasge{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	// Register a new client through the register channel
	h.hub.Register <- cl

	// Broadcast that mesasge
	h.hub.Broadcast <- m

	go cl.writeMesasge()
	cl.readMessage(h.hub)
}

// Handle Get Rooms
func (h *Handler) GetRooms(c *gin.Context) {
	rooms := make([]GetRoomsRes, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, GetRoomsRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	c.JSON(http.StatusOK, rooms)
}

// Handle Get Clients
func (h *Handler) GetClients(c *gin.Context) {
	var clients []GetClientRes
	roomId := c.Param("roomId")

	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients = make([]GetClientRes, 0)
		c.JSON(http.StatusOK, clients)
	}

	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, GetClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}
