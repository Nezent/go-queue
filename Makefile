APP_NAME := go-queue
SRC_FILE := cmd/main.go

build:
	@go build -o bin/$(APP_NAME) $(SRC_FILE)

run: build
	@./bin/$(APP_NAME)

start:
	@docker start $(shell docker ps -aq)

stop:
	@docker stop $(shell docker ps -q)