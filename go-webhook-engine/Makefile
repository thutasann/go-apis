APP_NAME=go-webhook-engine
CMD_PATH=./cmd/server
BINARY_NAME=server

.PHONY: run build test clean tidy fmt lint

run:
	go run $(CMD_PATH)/main.go

build:
	go build -o bin/$(BINARY_NAME) $(CMD_PATH)/main.go

test:
	go test ./...

clean:
	rm -rf bin/

tidy:
	go mod tidy

fmt:
	go fmt ./...

lint:
	golangci-lint run

deps:
	go get github.com/redis/go-redis/v9
	go get go.mongodb.org/mongo-driver/mongo
	go get go.mongodb.org/mongo-driver/bson
