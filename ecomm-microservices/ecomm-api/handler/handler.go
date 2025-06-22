package handler

import (
	"context"

	"github.com/dhij/ecomm/ecomm-api/server"
)

type handler struct {
	ctx    context.Context
	server *server.Server
}

func NewHandler(server *server.Server) *handler {
	return &handler{
		ctx:    context.Background(),
		server: server,
	}
}
