APP_NAME = main
BUILD_DIR = build

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) *.go

dev:
	@echo "Running in development mode..."
	@go run *.go

run: build
	./$(BUILD_DIR)/$(APP_NAME)

clean:
	rm -rf $(BUILD_DIR)
