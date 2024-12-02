.PHONY: build run run_file test clean

# Binary name and path
BINARY_NAME=lox
BINARY_PATH=bin/$(BINARY_NAME)

# Ensure the bin directory exists
$(shell mkdir -p bin)

build:
	@go build -o $(BINARY_PATH) .

run: build
	@./$(BINARY_PATH)

run_file: build
	@./$(BINARY_PATH) main.lox

test:
	@go test ./... -v

clean:
	@rm -rf bin/
	@go clean
