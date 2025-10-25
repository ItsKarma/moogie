# Moogie Database Seeder

A Go application that populates the Moogie PostgreSQL database with realistic sample data for testing and development.

## What it does

The seeder generates:

- **15 monitoring jobs** across 5 different types (HTTP, TCP, DNS, SSL, Ping)
- **30+ days of execution history** with realistic success rates and response times
- **Realistic configurations** using the same YAML-based job specifications as the real system

## Job Types Generated

### HTTP Checks

- API health endpoints
- Payment service monitoring
- Authentication service checks
- Web application monitoring

### TCP Checks

- Database connections (PostgreSQL, MySQL)
- Cache connections (Redis)
- Message queue connections (RabbitMQ)
- Search service connections (Elasticsearch)

### DNS Checks

- Domain resolution validation
- Internal DNS checks
- External DNS verification

### SSL Certificate Checks

- Certificate expiration monitoring
- SSL configuration validation
- Gateway and proxy SSL checks

### Ping Checks

- Server connectivity tests
- Load balancer health checks
- Network node monitoring

## Usage

### Using Docker Compose (Recommended)

```bash
# Run the seeder as a one-time job
docker compose run --rm seeder

# Or use the profile to run as part of the stack
docker compose --profile seed up
```

### Manual Build and Run

```bash
cd seed
go mod tidy
go build -o seeder main.go

# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=moogie
export DB_PASSWORD=moogie123
export DB_NAME=moogie

# Run the seeder
./seeder
```

## Environment Variables

| Variable      | Default     | Description       |
| ------------- | ----------- | ----------------- |
| `DB_HOST`     | `localhost` | PostgreSQL host   |
| `DB_PORT`     | `5432`      | PostgreSQL port   |
| `DB_USER`     | `moogie`    | Database username |
| `DB_PASSWORD` | `moogie123` | Database password |
| `DB_NAME`     | `moogie`    | Database name     |

## Generated Data

### Jobs

- 15 total jobs across all check types
- ~90% enabled (realistic production scenario)
- Randomized monitoring schedules (5-60 minute intervals)
- Realistic alert configurations with email notifications

### Executions

- 30 days of historical data
- Success rates between 70-99% (varies by job)
- Realistic response times:
  - Success: 50-500ms
  - Timeout failures: 5000-10000ms
  - Quick failures: 100-1000ms
- Detailed error messages and response metadata

## Data Cleanup

The seeder automatically clears existing data before seeding:

- Deletes all executions
- Deletes all jobs
- Resets auto-increment sequences

⚠️ **Warning**: This will delete all existing monitoring data!

## Development

The seeder uses:

- `github.com/lib/pq` for PostgreSQL connectivity
- `github.com/brianvoe/gofakeit/v6` for realistic fake data generation
- Standard library for JSON marshaling and time manipulation

## Sample Output

```
🌱 Starting Moogie database seeder...
✅ Connected to database successfully
🧹 Clearing existing data...
📝 Creating monitoring jobs...
  ✓ Created http job: api-health-check-0 (ID: 1)
  ✓ Created tcp job: database-connection-1 (ID: 2)
  ✓ Created dns job: primary-dns-resolution-2 (ID: 3)
⚡ Creating execution history...
  ✓ Created 1,247 executions for job ID 1
  ✓ Created 892 executions for job ID 2
  📊 Total executions created: 32,674
🎉 Database seeding completed successfully!
```
