package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/thutasann/go-react-ssr-engine/internal/config"
)

func main() {
	dev := flag.Bool("dev", false, "enable dev mode")
	port := flag.Int("port", 0, "overdie port")
	flag.Parse()

	cfg, err := config.Load("reactgo.config.json")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if *dev {
		cfg.Dev = true
	}
	if *port != 0 {
		cfg.Port = *port
	}

	fmt.Printf("reactgo | port=%d dev=%v workers=%d cache=%d\n",
		cfg.Port, cfg.Dev, cfg.WorkerPoolSize, cfg.CacheMaxEntries)

	// Holds the process alive until SIGINT or SIGTERM.
	// Phase 2+ will start the server before this block.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutdown")
}
