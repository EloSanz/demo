# Demo API (Go Port)

This project is a Go port of an enterprise-grade Spring Boot application, maintaining Clean Architecture and DDD principles.

## Features

- **Users Domain**: CRUD operations with SQLite local persistence.
- **Rick and Morty Integration**: Consume external API with filtering and pagination.
- **External Sync**: Fetch users from JSONPlaceholder and upsert into local database.
- **Storage Service**: AWS S3 integration for file upload/download/delete/list.
- **Thin Web Wrapper**: Custom `pkg/web` following internal infrastructure conventions.
- **Developer Experience**: Hot-reload with `air` and structured request logging.

## Requirements

- Go 1.25+
- Docker (optional)

## Getting Started

### 1. Environment Configuration

The following optional environment variables are supported:

- `PORT`: Server port (default: `8080`).
- `DATABASE_PATH`: Path to SQLite file (default: `:memory:`).
- `AWS_S3_BUCKET`: S3 Bucket name (default: `myawsbucketelito`).
- `AWS_REGION`: AWS Region (reads from `~/.aws/config` if not set).

### 2. Run Local Development (with Hot-Reload)

First, install `air`:
```bash
go install github.com/air-verse/air@latest
```

Then start the application:
```bash
air
```

### 3. Run with Docker

Build the image:
```bash
docker build -t demo-go .
```

Run the container:
```bash
docker run -p 8080:8080 demo-go
```

## Running Tests

To execute the unit tests (currently 14 test cases):
```bash
go test ./... -v
```

## API Documentation

A Postman collection is available at [postman_collection.json](./postman_collection.json).

## Project Structure

```text
├── cmd/api/          # Application entrypoint
├── infrastructure/   # External adapters (DB, HTTP Clients, AWS)
├── internal/         # Bounded Contexts (User, RickAndMorty, Storage)
├── pkg/web/          # HTTP Infrastructure wrapper
└── migrations/       # SQL migrations
```

## License

MIT
