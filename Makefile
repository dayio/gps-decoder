APP_NAME=gps-decoder
MAIN_FILE=cmd/$(APP_NAME)/main.go
BIN_DIR=bin
DATA_DIR=data

.PHONY: all build run test lint clean deps help

all: lint test build

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(MAIN_FILE)

run: build
	./$(BIN_DIR)/$(APP_NAME)

test:
	go test -v -race ./...

lint:
	golangci-lint run

clean:
	rm -rf $(BIN_DIR)
	go clean -testcache

deps:
	go mod tidy
	go mod download