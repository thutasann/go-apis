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

	ctx := context.Background()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal(err)
	}

	db := mongoClient.Database(cfg.MongoDBName)
	repo := repository.NewMongoEventRepository(db, "events")

	// Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	queue := queue.NewRedisQueue(redisClient, "webhook:events")

	// Worker Pool
	pool := worker.NewPool(10, queue, repo)

	ctx, cancel := context.WithCancel(context.Background())
	pool.Start(ctx)

	log.Println("worker pool started")

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

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("shutting down...")

	cancel()
	pool.Shutdown()

	shutdownHTTP, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_ = server.Shutdown(shutdownHTTP)

	shutdownCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_ = mongoClient.Disconnect(shutdownCtx)

	log.Println("shutdown complete")
}
