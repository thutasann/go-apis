package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/thutasann/rssagg/handlers"
	config "github.com/thutasann/rssagg/internal"
	"github.com/thutasann/rssagg/internal/database"
	"github.com/thutasann/rssagg/middlewares"
	"github.com/thutasann/rssagg/utilities"

	_ "github.com/lib/pq"
)

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
	apiCfg := &config.APIConfig{
		DB: database.New(conn),
	}

	// Start Scraping
	go utilities.StartScraping(apiCfg.DB, 10, time.Minute)

	// Initialize Handlers with APIConfig
	h := handlers.Handler{API: apiCfg}

	// Chi Router
	router := chi.NewRouter()

	// CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Use(middlewares.ResponseTimeMiddleware)

	// Auth Middleware
	middlewareHandler := &middlewares.AuthMiddleware{Cfg: h.API}

	// v1 Router
	v1Router := chi.NewRouter()
	v1Router.Get("/health", h.HealthHandler)
	v1Router.Post("/users", h.CreateUserHandler)
	v1Router.Get("/users", middlewareHandler.AuthMiddleware(h.GetUserByAPIKeyHandler))
	v1Router.Post("/feeds", middlewareHandler.AuthMiddleware(h.HandleCreateFeed))
	v1Router.Get("/feeds", middlewareHandler.AuthMiddleware(h.HandlerGetFeeds))
	v1Router.Get("/posts", middlewareHandler.AuthMiddleware(h.GetPostsForUserHandler))
	v1Router.Post("/feed_follows", middlewareHandler.AuthMiddleware(h.HandleCreateFeedFollows))
	v1Router.Get("/feed_follows", middlewareHandler.AuthMiddleware(h.HandleGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", middlewareHandler.AuthMiddleware(h.HandleDeleteFeedFollows))

	// Mount the Router
	router.Mount("/api/v1", v1Router)

	// HTTP Server
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	// Listen and Serve
	log.Printf(":::: Server Starting on PORT %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Server Cannot Start --> ", err)
	}
}
