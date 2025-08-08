.PHONY: help format lint run build test clean

# Default target
help:
	@echo "Available commands:"
	@echo "  format  - Format Go code using gofmt and goimports"
	@echo "  lint    - Run golangci-lint"
	@echo "  run     - Run the application with air (hot reload)"
	@echo "  build   - Build the application"
	@echo "  test    - Run tests"
	@echo "  clean   - Clean build artifacts"

# Format Go code
format:
	@echo "Formatting Go code..."
	gofmt -w .
	goimports -w .

# Lint code with golangci-lint
lint:
	@echo "Running golangci-lint..."
	golangci-lint run

# Run with air for hot reload
run:
	@echo "Starting application with air..."
	air

# Build the application
build:
	@echo "Building application..."
	go build -o ./bin/partiuFit .

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf ./bin
	rm -rf ./tmp