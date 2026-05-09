APP_NAME = main
BUILD_DIR = build
MAIN_FILE = cmd/main.go

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

dev:
	@echo "Running in development mode..."
	@go run $(MAIN_FILE)

run: build
	./$(BUILD_DIR)/$(APP_NAME)

clean:
	rm -rf $(BUILD_DIR)
