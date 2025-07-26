package handlers

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func RoomCreate(c *fiber.Ctx) error {
	return nil
}

func Room(c *fiber.Ctx) error {
	return nil
}

func RoomChat(c *fiber.Ctx) error {
	return nil
}

func RoomWebSocket(c *websocket.Conn) {}

func RoomChatWebsocket(c *websocket.Conn) {}

func RoomViewerWebsocket(c *websocket.Conn) {}
