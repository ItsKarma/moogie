# Moogie API

A monitoring dashboard API built with Go, Gin, PostgreSQL, and WebSocket support for real-time updates.

## Features

- 🚀 **RESTful API** - Built with Gin web framework
- 📊 **PostgreSQL Database** - With GORM ORM and Goose migrations
- 📡 **WebSocket Support** - Real-time updates for dashboard
- 📖 **OpenAPI/Swagger Docs** - Auto-generated API documentation
- 🔄 **Hot Reload** - Development mode with Air
- 🐳 **Docker Ready** - Containerized deployment

## Quick Start

### Prerequisites

- Go 1.22+
- PostgreSQL 13+
- Make (optional, for convenience commands)

### Environment Setup

1. Copy the example environment file:

```bash
cp .env.example .env
```

2. Update `.env` with your database credentials:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=moogie
DB_PASSWORD=moogie
DB_NAME=moogie
```

### Database Setup

1. Create the database:

```sql
CREATE DATABASE moogie;
CREATE USER moogie WITH PASSWORD 'moogie';
GRANT ALL PRIVILEGES ON DATABASE moogie TO moogie;
```

2. Run migrations:

```bash
make migrate-up DB_URL="postgres://moogie:moogie@localhost:5432/moogie?sslmode=disable"
```

### Development

1. Install dependencies:

```bash
make deps
```

2. Run in development mode (with hot reload):

```bash
make dev
```

3. Or build and run:

```bash
make build
make run
```

The API will be available at `http://localhost:8080`

## API Endpoints

### Jobs

- `GET /api/v1/jobs` - List all jobs with metrics
- `GET /api/v1/jobs/:id` - Get job details with execution history

### Executions

- `POST /api/v1/executions` - Create execution result (called by runner)

### Dashboard

- `GET /api/v1/dashboard/summary` - Get dashboard summary metrics

### WebSocket

- `GET /ws` - WebSocket endpoint for real-time updates

### Health

- `GET /health` - Health check endpoint

## Date Range Filtering

Most endpoints support date range filtering with query parameters:

- `from` - Start date (YYYY-MM-DD format)
- `to` - End date (YYYY-MM-DD format)

Example: `GET /api/v1/jobs?from=2025-10-15&to=2025-10-22`

## WebSocket Messages

The WebSocket endpoint (`/ws`) broadcasts real-time updates:

```json
{
  "type": "execution_created",
  "data": {
    "id": 123,
    "job_id": 1,
    "status": "success",
    "response_time": 250,
    "timestamp": "2025-10-22T10:30:00Z"
  }
}
```

Message types:

- `execution_created` - New execution result
- `job_updated` - Job configuration updated
- `dashboard_updated` - Dashboard metrics updated

## Database Schema

### Jobs Table

```sql
CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(100) NOT NULL,
    config JSONB NOT NULL,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Executions Table

```sql
CREATE TABLE executions (
    id SERIAL PRIMARY KEY,
    job_id INTEGER NOT NULL REFERENCES jobs(id),
    status VARCHAR(20) NOT NULL CHECK (status IN ('success', 'failure')),
    response_time INTEGER DEFAULT 0,
    details JSONB,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

## Creating Execution Results

The runner service posts execution results to the API:

```bash
curl -X POST http://localhost:8080/api/v1/executions \
  -H "Content-Type: application/json" \
  -d '{
    "job_name": "example-ping-check",
    "status": "success",
    "response_time": 25,
    "details": {"host": "example.com", "packet_loss": 0},
    "timestamp": "2025-10-22T10:30:00Z"
  }'
```

## Development Commands

```bash
# Download dependencies
make deps

# Build the application
make build

# Run the application
make run

# Development mode with hot reload
make dev

# Run tests
make test

# Database migrations
make migrate-up DB_URL="postgres://..."
make migrate-down DB_URL="postgres://..."
make migrate-create NAME="add_new_field"

# Generate API documentation
make docs

# Clean build artifacts
make clean
```

## API Documentation

Generate and view Swagger documentation:

```bash
make docs
```

Then visit `http://localhost:8080/swagger/index.html` when the server is running.

## Configuration

Environment variables:

| Variable          | Description             | Default                 |
| ----------------- | ----------------------- | ----------------------- |
| `APP_ENV`         | Application environment | `development`           |
| `APP_PORT`        | Server port             | `8080`                  |
| `DB_HOST`         | Database host           | `localhost`             |
| `DB_PORT`         | Database port           | `5432`                  |
| `DB_USER`         | Database user           | `moogie`                |
| `DB_PASSWORD`     | Database password       | `moogie`                |
| `DB_NAME`         | Database name           | `moogie`                |
| `DB_SSLMODE`      | Database SSL mode       | `disable`               |
| `ALLOWED_ORIGINS` | CORS allowed origins    | `http://localhost:3000` |

## Project Structure

```
api/
├── cmd/server/          # Application entry point
├── internal/
│   ├── handlers/        # HTTP request handlers
│   ├── models/          # Database models and DTOs
│   ├── services/        # Business logic layer
│   └── websocket/       # WebSocket hub implementation
├── pkg/
│   ├── config/          # Configuration management
│   └── database/        # Database connection
├── migrations/          # Database migration files
├── docs/                # Generated API documentation
└── Makefile            # Development commands
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Run `make test` to verify
6. Submit a pull request

## License

This project is licensed under the MIT License.
