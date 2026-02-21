.PHONY: dev build run clean test

# Dev mode - hot reload flag, race detector on
dev:
	go run -race cmd/server/main.go -dev

# Production binary - static linked, optimized
build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/reactgo cmd/server/main.go

# Run production binary
run: build
	./bin/reactgo

# Wipe build artifacts
clean:
	rm -rf bin/ .build/

# Tests with race detector
test:
	go test -race -count=1 ./...

# Init module (run once)
init:
	go mod init github.com/yourname/reactgo
	go mod tidy