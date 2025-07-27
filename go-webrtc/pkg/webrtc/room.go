package webrtc

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

var (
	RoomsLock sync.RWMutex
	Rooms     map[string]*Room
	Streams   map[string]*Room
)

func RoomConn(c *websocket.Conn, p *Peers) {}
