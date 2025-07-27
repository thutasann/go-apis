package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func Stream(c *fiber.Ctx) error { return nil }

func StreamWebsocket(c *websocket.Conn) {}

func StreamViewerWebsocket(c *websocket.Conn) {}
