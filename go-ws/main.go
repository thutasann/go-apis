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

// Private: setup api
func setupAPI() {
	manager := NewManager()

	fs := http.FileServer(http.Dir("./client"))
	http.Handle("/", withNoCache(fs))

	http.HandleFunc("/ws", manager.serverWS)
}

// Middleware to add no-cache headers
func withNoCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setNoCacheHeaders(w)
		h.ServeHTTP(w, r)
	})
}

// Utility function to set no-cache headers
func setNoCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}
