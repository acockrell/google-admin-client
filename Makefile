.PHONY: all build clean test lint fmt vet install help docker-build

# Binary name
BINARY_NAME=gac
DOCKER_IMAGE=gac:latest

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build parameters
BUILD_DIR=build
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

all: test build

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) -v

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f $(BINARY_NAME)

## test: Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## lint: Run linter (requires golangci-lint)
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run ./...

## gosec: Run security scanner (requires gosec)
gosec:
	@echo "Running security scanner..."
	@which gosec > /dev/null || (echo "gosec not installed. Install with: go install github.com/securego/gosec/v2/cmd/gosec@latest" && exit 1)
	gosec ./...

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

## tidy: Tidy go modules
tidy:
	@echo "Tidying go modules..."
	$(GOMOD) tidy

## deps: Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download

## update-deps: Update dependencies
update-deps:
	@echo "Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

## install: Install binary to $GOPATH/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

## run: Build and run the binary
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

## check: Run all checks (fmt, vet, lint, gosec, test)
check: fmt vet lint gosec test

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
