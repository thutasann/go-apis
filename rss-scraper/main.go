package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

// RSS Scraper
func main() {
	fmt.Println(":::: RSS Scraper ::::")

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT Is not found in the ENV")
	}

	router := chi.NewRouter()

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf(":::: Server Starting on PORT %v", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("Server Cannot Start --> ", err)
	}

}
