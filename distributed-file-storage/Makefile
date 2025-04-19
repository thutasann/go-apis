dev:
	@go run main.go

build:
	clear
	@go build -o bin/fs

run: build
		@./bin/fs

test: 
	@go test ./...

clean:
	@rm -rf bin