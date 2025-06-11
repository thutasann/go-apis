package main

import (
	"log"

	"github.com/thutasann/go_docker_chat/db"
)

// Go Docker Postgres Chat Application
func main() {
	_, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("Database Initialize failed: %s", err)
	}

	log.Println(":::: Database Initialized ::::")
}
