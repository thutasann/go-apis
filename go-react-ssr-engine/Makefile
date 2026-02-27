.PHONY: dev build run clean test deps

# Install Go dependencies
deps:
	go get rogchap.com/v8go@v0.9.0
	go get github.com/evanw/esbuild@v0.21.4
	go get github.com/fsnotify/fsnotify@v1.7.0
	go get github.com/valyala/fasthttp@v1.55.0
	go mod tidy

# Dev mode - race detector on
dev:
	go run -race cmd/server/main.go -dev

# Production binary
build:
	CGO_ENABLED=1 go build -ldflags="-s -w" -o bin/reactgo cmd/server/main.go

# Run production binary
run: build
	./bin/reactgo

# Wipe build artifacts
clean:
	rm -rf bin/ .build/

# Tests with race detector
test:
	go test -race -count=1 ./...