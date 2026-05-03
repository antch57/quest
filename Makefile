BINARY_NAME := quest
BUILD_DIR   := ./bin
CMD_PATH    := .

.PHONY: all build run install lint vet fmt test coverage clean tidy help

## all: build the binary (default)
all: build

## build: compile the binary into ./bin/quest
build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)

## run: build and run the CLI (pass args with ARGS="...")
run: build
	$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

## install: install the binary to $GOPATH/bin
install:
	go install $(CMD_PATH)

## lint: run golangci-lint on all Go files
lint:
	golangci-lint run ./...

## vet: run go vet
vet:
	go vet ./...

## fmt: format all Go source files
fmt:
	gofmt -w .

## test: run all tests with verbose output
test:
	go test -v ./...

## coverage: run tests and open an HTML coverage report
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

## clean: remove build artifacts and coverage output
clean:
	rm -rf $(BUILD_DIR) coverage.out

## tidy: tidy and verify go modules
tidy:
	go mod tidy
	go mod verify

## help: list available targets
help:
	@awk 'BEGIN {FS=": "; cyan="\033[36m"; yellow="\033[33m"; reset="\033[0m"; printf "\n" cyan "Available targets" reset "\n\n"} /^## / {gsub(/^## /, "", $$0); split($$0, a, ": "); if (length(a) > 1) printf "  " yellow "%-14s" reset " %s\n", a[1], a[2];} END {printf "\n"}' Makefile
