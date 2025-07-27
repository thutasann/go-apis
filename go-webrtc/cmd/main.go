package main

import (
	"log"

	"github.com/thuta/gowebrtc/internal/server"
)

// GoLang videochat app
func main() {
	if err := server.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}
