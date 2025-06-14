BUILD_DIR = bin
APP_NAME = cli_todo
SRC_DIR = cmd/todo

dev:
	go run cmd/todo/main.go

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)/main.go

run: build
	./$(BUILD_DIR)/$(APP_NAME)

.PHONY: postgresinit postgres createdb dropdb migrateup migratedown dev build run