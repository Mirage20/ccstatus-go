# ccstatus-go Makefile
BINARY_NAME := ccstatus
BUILD_DIR := ./build
BIN_DIR := ./bin
CMD_PATH := ./cmd/ccstatus

# Go commands
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test

# golangci-lint
GOLANGCI_VERSION := v2.4.0
GOLANGCI_BIN := $(BIN_DIR)/golangci-lint

# Version information
VERSION ?= dev
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_COMMIT_DATE := $(shell git show -s --format=%ci HEAD 2>/dev/null | cut -d' ' -f1 || echo "")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags with version injection
LDFLAGS := -s -w \
	-X 'main.version=$(VERSION)' \
	-X 'main.commit=$(GIT_COMMIT)' \
	-X 'main.commitDate=$(GIT_COMMIT_DATE)' \
	-X 'main.buildDate=$(BUILD_DATE)'

# Define release platforms
PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/arm64

.PHONY: all
all: lint build ## Run lint and build

.PHONY: build
build: ## Build the binary
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)

.PHONY: clean
clean: ## Clean build artifacts and binaries
	$(GOCLEAN)
	rm -rf $(BUILD_DIR) $(BIN_DIR)

.PHONY: test
test: ## Run tests
	$(GOTEST) -v ./...

$(GOLANGCI_BIN): ## Install golangci-lint to project bin
	mkdir -p $(BIN_DIR)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(BIN_DIR) $(GOLANGCI_VERSION)

.PHONY: lint
lint: $(GOLANGCI_BIN) ## Lint the codebase without fixing
	$(GOLANGCI_BIN) run ./...

.PHONY: lint-fix
lint-fix: $(GOLANGCI_BIN) ## Fix all auto-fixable issues
	$(GOLANGCI_BIN) run --fix ./...

.PHONY: run
run: build ## Build and run with sample input
	cat ccstatus.json | $(BUILD_DIR)/$(BINARY_NAME)

.PHONY: release
release: ## Build for all platforms
	@echo "Building $(BINARY_NAME) $(VERSION) for all platforms..."
	@mkdir -p $(BUILD_DIR)/release
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d'/' -f1); \
		GOARCH=$$(echo $$platform | cut -d'/' -f2); \
		output_name=$(BUILD_DIR)/release/$(BINARY_NAME)-$$GOOS-$$GOARCH; \
		if [ "$$GOOS" = "windows" ]; then \
			output_name="$$output_name.exe"; \
		fi; \
		echo "Building $$GOOS/$$GOARCH..."; \
		GOOS=$$GOOS GOARCH=$$GOARCH $(GOBUILD) \
			-ldflags="$(LDFLAGS)" \
			-o $$output_name $(CMD_PATH) || exit 1; \
	done
	@echo "Creating checksums..."
	@cd $(BUILD_DIR)/release && shasum -a 256 * > checksums.txt
	@echo "Release builds complete! Artifacts in $(BUILD_DIR)/release/"

.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
