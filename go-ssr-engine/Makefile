# ===== CONFIG =====
APP_NAME := go-ssr-engine
BINARY := bin/server
PKG := ./...
PORT := 8080

# Detect CPU count (mac/linux)
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	CPU_COUNT := $(shell sysctl -n hw.ncpu)
else
	CPU_COUNT := $(shell nproc)
endif

# ===== DEFAULT =====
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make build        - Build binary"
	@echo "  make run          - Run server"
	@echo "  make dev          - Run with race detector"
	@echo "  make bench        - Run benchmarks"
	@echo "  make load         - Run wrk load test"
	@echo "  make ab           - Run apache benchmark"
	@echo "  make pprof-cpu    - Capture CPU profile"
	@echo "  make pprof-mem    - Capture memory profile"
	@echo "  make clean        - Clean build artifacts"

# ===== BUILD =====
.PHONY: build
build:
	@mkdir -p bin
	go build -o $(BINARY) ./cmd/server

# ===== RUN =====
.PHONY: run
run: build
	./$(BINARY)

# ===== DEV (RACE DETECTOR) =====
.PHONY: dev
dev:
	go run -race ./cmd/server

# ===== BENCH (engine benchmarks only) =====
.PHONY: bench
bench:
	go test -bench=. -benchmem ./benchmarks

# ===== LOAD TEST (wrk required) =====
.PHONY: load
load:
	wrk -t$(CPU_COUNT) -c1000 -d30s http://localhost:$(PORT)

# ===== Apache Benchmark fallback =====
.PHONY: ab
ab:
	ab -n 200000 -c 500 http://127.0.0.1:$(PORT)/

# ===== CPU PROFILE =====
.PHONY: pprof-cpu
pprof-cpu:
	go test -bench=. -cpuprofile cpu.out ./benchmarks
	go tool pprof cpu.out

# ===== MEMORY PROFILE =====
.PHONY: pprof-mem
pprof-mem:
	go test -bench=. -memprofile mem.out ./benchmarks
	go tool pprof mem.out

# ===== CLEAN =====
.PHONY: clean
clean:
	rm -rf bin
	rm -f cpu.out mem.out