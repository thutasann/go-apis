package main

import (
	"log"

	"github.com/thuta/ecomm/db"
	"github.com/thuta/ecomm/ecomm-api/handler"
	"github.com/thuta/ecomm/ecomm-api/server"
	"github.com/thuta/ecomm/ecomm-api/storer"
)

// Ecomm API
func main() {
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	defer db.Close()
	log.Println("successfully connected to database")

	st := storer.NewMySQLStorer(db.GetDB())
	srv := server.NewServer(st)
	hdl := handler.NewHandler(srv)
	handler.RegisterRoutes(hdl)
	handler.Start(":8080")
}
