# Backend

Go service for the Pack Calculator API.

## Prerequisites

- Go 1.24+
- Docker and Docker Compose
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI (for running migrations manually)

## Getting Started

### 1. Configure environment

Copy the local environment file:

```bash
cp .env.local .env
```

| Variable            | Description               | Default       |
| ------------------- | ------------------------- | ------------- |
| `APP_PORT`          | HTTP server port          | `8080`        |
| `DB_HOST`           | PostgreSQL host           | `localhost`   |
| `DB_PORT`           | PostgreSQL port           | `5432`        |
| `DB_NAME`           | Database name             | `pack`        |
| `DB_USER`           | Database user             | `pack`        |
| `DB_PASSWORD`       | Database password         | `awesomepass` |
| `CACHE_HOST`        | Redis host                | `localhost`   |
| `CACHE_PORT`        | Redis port                | `6379`        |
| `CACHE_PASSWORD`    | Redis password (optional) | empty         |
| `APP_ALLOW_ORIGINS` | CORS allowed origins      | `*`           |

### 2. Start dependencies

```bash
make docker-up
```

This starts PostgreSQL 16 and Redis 7 containers.

### 3. Run migrations

```bash
make migrate-up
```

To rollback:

```bash
make migrate-down
```

### 4. Run the application

```bash
make run
```

The API will be available at `http://localhost:8080`. OpenAPI docs are served at `/swagger/`.

## Testing

### Unit tests

```bash
make test
```

### Integration tests

Requires running PostgreSQL and Redis (via `make docker-up`):

```bash
make test-integration
```

### Mock generation

Regenerate mocks after changing domain interfaces:

```bash
make generate-mocks
```
