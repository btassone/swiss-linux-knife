.PHONY: build run clean test lint fmt

# Binary name
BINARY_NAME=swiss-linux-knife
MAIN_PATH=cmd/swiss-linux-knife

# Build the application
build:
	go build -o bin/$(BINARY_NAME) ./$(MAIN_PATH)

# Run the application
run:
	go run ./$(MAIN_PATH)

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY_NAME)-linux-amd64 ./$(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build -o bin/$(BINARY_NAME)-darwin-amd64 ./$(MAIN_PATH)
	GOOS=windows GOARCH=amd64 go build -o bin/$(BINARY_NAME)-windows-amd64.exe ./$(MAIN_PATH)