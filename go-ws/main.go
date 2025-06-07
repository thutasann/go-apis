package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// Application PORT
const PORT = ":4200"

// Websockets with Go
func main() {
	setupAPI()

	fmt.Printf("🚀 Serving the app in https://localhost%s", PORT)
	log.Fatal(http.ListenAndServeTLS(PORT, "server.crt", "server.key", nil))
}

// Private: setup api
func setupAPI() {
	ctx := context.Background()
	manager := NewManager(ctx)
	fs := http.FileServer(http.Dir("./client"))

	http.Handle("/", withNoCache(fs))
	http.HandleFunc("/ws", manager.serverWS)
	http.HandleFunc("/login", manager.loginHandler)
}
