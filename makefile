# -----------------------------------------
# 26-in-26 Project Makefile Template
# -----------------------------------------

# Default target
.PHONY: all
all: help

# -----------------------------------------
# Help
# -----------------------------------------
.PHONY: help
help:
	@echo "26-in-26 Project Makefile"
	@echo
	@echo "Available commands:"
	@echo "  make run        # Run the project"
	@echo "  make build      # Build / compile the project"
	@echo "  make test       # Run tests"

# -----------------------------------------
# Run / Build
# -----------------------------------------
.PHONY: run
run: 
	@echo "Running project..."
	@go run ./cmd/ascii-gen/main.go


.PHONY: build
build:
	@echo "Building project..."
	@go build ./cmd/...

# -----------------------------------------
# Tests
# -----------------------------------------
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...
