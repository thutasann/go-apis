BUILD_DIR = bin
APP_NAME = app
SRC_DIR = cmd


postgresinit:
	docker run --name postgres15 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:15-alpine

postgres:
	docker exec -it postgres15 psql

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root go-chat

dropdb:
	docker exec -it postgres15 dropdb go-chat

migrateup:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5433/go-chat?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:password@localhost:5433/go-chat?sslmode=disable" -verbose down

dev:
	go run cmd/main.go

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)/main.go

run: build
	./$(BUILD_DIR)/$(APP_NAME)

.PHONY: postgresinit postgres createdb dropdb migrateup migratedown dev build run