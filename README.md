# Demo API (Go Port)

This project is a high-performance, modular Go port of an enterprise Spring Boot application, following Clean Architecture and DDD principles.

## Features

- **GORM Persistence**: Multi-DB support (PostgreSQL & SQLite) with automatic migrations.
- **Observability**: Built-in Prometheus metrics at `/metrics` and structured logging.
- **Resilience**: Graceful shutdown and global recovery middleware.
- **Storage**: AWS S3 integration.
- **Modular Design**: Decoupled database and routing initialization.
- **Testing**: Full suite of unit and integration tests with HTTP mocking.

## Requirements

- Go 1.25+
- Docker & Docker Compose (optional for full stack run)

## Getting Started

### 1. Fast Development (SQLite)

```bash
# Uses default config (SQLite)
air
```

### 2. Full Production Stack (Postgres)

```bash
docker-compose up -d
```

### 3. Running Tests

```bash
# Run all tests (Unit + Integration)
go test ./... -v
```

## Infrastructure & Monitoring

- **Health Check**: `GET /health` (Checks DB connectivity).
- **Metrics**: `GET /metrics` (Prometheus format).
- **Graceful Shutdown**: The server waits 20s for active requests before stopping.

## Project Structure

```text
├── cmd/api/          # Main entrypoint (Orchestration only)
├── internal/
│   ├── api/          # Routing Index & Middleware Assembly
│   ├── database/     # DB Initialization & GORM Config
│   ├── config/       # Structured Env Var Configuration
│   └── [domain]/     # Bounded Contexts (Logic + Handlers + Tests)
├── infrastructure/   # External Adapters (Repo, HTTP Clients, AWS)
└── pkg/web/          # HTTP Utilities (JSON, Error handling, Metrics middleware)
```
