# Build variables
BINARY_NAME=agbcloud
VERSION?=dev
GIT_COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS=-ldflags "-X github.com/agbcloud/agbcloud-cli/cmd.Version=$(VERSION) -X github.com/agbcloud/agbcloud-cli/cmd.GitCommit=$(GIT_COMMIT) -X github.com/agbcloud/agbcloud-cli/cmd.BuildDate=$(BUILD_DATE)"

# Default target
.PHONY: all
all: build

# Build for current platform
.PHONY: build
build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) .

# Build for all platforms
.PHONY: build-all
build-all: build-linux build-darwin build-windows

# Build for Linux
.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-arm64 .

# Build for macOS
.PHONY: build-darwin
build-darwin:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 .

# Build for Windows
.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-arm64.exe .

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf bin/

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run linter
.PHONY: lint
lint:
	golangci-lint run

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Install dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Install the binary
.PHONY: install
install:
	go install $(LDFLAGS) .

# Development build and run
.PHONY: dev
dev: build
	./bin/$(BINARY_NAME)

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - Build for current platform"
	@echo "  build-all    - Build for all platforms"
	@echo "  build-linux  - Build for Linux (amd64, arm64)"
	@echo "  build-darwin - Build for macOS (amd64, arm64)"
	@echo "  build-windows- Build for Windows (amd64, arm64)"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  lint         - Run linter"
	@echo "  fmt          - Format code"
	@echo "  deps         - Install dependencies"
	@echo "  install      - Install binary to GOPATH/bin"
	@echo "  dev          - Build and run for development"
	@echo "  help         - Show this help" 