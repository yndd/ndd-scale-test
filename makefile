# ====================================================================================
# Setup Project

PROJECT_NAME := ndd-scale-test
PROJECT_REPO := github.com/yndd/$(PROJECT_NAME)

BIN_DIR = $(shell pwd)/bin
BINARY = $(shell pwd)/bin/nddscaletest

all: build

build: ## Build binaries: ndd-gen
	mkdir -p $(BIN_DIR)
	go build -o $(BINARY) ./cmd/main.go 

fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet against code.
	go vet ./...

test:
	go test -race ./... -v

lint:
	golangci-lint run

clint:
	docker run -it --rm -v $$(pwd):/app -w /app golangci/golangci-lint:v1.40.1 golangci-lint run -v