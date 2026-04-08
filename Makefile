.PHONY: build run clean test lint install

BINARY_NAME=gitplus
BUILD_DIR=dist

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/gitplus

run:
	go run ./cmd/gitplus

clean:
	rm -rf $(BUILD_DIR)

test:
	go test -v ./...

lint:
	golangci-lint run

install:
	go install ./cmd/gitplus