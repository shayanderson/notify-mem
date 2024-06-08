.DEFAULT_GOAL := help

PROJECT := notifymem
SOURCES := $(wildcard *.go cmd/*/*.go)
VERSION := $(shell git describe --tags 2>/dev/null || echo "Unknown")

build: $(SOURCES) ## Build the project
	@echo "Building $(PROJECT) ($(VERSION))"
	CGO_ENABLED=0 go build -ldflags "-X 'main.version=$(VERSION)'" -o bin/$(PROJECT) cmd/notifymem/*.go

.PHONY: help
help: ## Display help
	@COL_W=20
	@grep -h '##' $(MAKEFILE_LIST) | \
	  grep -v grep | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-'$$COL_W's\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: build ## Run the project
	./bin/$(PROJECT)

.PHONY: tidy
tidy: ## Run go mod tidy
	go mod tidy -v
