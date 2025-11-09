.PHONY: build install clean test

# Build variables
BINARY_NAME=docker-cri
VERSION?=0.1.0
BUILD_DIR=bin
GOOS?=linux
GOARCH?=amd64

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(BINARY_NAME) -ldflags "-X main.version=$(VERSION)" ./main.go
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

install: build
	@echo "Installing $(BINARY_NAME)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/$(BINARY_NAME)
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "Installation complete"

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "Clean complete"

test:
	@echo "Running tests..."
	@go test -v ./...

deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

lint:
	@echo "Running linter..."
	@golangci-lint run || echo "golangci-lint not installed, skipping..."

.DEFAULT_GOAL := build

