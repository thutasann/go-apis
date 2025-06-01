package main

import (
	"fmt"
	"log"
	"net/http"
)

const PORT = ":4200"

// Websockets with Go
func main() {
	setupAPI()

	fmt.Printf("🚀 Serving the app in http://localhost%s", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}

func setupAPI() {
	http.Handle("/", http.FileServer(http.Dir("./client")))
}
