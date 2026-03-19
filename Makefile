.PHONY: build test install clean release run help

BINARY_NAME=gscex
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION)"

help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

build: ## Build the binary for current platform
	go build $(LDFLAGS) -o $(BINARY_NAME) ./cmd/gscex/
	@echo "Built: $(BINARY_NAME)"

test: ## Run all tests
	go test -v ./...

install: ## Install to $GOPATH/bin
	go install $(LDFLAGS) ./cmd/gscex/
	@echo "Installed to $$(go env GOPATH)/bin/$(BINARY_NAME)"

clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)
	rm -rf dist/
	@echo "Cleaned"

run: build ## Build and run the TUI
	./$(BINARY_NAME) tui

init: build ## Build and run init
	./$(BINARY_NAME) init

release: ## Build release for all platforms (requires VERSION arg)
	@if [ -z "$(VERSION)" ] || [ "$(VERSION)" = "dev" ]; then \
		echo "Usage: make release VERSION=v1.0.0"; \
		exit 1; \
	fi
	./scripts/build-release.sh $(VERSION)

quick-release: build ## Quick local release build (no GitHub)
	./scripts/build-release.sh dev

fmt: ## Format Go code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

deps: ## Download dependencies
	go mod download
	go mod tidy

.DEFAULT_GOAL := help
