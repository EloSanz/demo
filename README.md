# Demo API (Go Port)

This project is a high-performance, modular Go port of an enterprise Spring Boot application, following Clean Architecture and DDD principles.

## Features

- **GORM Persistence**: Multi-DB support (PostgreSQL & SQLite) with automatic migrations.
- **Asynchronous Messaging**: Built-in support for **AWS SQS** and **In-Memory Go Channels** via a notification domain.
- **Observability**: Prometheus metrics at `/metrics` and structured logging with error tracing.
- **Resilience**: Graceful shutdown (20s timeout) and global recovery middleware.
- **Storage**: AWS S3 integration.
- **Modular Design**: Decoupled database, notification engine, and routing initialization.
- **Testing**: Full suite of unit and integration tests with HTTP mocking.

## Getting Started

### 1. Run with In-Memory Assets (Fastest)

```bash
# Uses SQLite and In-Memory Notifications by default
air
```

### 2. Run with Full AWS/Postgres Stack

```bash
# Set your environment variables
export NOTIFICATION_ENGINE=sqs
export AWS_SQS_QUEUE_URL="your-queue-url"
docker-compose up -d
```

### 3. Running Tests

```bash
go test ./... -v
```

## Monitoring & Infrastructure

- **Health Check**: `GET /health`
- **Metrics**: `GET /metrics`
- **Notifications**: `POST /api/notifications` (Async)

## Project Structure

```text
├── cmd/api/          # Orchestration entrypoint
├── internal/
│   ├── api/          # Routing Index & Middleware
│   ├── notification/ # Async Messaging Domain (SQS/Memory)
│   ├── database/     # DB Initialization
│   ├── [domain]/     # Logic, Handlers, and Repos
└── pkg/web/          # HTTP Utilities & Adapters
```
