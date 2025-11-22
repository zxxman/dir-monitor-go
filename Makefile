.PHONY: build clean test run install fmt vet lint staticcheck quality deps install-service uninstall-service docs help

BINARY_NAME=dir-monitor-go
CONFIG_DIR=configs
LOG_DIR=logs
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT_HASH=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.CommitHash=$(COMMIT_HASH)"

build:
	@echo "Building $(BINARY_NAME) ..."
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/dir-monitor-go
	@echo "Build complete: $(BINARY_NAME)"



clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(LOG_DIR)/*.log $(LOG_DIR)/*.out $(LOG_DIR)/*.err 2>/dev/null || true
	@echo "Clean complete"

test:
	@echo "Running tests..."
	@go test -v ./...

run: build
	@echo "Running ..."
	@./$(BINARY_NAME) -config $(CONFIG_DIR)/config.json

deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed"

docs:
	@echo "Generating documentation..."
	@mkdir -p docs/api
	@echo "Documentation generation complete"

install-service: build
	@echo "Installing service..."
	@sudo bash deploy/deploy-service.sh
	@echo "Service installation complete"

uninstall-service:
	@echo "Uninstalling service..."
	@sudo bash deploy/uninstall-service.sh
	@echo "Service uninstallation complete"

install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "Installation complete"

init-dirs:
	@mkdir -p $(LOG_DIR)

fmt:
	@echo "Formatting Go files..."
	@go fmt ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

staticcheck:
	@echo "Running staticcheck..."
	@staticcheck ./... 2>/dev/null || echo "staticcheck not found; install with 'go install honnef.co/go/tools/cmd/staticcheck@latest'"

lint: fmt vet staticcheck

quality: fmt vet staticcheck

help:
	@echo "Available targets:"
	@echo "  build          - Build the binary"
	@echo "  clean          - Clean build artifacts and logs"
	@echo "  test           - Run tests"
	@echo "  run            - Build and run the application"
	@echo "  deps           - Install dependencies"
	@echo "  install        - Install binary to /usr/local/bin"
	@echo "  install-service- Install as a system service"
	@echo "  uninstall-service- Uninstall the system service"
	@echo "  fmt            - Format Go code"
	@echo "  vet            - Run go vet"
	@echo "  staticcheck    - Run staticcheck (if installed)"
	@echo "  lint           - Run fmt, vet, staticcheck"
	@echo "  quality        - Alias for lint"
	@echo "  docs           - Generate documentation"
	@echo "  help           - Show this help"

default: build
