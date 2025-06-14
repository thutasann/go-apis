package main

import (
	"log"

	"github.com/thutasann/go_docker_chat/db"
	"github.com/thutasann/go_docker_chat/internal/user"
	"github.com/thutasann/go_docker_chat/internal/ws"
	"github.com/thutasann/go_docker_chat/router"
)

// Go Docker Postgres Chat Application
func main() {
	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Database Initialize failed: %s", err)
	}
	log.Println(":::: Database Initialized ::::")

	userRep := user.NewRepository(dbConn.GetDB())
	userSvc := user.NewService(userRep)
	userHandler := user.NewHandler(userSvc)

	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)
	go hub.Run()

	router.InitRouter(userHandler, wsHandler)
	router.Start("0.0.0.0:1335")
}
