# Demo API (Go Port)

This project is a Go port of an enterprise-grade Spring Boot application, maintaining Clean Architecture and DDD principles, now powered by **GORM**.

## Features

- **Users Domain**: CRUD operations with **GORM** persistence.
- **Auto-Migrations**: Database schema managed automatically by GORM.
- **Rick and Morty Integration**: Consume external API with filtering and pagination.
- **External Sync**: Fetch users from JSONPlaceholder and upsert into local database.
- **Storage Service**: AWS S3 integration for file upload/download/delete/list.
- **Thin Web Wrapper**: Custom `pkg/web` following internal infrastructure conventions.
- **Developer Experience**: Hot-reload with `air` and structured request logging.
- **Testing**: Comprehensive Unit and **Integration Tests** (using `httptest` and in-memory SQLite).

## Requirements

- Go 1.25+
- Docker (optional)

## Getting Started

### 1. Environment Configuration

- `PORT`: Server port (default: `8080`).
- `DATABASE_PATH`: Path to SQLite file (default: `demo.db`).
- `AWS_S3_BUCKET`: S3 Bucket name (default: `myawsbucketelito`).

### 2. Run Local Development

```bash
air
```

### 3. Running Tests

Execute both unit and integration tests:
```bash
go test ./... -v
```

## Project Structure

```text
├── cmd/api/          # Application entrypoint (GORM Init & AutoMigrate)
├── infrastructure/   # External adapters (GORM Repository, HTTP Clients, AWS)
├── internal/         # Bounded Contexts (User, RickAndMorty, Storage)
├── pkg/web/          # HTTP Infrastructure wrapper
└── tests/            # Integration logic (within handler packages)
```
