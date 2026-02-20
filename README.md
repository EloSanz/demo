# Spring Boot 4 Template Project

A production-ready Spring Boot 4 template with clean architecture, designed for take-home challenges and rapid prototyping.

## 🏗️ Architecture

```
Controller → Service (Interface) → Service (Implementation) → Repository
                ↓
          External API Client (HTTP Interface)
```

## 🎯 Features

- ✅ **Clean Architecture**: Controller → Service → Repository pattern
- ✅ **Interface-based Design**: Easy to mock and test
- ✅ **Database Agnostic**: JPA-based, swap databases by changing config
- ✅ **Declarative HTTP Clients**: Spring Boot 4 HTTP Interfaces (no RestTemplate boilerplate)
- ✅ **Configuration from YAML**: External API URLs configured in `application.yml`
- ✅ **Validation**: Built-in with `@Valid` annotations
- ✅ **Exception Handling**: Global exception handler
- ✅ **Lombok**: Reduced boilerplate code
- ✅ **Logging**: SLF4J with Logback

## 📦 Tech Stack

- **Java 21** (LTS)
- **Spring Boot 4.0.2**
- **Spring Data JPA** (database abstraction)
- **H2 Database** (in-memory, easily swappable)
- **Lombok** (reduce boilerplate)
- **Gradle 8.14**

## 🚀 Quick Start

### Prerequisites
- Java 21 installed
- Gradle 8.14+ (or use the wrapper `./gradlew`)

### Run the application (H2)
```bash
./gradlew bootRun
```
The application will start on `http://localhost:8080`.

### Run with PostgreSQL (Local)
If you have PostgreSQL 17 installed on your Mac/Windows, you can use the automated task that ensures the database exists before starting:
```bash
./gradlew bootRunPostgresLocal --console=plain
```
*Note: The `--console=plain` flag is recommended to keep the terminal clean while the app is running.*

### Run with Docker Compose
To start the entire stack (App + PostgreSQL 17) in containers:
```bash
docker-compose up -d --build
```

### Access H2 Console
```
URL: http://localhost:8080/h2-console
JDBC URL: jdbc:h2:mem:testdb
Username: sa
Password: (leave empty)
```

## 📡 API Endpoints

### Local Database Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/users` | Get all users from local DB |
| GET | `/api/users/{id}` | Get user by ID from local DB |
| POST | `/api/users` | Create new user |
| PUT | `/api/users/{id}` | Update user |
| DELETE | `/api/users/{id}` | Delete user |

### External API Integration

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/users/external` | Fetch users from JSONPlaceholder API |
| POST | `/api/users/sync/{id}` | Sync user from external API to local DB |

## 📝 Example Requests

### Create a User
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "555-1234",
    "website": "johndoe.com"
  }'
```

### Get All Users
```bash
curl http://localhost:8080/api/users
```

### Fetch from External API
```bash
curl http://localhost:8080/api/users/external
```

### Sync User from External API
```bash
curl -X POST http://localhost:8080/api/users/sync/1
```

## 🔧 Configuration

### Database Configuration (`application.yml`)

**Current (H2 in-memory):**
```yaml
spring:
  datasource:
    url: jdbc:h2:mem:testdb;DB_CLOSE_DELAY=-1
    driverClassName: org.h2.Driver
    username: sa
    password:
```

**Switch to PostgreSQL:**
```yaml
spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/mydb
    driverClassName: org.postgresql.Driver
    username: postgres
    password: yourpassword
```

Then add dependency in `build.gradle`:
```gradle
runtimeOnly 'org.postgresql:postgresql'
```

### External API Configuration

Add your API base URLs in `application.yml`:
```yaml
external:
  api:
    your-api-name:
      base-url: https://api.example.com
```

Then create an HTTP Interface:
```java
@HttpExchange
public interface YourApiClient {
    @GetExchange("/endpoint")
    YourResponse getData();
}
```

Configure the bean in `HttpClientConfig.java`.

## 📂 Project Structure

```
src/main/java/com/example/demo/
├── client/              # HTTP Interface clients for external APIs
│   └── ExternalUserClient.java
├── config/              # Configuration classes
│   └── HttpClientConfig.java
├── controller/          # REST controllers
│   └── UserController.java
├── domain/              # JPA entities
│   └── User.java
├── dto/                 # Data Transfer Objects
│   ├── UserRequest.java
│   └── UserResponse.java
├── exception/           # Exception handlers
│   └── GlobalExceptionHandler.java
├── repository/          # JPA repositories
│   └── UserRepository.java
└── service/             # Business logic
    ├── UserService.java (interface)
    └── impl/
        └── UserServiceImpl.java
```

## 🧪 Testing

Run tests:
```bash
./gradlew test
```

## 🛠️ Gradle Commands

Here are the most common and useful commands for working with this project:

| Command | Description |
|---------|-------------|
| `./gradlew bootRun` | Starts the application locally using H2. |
| `./gradlew bootRunPostgresLocal` | Ensures local DB exists and starts the app with PostgreSQL. |
| `./gradlew build` | Compiles code, runs all tests, executes linters, and builds the executable JAR. |
| `./gradlew clean` | Deletes the `build/` directory, wiping all generated files. |
| `./gradlew test` | Runs the unit tests. |
| `./gradlew integrationTest` | Runs the integration tests. |
| `./gradlew spotlessApply` | Auto-formats the code (Google Java Format, removes unused imports, etc.). |
| `./gradlew spotlessCheck` | Verifies if the code is correctly formatted without modifying the files. |
| `./gradlew spotbugsMain` | Runs static code analysis to detect bugs and bad practices in the main source code. |
| `./gradlew classes` | Compiles the main Java source code to ensure there are no compilation errors. |

## 🎓 Key Concepts for Take-Home Challenges

### 1. **Database Agnostic Design**
- Uses JPA, so you can swap H2 for PostgreSQL/MySQL/MongoDB easily
- Just change `application.yml` and add the driver dependency

### 2. **Declarative HTTP Clients (Spring Boot 4)**
- No more `RestTemplate` or `WebClient` boilerplate
- Define interfaces with `@HttpExchange` annotations
- URLs come from `application.yml`

### 3. **Interface-Based Services**
- Easy to mock for testing
- Follows SOLID principles
- Clean separation of concerns

### 4. **DTO Pattern**
- Separate request/response DTOs from domain entities
- Validation on DTOs, not entities
- Clean API contracts

## 📚 Adding New Features

### Add a new entity (e.g., Product)

1. **Create domain entity**: `domain/Product.java`
2. **Create DTOs**: `dto/ProductRequest.java`, `dto/ProductResponse.java`
3. **Create repository**: `repository/ProductRepository.java`
4. **Create service interface**: `service/ProductService.java`
5. **Create service implementation**: `service/impl/ProductServiceImpl.java`
6. **Create controller**: `controller/ProductController.java`

### Add external API client

1. **Add base URL** in `application.yml`:
   ```yaml
   external:
     api:
       my-api:
         base-url: https://api.example.com
   ```

2. **Create HTTP Interface**: `client/MyApiClient.java`
   ```java
   @HttpExchange
   public interface MyApiClient {
       @GetExchange("/data")
       MyResponse getData();
   }
   ```

3. **Configure bean** in `HttpClientConfig.java`

## 🛠️ Customization Tips

- **Change port**: Add `server.port: 8081` in `application.yml`
- **Enable CORS**: Add `@CrossOrigin` on controllers
- **Add Swagger**: Add `springdoc-openapi-starter-webmvc-ui` dependency
- **Add security**: Add `spring-boot-starter-security` dependency

## 📄 License

This is a template project for educational and take-home challenge purposes.
