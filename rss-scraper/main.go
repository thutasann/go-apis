package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/thutasann/rssagg/handlers"
	"github.com/thutasann/rssagg/internal/database"

	_ "github.com/lib/pq"
)

// API Config
type apiConfig struct {
	// Database
	DB *database.Queries
}

// RSS Scraper
func main() {
	fmt.Println(":::: RSS Scraper ::::")
	godotenv.Load(".env")

	// PORT
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT Is not found in the ENV")
	}

	// DB_URL
	dbString := os.Getenv("DB_URL")
	if dbString == "" {
		log.Fatal("DB_URL is missing in the ENV")
	}

	// DB Connect
	conn, err := sql.Open("postgres", dbString)
	if err != nil {
		log.Fatal("DB Connection Failed..")
	}

	// API Config
	apiConfig := apiConfig{
		DB: database.New(conn),
	}

	// Router
	router := chi.NewRouter()

	// Cors
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// v1 router
	v1Router := chi.NewRouter()
	v1Router.Get("/health", handlers.HandlerReadiness)
	v1Router.Get("/err", handlers.HandlerErr)

	// Mount the Router
	router.Mount("/api/v1", v1Router)

	// HTTP Server
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	// Listen and Serve
	log.Printf(":::: Server Starting on PORT %v", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal("Server Cannot Start --> ", err)
	}
}
