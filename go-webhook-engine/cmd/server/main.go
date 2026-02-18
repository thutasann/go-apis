package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thutasann/go-webhook-engine/internal/config"
	http2 "github.com/thutasann/go-webhook-engine/internal/delivery/http"
	"github.com/thutasann/go-webhook-engine/internal/queue"
	"github.com/thutasann/go-webhook-engine/internal/repository"
	"github.com/thutasann/go-webhook-engine/internal/worker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.Load()

	// Root context for entire app lifecycle
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel()

	mongoClient, err := mongo.Connect(rootCtx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal(err)
	}

	db := mongoClient.Database(cfg.MongoDBName)
	repo := repository.NewMongoEventRepository(db, "events")

	if err := repo.EnsureIndexes(rootCtx); err != nil {
		log.Fatal(err)
	}

	// ----- Redis -----
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	queue := queue.NewRedisQueue(redisClient, "webhook:events")

	// ----- Worker Pool -----
	pool := worker.NewPool(10, queue, repo)
	pool.Start(rootCtx)

	log.Println("worker pool started")

	// ----- HTTP -----
	handler := http2.NewHandler(repo, queue)

	mux := http.NewServeMux()
	mux.HandleFunc("/webhook", handler.Webhook)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		log.Println("HTTP server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// ----- Graceful Shutdown -----
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("shutting down...")

	// Cancel root context (stops workers)
	rootCancel()

	// Shutdown HTTP server
	httpCtx, httpCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer httpCancel()

	if err := server.Shutdown(httpCtx); err != nil {
		log.Printf("HTTP shutdown error: %v\n", err)
	}

	// Disconnect Mongo
	mongoCtx, mongoCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer mongoCancel()

	if err := mongoClient.Disconnect(mongoCtx); err != nil {
		log.Printf("Mongo disconnect error: %v\n", err)
	}

	log.Println("shutdown complete")
}
