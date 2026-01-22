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
	@echo "  make clean      # Clean build artifacts"
	@echo "  make poc        # Run POC experiments"
	@echo "  make docs       # Generate / open documentation"

# -----------------------------------------
# Run / Build
# -----------------------------------------
.PHONY: run
run:
	@echo "Running project..."
	@echo "Define your own run commands per project"

.PHONY: build
build:
	@echo "Building project..."
	@echo "Define your own build commands per project"

# -----------------------------------------
# Tests
# -----------------------------------------
.PHONY: test
test:
	@echo "Running tests..."
	@echo "Define your own test commands per project"

# -----------------------------------------
# POC / Experiments
# -----------------------------------------
.PHONY: poc
poc:
	@echo "Running POC / experiments..."
	@echo "Customize this target per project"

# -----------------------------------------
# Documentation
# -----------------------------------------
.PHONY: docs
docs:
	@echo "Opening / generating documentation..."
	@echo "Customize per project"

# -----------------------------------------
# Clean
# -----------------------------------------
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@echo "Define your own clean commands per project"
