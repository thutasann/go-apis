BUILD_DIR = bin
APP_NAME = cli_todo
SRC_DIR = cmd/todo

dev:
	go run $(SRC_DIR)/main.go

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)/main.go

run: build
	./$(BUILD_DIR)/$(APP_NAME)


# ------- Commands -------

add:
	go run ./$(SRC_DIR) -add

complete:
	go run ./$(SRC_DIR) -complete=$(ID)

.PHONY: postgresinit postgres createdb dropdb migrateup migratedown dev build run