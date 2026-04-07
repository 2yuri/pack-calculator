# Pack Calculator

## Project Structure

```
.
├── backend/                  # Go Service
│   ├── cmd/                  # Entry point (main.go)
│   ├── internal/
│   │   ├── api/              # HTTP Layer
│   │   │   └── handler/
│   │   │       ├── product/  # Products CRUD endpoints
│   │   │       ├── pack/     # Packs CRUD endpoints
│   │   │       ├── order/    # Order calculation endpoint
│   │   │       └── swagger/  # OpenAPI docs
│   │   ├── domain/           # Models, interfaces and errors
│   │   ├── service/          # Business logic layer
│   │   ├── repository/       # Data access (PostgreSQL + Redis)
│   │   └── mocks/            # Generated test mocks (mockgen)
│   ├── pkg/
│   │   ├── config/           # Environment-based configuration
│   │   ├── db/               # PostgreSQL connection
│   │   ├── cache/            # Redis connection
│   │   └── rest/             # HTTP response helpers
│   └── migrations/           # SQL migrations (golang-migrate)
├── frontend/                 # React SPA
```

## Architecture Decisions

**Layered architecture with domain interfaces.** The `domain` package defines all models and interfaces. Services implement business logic, repositories handle data access. Handlers depend on services, services depend on repositories, nothing depends on the HTTP or database layer directly.

**Dependency injection via constructor structs.** Each layer receives its dependencies through a `Deps` struct passed to the constructor, making testing easier with generated mocks (`mockgen`).

**Mockgen as code generator.** To be honest i use this package since 2021 when i write a article about it. In that time it was not maintained by Uber, but since then i think it is a good tool.

**Redis as a pack cache.** Pack sizes are cached in Redis and synced on startup via `SyncCache()`. My idea was to use a cache to avoid hitting PostgreSQL on every request and faster http response.

**Dynamic programming for pack calculation.** The order service uses a DP algorithm to find the optimal pack combination. It minimizes total items shipped (waste) first, then minimizes pack count as a secondary objective. Handles edge cases like all packs being larger than the order quantity.

**Soft deletes.** Products and packs use a `deleted_at` for soft deletes.

## Stack

| Layer          | Technology                       |
| -------------- | -------------------------------- |
| HTTP           | Echo v4                          |
| Logging        | Zap                              |
| Database       | PostgreSQL 16 and golang-migrate |
| Testing        | mockgen and testify              |
| Cache          | Redis 7                          |
| Infrastructure | Docker Compose                   |

## Makefile Commands

| Command                 | Description                                 |
| ----------------------- | ------------------------------------------- |
| `make run`              | Start the application                       |
| `make build`            | Build the application                       |
| `make test`             | Run all tests                               |
| `make test-integration` | Run integration tests                       |
| `make generate-mocks`   | Generate mocks for all domain interfaces    |
| `make migrate-up`       | Run database migrations                     |
| `make migrate-down`     | Rollback database migrations                |
| `make docker-up`        | Start the application in a Docker container |
| `make docker-down`      | Stop the Docker container                   |

## Testing

```
          ╱  ╲
         ╱    ╲           E2E
        ╱ E2E  ╲          - Playwright — full user flows
       ╱────────╲
      ╱          ╲        Integration
     ╱Integration ╲       - testify + real PostgreSQL/Redis
    ╱──────────────╲
   ╱                ╲     Unit
  ╱   Unit  Tests    ╲    - testify + mockgen
 ╱────────────────────╲
```

### Unit tests

Unit tests are written with `testify` and `mockgen`. They cover two layers:

- **Service layer** (`internal/service/`) — Each service method is tested with mocked repository dependencies using `gomock`. Tests verify business logic in isolation.
- **Handler layer** (`internal/api/handler/`) — Table-driven tests with mocked services. Each test case defines its own mock setup, expected HTTP status code, and response body.

```bash
make test
```

### Integration tests

Integration tests are written with `testify` and live in the repository layer (`internal/repository/`). They run against real PostgreSQL and Redis instances, using helper functions (`setupTestDB`, `setupTestRedis`) and `t.Cleanup()` to manage test data.

Requires `TEST_DATABASE_URL` and `TEST_REDIS_ADDR` environment variables.

```bash
make test-integration
```

### E2E tests

End-to-end tests use [Playwright](https://playwright.dev/) and run against the full stack (frontend + backend + database). They exercise real user flows through the browser.

```bash
cd frontend && pnpm test:e2e
```

### Mock generation

Mocks are generated from domain interfaces using `mockgen` and stored in `internal/mocks/`.

```bash
make generate-mocks
```

## Deploy information

- So i have deployed this application on my personal VPS. I am using Docker to host the application, Redis and PostgreSQL.
- I thought about using Kubernetes, but i think it is overkill for this application challenge.
