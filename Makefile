.ONESHELL:
.DEFAULT_GOAL := help

PROJECT := notifymem
BINARY := bin/$(PROJECT)-linux-amd64
SOURCES := $(wildcard *.go cmd/*/*.go)
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "Unknown")

build: $(SOURCES) ## Build the project
	@echo "Building $(PROJECT) ($(VERSION))"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.version=$(VERSION)'" -o $(BINARY) cmd/notifymem/*.go

.PHONY: help
help: ## Display help
	@COL_W=10
	@grep -h '##' $(MAKEFILE_LIST) | \
	  grep -v grep | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-'$$COL_W's\033[0m %s\n", $$1, $$2}'

.PHONY: logs
logs: ## Display logs for systemd service (journalctl)
	journalctl -u $(PROJECT) -f

.PHONY: reload
reload: ## Reload systemd service
	sudo systemctl daemon-reload

.PHONY: restart
restart: ## Restart systemd service
	sudo systemctl restart $(PROJECT)

.PHONY: run
run: build ## Run the project
	./$(BINARY)

.PHONY: start
start: ## Start systemd service
	sudo systemctl start $(PROJECT)

.PHONY: status
status: ## Display status of systemd service
	sudo systemctl status $(PROJECT)

.PHONY: stop
stop: ## Stop systemd service
	sudo systemctl stop $(PROJECT)

.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: tidy
tidy: ## Run go mod tidy
	go mod tidy -v
