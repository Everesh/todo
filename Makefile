.PHONY: build build-all release clean clean-release clean-all test help

# Application name
APP_NAME := todo

# Build directory
BUILD_DIR := build
RELEASE_DIR := releases

# Version (can be overridden)
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOFLAGS := -ldflags "-X main.version=$(VERSION)"

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build for current platform
	@echo "Building $(APP_NAME) for $(shell go env GOOS)/$(shell go env GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME)$(if $(filter Windows%,$(OS)),.exe,) ./cmd/todo

build-all: ## Build for all platforms
	@echo "Building $(APP_NAME) for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@echo "Building Windows (amd64)..."
	@GOOS=windows GOARCH=amd64 $(GOBUILD) $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe ./cmd/todo
	@echo "Building Windows (386)..."
	@GOOS=windows GOARCH=386 $(GOBUILD) $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-386.exe ./cmd/todo
	@echo "Building Windows (arm64)..."
	@GOOS=windows GOARCH=arm64 $(GOBUILD) $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-windows-arm64.exe ./cmd/todo
	@echo "Building macOS (amd64)..."
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 ./cmd/todo
	@echo "Building macOS (arm64)..."
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 ./cmd/todo
	@echo "Building Linux (amd64)..."
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 ./cmd/todo
	@echo "Building Linux (386)..."
	@GOOS=linux GOARCH=386 $(GOBUILD) $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-386 ./cmd/todo
	@echo "Building Linux (arm64)..."
	@GOOS=linux GOARCH=arm64 $(GOBUILD) $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-arm64 ./cmd/todo
	@echo "Building Linux (arm)..."
	@GOOS=linux GOARCH=arm $(GOBUILD) $(GOFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux-arm ./cmd/todo
	@echo "Done! Binaries are in $(BUILD_DIR)/"

release: build-all ## Prepare release binaries with simple names
	@echo "Preparing release binaries..."
	@rm -rf $(RELEASE_DIR)
	@mkdir -p $(RELEASE_DIR)
	@echo "Preparing Windows binaries..."
	@mkdir -p $(RELEASE_DIR)/windows-amd64 && cp $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe $(RELEASE_DIR)/windows-amd64/$(APP_NAME).exe
	@mkdir -p $(RELEASE_DIR)/windows-386 && cp $(BUILD_DIR)/$(APP_NAME)-windows-386.exe $(RELEASE_DIR)/windows-386/$(APP_NAME).exe
	@mkdir -p $(RELEASE_DIR)/windows-arm64 && cp $(BUILD_DIR)/$(APP_NAME)-windows-arm64.exe $(RELEASE_DIR)/windows-arm64/$(APP_NAME).exe
	@echo "Preparing macOS binaries..."
	@mkdir -p $(RELEASE_DIR)/darwin-amd64 && cp $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 $(RELEASE_DIR)/darwin-amd64/$(APP_NAME) && chmod +x $(RELEASE_DIR)/darwin-amd64/$(APP_NAME)
	@mkdir -p $(RELEASE_DIR)/darwin-arm64 && cp $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(RELEASE_DIR)/darwin-arm64/$(APP_NAME) && chmod +x $(RELEASE_DIR)/darwin-arm64/$(APP_NAME)
	@echo "Preparing Linux binaries..."
	@mkdir -p $(RELEASE_DIR)/linux-amd64 && cp $(BUILD_DIR)/$(APP_NAME)-linux-amd64 $(RELEASE_DIR)/linux-amd64/$(APP_NAME) && chmod +x $(RELEASE_DIR)/linux-amd64/$(APP_NAME)
	@mkdir -p $(RELEASE_DIR)/linux-386 && cp $(BUILD_DIR)/$(APP_NAME)-linux-386 $(RELEASE_DIR)/linux-386/$(APP_NAME) && chmod +x $(RELEASE_DIR)/linux-386/$(APP_NAME)
	@mkdir -p $(RELEASE_DIR)/linux-arm64 && cp $(BUILD_DIR)/$(APP_NAME)-linux-arm64 $(RELEASE_DIR)/linux-arm64/$(APP_NAME) && chmod +x $(RELEASE_DIR)/linux-arm64/$(APP_NAME)
	@mkdir -p $(RELEASE_DIR)/linux-arm && cp $(BUILD_DIR)/$(APP_NAME)-linux-arm $(RELEASE_DIR)/linux-arm/$(APP_NAME) && chmod +x $(RELEASE_DIR)/linux-arm/$(APP_NAME)
	@echo "Done! Release binaries are in $(RELEASE_DIR)/"
	@echo "Each platform folder contains a simple '$(APP_NAME)' or '$(APP_NAME).exe' binary"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "Done!"

clean-release: ## Clean release artifacts
	@echo "Cleaning release artifacts..."
	@rm -rf $(RELEASE_DIR)
	@echo "Done!"

clean-all: clean clean-release ## Clean both build and release artifacts

test: ## Run tests
	$(GOCMD) test -v ./...

