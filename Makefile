.PHONY: build run status stop logs test test/local clean help

SERVICE_NAME=api
VERSION=$(shell cat VERSION)
PWD=$(shell pwd)

build:
	@echo "Building the project..."
	docker buildx build -t "${SERVICE_NAME}:local" .

run:
	@echo "Running the project..."
	docker compose --project-name=devicemanager -f compose.yml up -d

status:
	@echo "Project status:"
	docker compose --project-name=devicemanager ps

stop:
	docker compose --project-name=devicemanager stop

logs:
	docker compose --project-name=devicemanager logs -f ${SERVICE_NAME}

test:
	docker run --rm -v ${PWD}:/app -w /app ${SERVICE_NAME}:local sh -c 'go test -v ./...'

test/local:
	go test -covermode=atomic ./internal/... -coverprofile=cover.out  && go tool cover -func=cover.out

clean:
	docker rmi -f ${SERVICE_NAME}:local 2> /dev/null || true

help:
	@echo "Usage:"
	@echo "  make build     - Build Docker image"
	@echo "  make run       - Run containers"
	@echo "  make status    - Show container status"
	@echo "  make logs      - Tail logs"
	@echo "  make test      - Run Go tests"
	@echo "  make stop      - Stop containers"
	@echo "  make clean     - Remove containers"