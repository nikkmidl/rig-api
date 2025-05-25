# Rig Security Microservice

Rig API is a microservice designed to scan organizational repositories and validate allowed or blocked users according to the policy

## Features

- gRPC + gRPC Gateway RESTful API
- Easy integration with Docker

## Prerequisites

- [Docker](https://www.docker.com/get-started) (for containerized setup)
- [Go 1.24+](https://go.dev/doc/install) (for local setup)
- [Protobuf compiler](https://protobuf.dev/installation/)
- [Make](https://www.gnu.org/software/make/)

### Configure environment

Make sure to input all necessary environment variables in `app.env` file. Default values, except github api key, are present.

## Running with Docker

Run the container:

   
    docker compose up -d
   

The service will be available at `http://localhost:8080` if environment uses default values provided.

## Running Locally

1. Install dependencies:

    ```bash
    make install
    ```

2. Generate protobufs:

    ```bash
    make generate
    ```

3. Ensure all imports are satisfied:

    ```bash
    make tidy
    ```

4. Run the application:

    ```bash
    make run
    ```

The service will be available at `http://localhost:8080`.

## Configuration

Configuration options can be set via environment variables in the `app.env` file. Available options are listed in existing file.

## API Documentation

Swagger json specification is generated alongside protobufs and available in `/openapiv2` folder

## Example
From terminal execute:
```
curl http://localhost:8080/list/{your_organization}
```
response:
```
{"repos":[{"repoName":"sample-repo", "allowedUsers":["nikkmidl"], "blockedUsers":[]}]}
```