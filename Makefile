.PHONY: build run run/postgres status stop stop/postgres logs test clean help

SERVICE_NAME=api
VERSION=$(shell cat VERSION)
PWD=$(shell pwd)

build:
	@echo "Building the project..."
	docker buildx build -t "${SERVICE_NAME}:local" .

run/postgres:
	docker compose --project-name=devicemanager -f compose.yml up -d postgres

run:
	@echo "Running the project..."
	docker compose --project-name=devicemanager -f compose.yml up -d api

status:
	@echo "Project status:"
	docker compose --project-name=devicemanager ps

stop/postgres:
	docker compose --project-name=devicemanager -f compose.yml stop -d postgres

stop:
	docker compose --project-name=devicemanager stop api

logs:
	docker compose --project-name=devicemanager logs -f ${SERVICE_NAME}

test:
	docker run --rm -v ${PWD}:/app -w /app ${SERVICE_NAME}:local sh -c 'go test -cover -v ./...'

clean:
	docker rmi -f ${SERVICE_NAME}:local 2> /dev/null || true

help:
	@echo "Usage:"
	@echo "  make build     		- Build Docker image of the device manager api"
	@echo "  make run       		- Run api container"
	@echo "  make run/postgres		- Run postgres as dependency to run the api"
	@echo "  make status    		- Show container status"
	@echo "  make logs      		- Tail logs of api container"
	@echo "  make test      		- Run Go tests on docker container"
	@echo "  make stop      		- Stop api container"
	@echo "  make stop/postgres		- Stop postgres as dependency to run the api"
	@echo "  make clean     		- Remove the api containers"