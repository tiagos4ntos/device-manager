# Device Manager API

This project provides an API for managing an inventory of devices and its state (available, in-use or inactive).

## Features

- Register, update, and delete devices
- Query all devices and filter by ID, Brand or State

## Requirements
- [Golang](https://go.dev/dl/) v1.25.0
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Make](https://www.gnu.org/software/make/)

## Setup

Clone the repository:

```sh
git clone https://github.com/tiagos4ntos/device-manager
cd device-manager
```

## Dependencies

This application requires a Postgres Database Server and a valid `.env` file for configuration. 

Both are provided and managed automatically by the included Docker Compose setup and the `.env.template` file.

### Environment Setup

1. Copy `.env.template` to `.env` in the project root directory.
2. Update any necessary values in the `.env` file to match your environment or preferences.


### Environment Variables

Below is a description of each environment variable used by the application:

| Variable                   | Description                                 | Example Value         |
|----------------------------|---------------------------------------------|-----------------------|
| `APP_NAME`                 | Name of the application                     | `device-manager`      |
| `SERVER_PORT`              | Port on which the server will run           | `8080`                |
| `HTTP_TIMEOUT_IN_SECONDS`  | HTTP request timeout in seconds             | `10`                  |
| `DATABASE_HOST`            | Hostname for the Postgres database          | `postgres`            |
| `DATABASE_PORT`            | Port for the Postgres database              | `5432`                |
| `DATABASE_USER`            | Username for the Postgres database          | `your_username`       |
| `DATABASE_PASS`            | Password for the Postgres database          | `your_password`       |
| `DATABASE_NAME`            | Name of the Postgres database               | `device-manager`      |


### Configure postgres database user and password

The same values used on `.env` file to set `DATABASE_USER` and `DATABASE_PASS` shall be used on the `compose.yml` file on the lines 6 and 7, where the Postgres User and Password is required to run in your local docker environment.

```yaml
POSTGRES_USER: your_username
POSTGRES_PASSWORD: your_password
```

To start or stop Postgres database, use commands described below.

```sh
make run/postgres   # Start the local Postgres on your Docker
make stop/postgres   # Start the local Postgres on your Docker
```


## Running with Docker

Build and start the application using Docker and Makefile:

The first step is run the command `build`, then use command `run` to start the application and then `logs` to see the logs of device manager api:

```sh
make build   # Builds the Docker image
make run   # Starts the containers
make logs    # Stops the containers
```

Or even:

```sh
make build run logs
```

If you need, it's possible to use `make clean` to erase the docker container to build again the image.

## All Makefile Commands

| Command               | Description                                                           |
|-----------------------|-----------------------------------------------------------------------|
| `make build`          | Build Docker image of the device manager API                          |
| `make run`            | Run API container                                                     |
| `make status`         | Show container status                                                 |
| `make logs`           | Tail logs of API container                                            |
| `make test`           | Run Go tests in Docker container showing % coverage                   |
| `make stop`           | Stop API container                                                    |
| `make clean`          | Remove the API containers                                             |
| `make run/postgres`   | Run Postgres as a dependency for the API                              |
| `make stop/postgres`  | Stop Postgres as a dependency for the API                             |


## API Documentation

Once the application is running, you can explore the API documentation and interact with the endpoints directly via Swagger UI at: http://localhost:8080/api/index.html

See [API Docs](docs/api.md) for endpoints and usage details.


## Testing the API

There's a [Postman Collection](docs/device-manager.postman_collection.json) and [Environment](docs/local.postman_environment.json) that can be imported to test each api endpoint.


## TO DOs

- Implement unit tests in network layer
- Implement integrated tests with newman
- Implement Api Key Authentication
- Refactoring unit tests to turn more simple and reusable
- Implement Rate Limit
- Implement Cache
- Migrate from lib/pq to pgx
- Implement pagination on List Devices endpoint



## License

APACHE
