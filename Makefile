.PHONY: build run test clean

# Binary name
BINARY_NAME=yuyan

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod
GOGET=$(GOCMD) get

# Main entry point
MAIN_PATH=cmd/server/main.go

# Default target
all: build

# Build binary
build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

# Run directly
run:
	$(GORUN) $(MAIN_PATH)

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe

# Update dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Build for multiple platforms
build-all: clean
	# Linux
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	
	# Windows
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	
	# macOS
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)

# Initialize development environment
init:
	mkdir -p data
	cp -n config/config.sample.yaml config/config.yaml || true

# Help
help:
	@echo "Available targets:"
	@echo "  build       - Build binary"
	@echo "  run         - Run the application"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  deps        - Update dependencies"
	@echo "  build-all   - Build for multiple platforms"
	@echo "  init        - Initialize development environment"
	@echo "  help        - Show this help message" 