.PHONY: default dev build test docs clean
# Variables
APP_NAME=cotacao

# Tasks
default: run-with-docs

dev:
	@air 
run-server:
	@swag init -g cmd/main.go
	@go run cmd/main.go
run-client:

build-server:
	@go build -o $(APP_NAME) cmd/server/main.go
test:
	@go test ./... -cover