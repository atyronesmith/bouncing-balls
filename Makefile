# Bouncing Balls - Makefile
# Go project build automation

# Project settings
PROJECT_NAME := bouncing-balls
MODULE_NAME := github.com/asmith/bouncing-balls
CMD_DIR := cmd/bouncing-balls
PKG_DIR := pkg
BUILD_DIR := build
BINARY_NAME := bouncing-balls

# Go settings
GO := go
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# Build flags
LDFLAGS := -ldflags "-s -w"
BUILD_FLAGS := $(LDFLAGS)

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
WHITE := \033[0;37m
NC := \033[0m # No Color

.PHONY: help build run clean test fmt vet lint deps tidy check install uninstall package all

# Default target
all: clean fmt vet build

# Help target
help: ## Show this help message
	@echo "$(CYAN)Bouncing Balls - Available Commands:$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""

# Build the application
build: ## Build the application binary
	@echo "$(BLUE)Building $(PROJECT_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "$(GREEN)✓ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

# Run the application directly
run: ## Run the application without building
	@echo "$(BLUE)Running $(PROJECT_NAME)...$(NC)"
	$(GO) run ./$(CMD_DIR)

# Run the built binary
run-binary: build ## Build and run the binary
	@echo "$(BLUE)Running built binary...$(NC)"
	./$(BUILD_DIR)/$(BINARY_NAME)

# Clean build artifacts
clean: ## Clean build artifacts and cache
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	$(GO) clean -cache -testcache -modcache
	@echo "$(GREEN)✓ Clean complete$(NC)"

# Run tests
test: ## Run all tests
	@echo "$(BLUE)Running tests...$(NC)"
	$(GO) test -v ./...
	@echo "$(GREEN)✓ Tests complete$(NC)"

# Run tests with coverage
test-coverage: ## Run tests with coverage report
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(NC)"

# Format code
fmt: ## Format Go code
	@echo "$(BLUE)Formatting code...$(NC)"
	$(GO) fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

# Run go vet
vet: ## Run go vet
	@echo "$(BLUE)Running go vet...$(NC)"
	$(GO) vet ./...
	@echo "$(GREEN)✓ Vet complete$(NC)"

# Run golint (if available)
lint: ## Run golint (requires golint to be installed)
	@echo "$(BLUE)Running golint...$(NC)"
	@if command -v golint >/dev/null 2>&1; then \
		golint ./...; \
		echo "$(GREEN)✓ Lint complete$(NC)"; \
	else \
		echo "$(YELLOW)⚠ golint not installed. Install with: go install golang.org/x/lint/golint@latest$(NC)"; \
	fi

# Download dependencies
deps: ## Download and install dependencies
	@echo "$(BLUE)Downloading dependencies...$(NC)"
	$(GO) mod download
	@echo "$(GREEN)✓ Dependencies downloaded$(NC)"

# Tidy modules
tidy: ## Tidy and verify modules
	@echo "$(BLUE)Tidying modules...$(NC)"
	$(GO) mod tidy
	$(GO) mod verify
	@echo "$(GREEN)✓ Modules tidied$(NC)"

# Run all checks
check: fmt vet test ## Run all code quality checks
	@echo "$(GREEN)✓ All checks passed$(NC)"

# Install the binary to GOPATH/bin
install: build ## Install the binary to GOPATH/bin
	@echo "$(BLUE)Installing $(PROJECT_NAME)...$(NC)"
	$(GO) install ./$(CMD_DIR)
	@echo "$(GREEN)✓ Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)$(NC)"

# Uninstall the binary
uninstall: ## Remove the installed binary
	@echo "$(YELLOW)Uninstalling $(PROJECT_NAME)...$(NC)"
	@rm -f $(shell go env GOPATH)/bin/$(BINARY_NAME)
	@echo "$(GREEN)✓ Uninstalled$(NC)"

# Create a distributable package
package: build ## Create a distributable package
	@echo "$(BLUE)Creating package...$(NC)"
	@mkdir -p $(BUILD_DIR)/package
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(BUILD_DIR)/package/
	@cp README.md $(BUILD_DIR)/package/ 2>/dev/null || echo "No README.md found"
	@tar -czf $(BUILD_DIR)/$(PROJECT_NAME)-$(GOOS)-$(GOARCH).tar.gz -C $(BUILD_DIR)/package .
	@echo "$(GREEN)✓ Package created: $(BUILD_DIR)/$(PROJECT_NAME)-$(GOOS)-$(GOARCH).tar.gz$(NC)"

# Cross-compile for different platforms
build-all: ## Cross-compile for multiple platforms
	@echo "$(BLUE)Cross-compiling for multiple platforms...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@for os in linux darwin windows; do \
		for arch in amd64 arm64; do \
			if [ "$$os" = "windows" ]; then \
				ext=".exe"; \
			else \
				ext=""; \
			fi; \
			echo "Building for $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch $(GO) build $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$$os-$$arch$$ext ./$(CMD_DIR); \
		done; \
	done
	@echo "$(GREEN)✓ Cross-compilation complete$(NC)"

# Development setup
dev-setup: deps ## Set up development environment
	@echo "$(BLUE)Setting up development environment...$(NC)"
	@if ! command -v golint >/dev/null 2>&1; then \
		echo "Installing golint..."; \
		$(GO) install golang.org/x/lint/golint@latest; \
	fi
	@if ! command -v goimports >/dev/null 2>&1; then \
		echo "Installing goimports..."; \
		$(GO) install golang.org/x/tools/cmd/goimports@latest; \
	fi
	@echo "$(GREEN)✓ Development environment ready$(NC)"

# Show project info
info: ## Show project information
	@echo "$(CYAN)Project Information:$(NC)"
	@echo "  Name: $(PROJECT_NAME)"
	@echo "  Module: $(MODULE_NAME)"
	@echo "  Go Version: $(shell go version)"
	@echo "  GOOS: $(GOOS)"
	@echo "  GOARCH: $(GOARCH)"
	@echo "  Build Dir: $(BUILD_DIR)"
	@echo "  Binary: $(BINARY_NAME)"
	@echo ""
	@echo "$(CYAN)Project Structure:$(NC)"
	@tree -I 'build|*.sum|go.mod' . 2>/dev/null || find . -type f -name "*.go" | head -20

# Quick development cycle
dev: clean fmt vet build run-binary ## Quick development cycle: clean, format, vet, build, and run

# Production build
prod: clean test build package ## Production build: clean, test, build, and package