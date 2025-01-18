[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=pangolin-do-golang_thumb-processor-workeri&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=pangolin-do-golang_tech-challenge) ![Known Vulnerabilities](https://snyk.io/test/github/pcbarretos/pangolin-do-golang/tech-challenge/badge.svg)

[Coverage Report](https://sonarcloud.io/summary/overall?id=pangolin-do-golang_thumb-processor-worker)


# Hackaton - Thumb Processor Worker

## Install

### Go

- [Go Install](https://go.dev/doc/install)

> Make sure you have Go 1.22.2 or higher

Execute

```shell
go mod tidy
```

## Run tests

```shell
go test -cover ./...
```

## Defining Envs

To correctly use the project, it is necessary to define a .env file, with the values for the envs:

* DB_USERNAME
* DB_PASSWORD
* DB_HOST
* DB_NAME
* DB_PORT
* THUMB_SERVICE_URL

We recommend using for development the following values:

```
DB_USERNAME=user
DB_PASSWORD=pass
DB_HOST=pgsql
DB_NAME=postgres
DB_PORT=5432
THUMB_SERVICE_URL=http://localhost:<PORT_SERVICE_IS_RUNNING>
```

## Executing with Docker (Compose)

```shell
docker compose build

docker compose up -d

curl --request GET --url http://localhost:8082/health

## Expected response
= Status Code 200
```

> If you're having trouble deploying the application with `docker-compose` (and not `docker compose`), use docker version 27.0.0 or higher.

## Accessing Swagger UI

Go to http://localhost:8080/swagger/index.html#/ after the application is running.

## Stack

- [Go](https://go.dev/)
- [Gin Web Framework](https://gin-gonic.com/) - Routes, JSON validation, Error management, Middleware support
- [PostgresSQL](https://www.postgresql.org/) - Database
- [swag](https://github.com/swaggo/swag) - Tool to generate swagger documentation
- [docker](https://www.docker.com/) - Containerization tool
- [docker-compose](https://docs.docker.com/compose/) - Tool to define and run multi-container Docker applications

## Swagger

This project makes use of the library [swag](https://github.com/swaggo/swag?tab=readme-ov-file#how-to-use-it-with-gin) to generate the swagger documentation.

### Install

Follow the steps described in the [official documentation](https://github.com/swaggo/swag?tab=readme-ov-file#getting-started)

### Generate

```shell
 swag init -g cmd/sqs/main.go 
```

### Access the documentation

The documentation can be founded at the path `/docs/swagger.yaml` or accessing this [link](./docs/swagger.yaml).

## Project structure

- `cmd`: Application entry point directory for the application's main entry points, dependency injection, or commands. The web subdirectory contains the main entry point to the REST API.
- `internal`: Directory to contain application code that should not be exposed to external packages.
    - `core`: Directory that contains the application's core business logic.
        - `thumb`: Directory contains definition of the entity's heights, interfaces, repository and service of the entity Thumb.
    - `adapters`: Directory to contain external services that will interact with the application core.
        - `db`: Directory contains the implementation of the repositories.
        - `rest`: Directory that contains the definition of the application's controllers and handlers for manipulating data provided by the controller
    - `domainerrors`: Directory that contains the definition of the application's domain errors.