.PHONY: default dev build test docs clean
# Variables
APP_NAME=cotacao

# Tasks
default: run-with-docs

dev:
	@air 
run-with-docs:
	@swag init -g cmd/main.go
	@go run cmd/main.go
build:
	@go build -o $(APP_NAME) cmd/main.go
test:
	@go test ./ ...
docs:
	sh ./scripts/generate_swagger_docs.sh