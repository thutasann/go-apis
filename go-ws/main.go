package main

import (
	"fmt"
	"log"
	"net/http"
)

// Application PORT
const PORT = ":4200"

// Websockets with Go
func main() {
	setupAPI()

	fmt.Printf("🚀 Serving the app in http://localhost%s", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}

// Private: setup api
func setupAPI() {
	manager := NewManager()

	fs := http.FileServer(http.Dir("./client"))
	http.Handle("/", withNoCache(fs))

	http.HandleFunc("/ws", manager.serverWS)
}
