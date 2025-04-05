.PHONY: dev build run clean test

APP_NAME=repl-cli
BUILD_DIR=build

dev:
	@echo "Running in development mode..."
	@go run main.go

build:
	@echo "Building production binary..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) main.go

run:
	@echo "Running production binary..."
	@./$(BUILD_DIR)/$(APP_NAME)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

test:
	@echo "Running tests..."
	@go test ./...
