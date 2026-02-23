# Spring Boot 4 Enterprise Template Project

A production-ready Spring Boot 4 template demonstrating clean architecture, DDD patterns, and production-grade features. Designed for technical challenges, microservices prototyping, and enterprise-grade applications.

## 🏛️ Architecture & Patterns

The project follows a **Clean Architecture** approach with clear separation of concerns:

```
Controller → Service (Interface) → Service (Implementation) → Repository (JPA)
    ↑                ↓                      ↓
OpenAPI / DTO    AOP Logging          External API Clients (HTTP Interfaces)
```

### Key Design Patterns:
- **DTO Pattern**: Clean separation between API contracts and internal data structures.
- **Interface-based Design**: Ensures loose coupling and easy mockability.
- **Repository Pattern**: Database-agnostic persistence layer.
- **Aspect-Oriented Programming (AOP)**: Modularized cross-cutting concerns (Logging).
- **HTTP Interfaces**: Declarative REST clients introduced in Spring Boot 3/4.
- **Strategy Pattern (via Spring Profiles)**: Easily switch between local H2 and production-ready PostgreSQL.

## 🎯 Features

- ✅ **Spring Boot 4**: Leveraging the latest declarative REST clients.
- ✅ **AOP Smart Logging**: Automated method-level logging (entry, exit, execution time, arguments) without polluting business logic.
- ✅ **JPA Auditing**: Automatic `createdAt` and `updatedAt` tracking via `BaseEntity`.
- ✅ **Pagination & Sorting**: Full support for paginated results in User API (`Pageable`).
- ✅ **Multi-API Integration**: Demonstrations with JSONPlaceholder and **Rick and Morty API**.
- ✅ **Database Agnostic**: H2 for development, PostgreSQL-ready with automated setup tasks.
- ✅ **OpenAPI 3.0 (Swagger)**: Comprehensive API documentation at `/swagger-ui.html`.
- ✅ **Infrastructure as Code**: `docker-compose.yml` for local environment consistency.
- ✅ **Code Coverage**: **JaCoCo** integration for unit and integration tests.
- ✅ **Quality Control**: Enforced via **Spotless** (formatting), **Checkstyle**, and **SpotBugs**.

## 📦 Tech Stack

- **Java 21** (LTS)
- **Spring Boot 4.0.2**
- **Spring Data JPA** (Hibernate 6)
- **MapStruct**: Type-safe bean mapping between Entities and DTOs.
- **Lombok**: Boilerplate reduction.
- **WireMock**: Reliable external API mocking for integration tests.
- **Gradle 8.14**

## 🚀 Getting Started

### Prerequisites
- Java 21+
- Docker (optional, for PostgreSQL)

### Run with H2 (In-memory)
```bash
./gradlew bootRun
```
Access at `http://localhost:8080`.

### Run with PostgreSQL (Local Automated)
Ensure you have PostgreSQL installed. This command creates the database if it doesn't exist and starts the app:
```bash
./gradlew bootRunPostgresLocal
```

### Run with Docker Compose
```bash
docker-compose up -d --build
```

## 📡 API Endpoints

### User Management (Local DB)
| Method | Endpoint | Description | Parameters |
|--------|----------|-------------|------------|
| GET | `/api/users` | Get paginated users | `page, size, sort` |
| GET | `/api/users/{id}` | Get user by ID | |
| POST | `/api/users` | Create user | `UserRequestDto` |
| PUT | `/api/users/{id}` | Update user | `UserRequestDto` |
| DELETE | `/api/users/{id}` | Delete user | |

### Rick and Morty Integration
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/rickandmorty/characters/{id}` | Fetch character from RM API |

### External Syncing
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/users/external` | List users from JSONPlaceholder |
| POST | `/api/users/sync/{id}` | Sync external user to local DB |

## 🔧 Quality & Testing

### Testing Strategy
- **Unit Tests**: Focus on business logic and mappers.
- **Integration Tests**: Full-slice tests using `BaseIntegrationTest`.
- **Mocking**: External APIs are mocked using **WireMock** via centralized test resources.

### Code Coverage (JaCoCo)
Coverage reports focus on **Business Logic** and **Controllers**. The following are excluded to keep results meaningful:
- `dto/**`, `domain/**` (POJOs/Entities)
- `mapper/**` (Generated code)
- `config/**`, `aspect/**`, `exception/**` (Infrastructure)
- `DemoApplication.java`

**Generate Reports:**
```bash
./gradlew jacocoFullReport     # Combined Unit + Integration report
```
Reports are available at `build/reports/jacoco/jacocoFullReport/html/index.html`.

### Quality Checks
```bash
./gradlew check              # Runs all checks (Test, Lints, Format, Coverage Verification)
./gradlew spotlessApply      # Fix formatting issues automatically
```

## 🛡️ Git Hooks & CI

### Pre-commit Hooks
This project uses `pre-commit` to ensure code quality before every commit.
1. **Install pre-commit**: `pip install pre-commit` (or `brew install pre-commit`)
2. **Install hooks**: `pre-commit install`

Hooks configured:
- **commit**: Runs `spotlessApply` to ensure formatting.
- **push**: Runs full `./gradlew check` to prevent breaking the build.

### Continuous Integration (GitHub Actions)
A CI pipeline is configured in `.github/workflows/ci.yml` that:
1. Validates code on every Push and Pull Request to `main`.
2. Runs the full build and all tests.
3. Generates and uploads the **JaCoCo Coverage Report** as a build artifact.

## 📂 Project Structure

```
src/main/java/com/example/demo/
├── aspect/              # Aspect Oriented Programming (Logging)
├── client/              # Declarative HTTP REST Clients
├── config/              # Infrastructure & OpenApi config
├── controller/          # REST Endpoints
├── domain/              # JPA Entities & BaseEntity
├── dto/                 # API Data Transfer Objects
├── exception/           # Global Exception Handling
├── mapper/              # MapStruct interfaces
├── repository/          # JPA Repositories
└── service/             # Business Logic (Interface + Impl)
```

## 📄 License
Template for educational and enterprise prototyping purposes.
