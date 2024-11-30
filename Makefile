build:
	@go build -o bin/lox main.go scanner.go token.go token_type.go error.go expr.go parser.go interpreter.go stmt.go environment.go lox.go

run: build
	@./bin/lox

run_file: build
	@./bin/lox main.lox

test:
	@go test ./... -v
