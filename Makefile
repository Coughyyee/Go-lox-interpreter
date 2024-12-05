.PHONY: build run run_file test clean

# Binary name and path
BINARY_NAME=lox
BINARY_PATH=bin/$(BINARY_NAME)

# Get the first argument passed to make, default to main.lox if none provided
args = $(filter-out $@,$(MAKECMDGOALS))
file = $(if $(call args),$(call args),main.lox)

# Ensure the bin directory exists
$(shell mkdir -p bin)

build:
	@go build -o $(BINARY_PATH) .

run: build
	@./$(BINARY_PATH)

run_file: build
	@./$(BINARY_PATH) $(file)

test:
	@go test ./... -v

clean:
	@rm -rf bin/

# Rule to handle any arguments
%:
	@:
