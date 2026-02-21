package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/thutasann/go-ssr-engine/internal/engine"
	apphttp "github.com/thutasann/go-ssr-engine/internal/http"
	"github.com/thutasann/go-ssr-engine/internal/worker"
)

func main() {
	// Use all CPUS
	runtime.GOMAXPROCS(runtime.NumCPU())

	tpl, err := engine.Compile([]byte("Hello {{first}} {{last}}"))
	if err != nil {
		log.Fatal(err)
	}

	// Build variable name -> index map from compiler
	varNameToIndex := map[string]uint16{
		"first": 0,
		"last":  1,
	}

	// Worker pool sizing strategy:
	// workers = CPU cores * 2 (good starting point for CPU-bound tasks)
	workers := runtime.NumCPU() * 2
	queueSize := 10000 // bounded queue to absorb burst

	pool := worker.New(workers, queueSize)
	pool.Start()

	handler := &apphttp.Handler{
		Pool:           pool,
		Tpl:            tpl,
		VarNameToIndex: varNameToIndex,
	}

	server := apphttp.NewServer(":8080", handler)

	go func() {
		log.Println("Server running on :8080")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Println("Shutting down...")
	pool.Shutdown()
}
