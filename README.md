# Moogie

Synthetics Runner with Real-time Dashboard

## Overview

Moogie is a synthetics monitoring system that runs native in Kubernetes utilizing Cron Jobs to execute your synthetic checks. Instead of configuring checks through a UI, Moogie uses YAML configuration files that follow Kubernetes-style resource definitions.

The system consists of:

- **📊 Dashboard UI** - Svelte-based monitoring dashboard
- **🔌 API Server** - Go/Gin REST API with WebSocket support
- **🗄️ Database** - PostgreSQL for storing job configs and execution results
- **🔄 Runner** - Kubernetes CronJobs that execute checks and report results

## Quick Start with Docker

### Prerequisites

- Docker & Docker Compose installed
- 8080, 3000, and 5432 ports available

### Start the Full Stack

1. **Clone and navigate to the project:**

```bash
git clone <repository-url>
cd moogie
```

2. **Start everything with Docker Compose:**

```bash
# Start the complete stack (database + api + ui)
docker-compose up --build

# Or run in background
docker-compose up -d --build
```

3. **Access the application:**

- 🌐 **Dashboard**: http://localhost:3000
- 🔌 **API**: http://localhost:8080/api/v1
- ❤️ **Health Check**: http://localhost:8080/health
- 📡 **WebSocket**: ws://localhost:8080/ws

### Optional: Database Admin Interface

To access PostgreSQL admin interface:

```bash
# Start with pgAdmin included
docker-compose --profile admin up -d

# Access pgAdmin at http://localhost:5050
# Email: admin@moogie.local
# Password: admin
```

### Development Commands

```bash
# View logs
docker-compose logs -f api
docker-compose logs -f ui

# Rebuild specific service
docker-compose up --build api
docker-compose up --build ui

# Stop everything
docker-compose down

# Stop and remove volumes (clean slate)
docker-compose down -v
```

### Database Access

PostgreSQL is accessible at:

- **Host**: localhost:5432
- **Database**: moogie
- **User**: moogie
- **Password**: moogie

### Seed Sample Data

To populate the database with realistic sample data for testing:

```bash
# Run the seeder to create 15 jobs with 30 days of execution history
docker-compose run --rm seeder

# Or include seeding when starting the stack
docker-compose --profile seed up -d
```

The seeder creates:

- 15 monitoring jobs across all check types (HTTP, TCP, DNS, SSL, Ping)
- 30+ days of execution history with realistic success rates (70-99%)
- Over 30,000 execution records with varied response times and error conditions

⚠️ **Warning**: The seeder clears all existing data before creating new sample data.

## Configuration

### Check Configuration Files

All check configurations are stored as YAML files in the `config/checks/` directory. These files follow a Kubernetes-style structure with three main sections:

- **apiVersion**: Specifies the API version (currently `moogie.io/v1`)
- **kind**: The type of check (HttpCheck, TcpCheck, DnsCheck, SslCheck, PingCheck)
- **metadata**: Contains the check name, labels, and other metadata
- **spec**: Defines the actual check configuration and parameters

### Supported Check Types

#### HttpCheck

For HTTP/HTTPS endpoint monitoring:

```yaml
apiVersion: moogie.io/v1
kind: HttpCheck
metadata:
  name: api-health-check
  labels:
    environment: production
    service: api
spec:
  url: https://api.example.com/health
  method: GET
  timeout: 30s
  expectedStatusCode: 200
  schedule: "*/5 * * * *" # Every 5 minutes
  retries: 3
  headers:
    User-Agent: "Moogie/1.0"
  alerts:
    onFailure: true
    email: alerts@example.com
```

#### TcpCheck

For TCP port connectivity checks:

```yaml
apiVersion: moogie.io/v1
kind: TcpCheck
metadata:
  name: database-connection
spec:
  host: db.example.com
  port: 5432
  timeout: 10s
  schedule: "*/2 * * * *" # Every 2 minutes
```

#### DnsCheck

For DNS resolution monitoring:

```yaml
apiVersion: moogie.io/v1
kind: DnsCheck
metadata:
  name: dns-resolution-check
spec:
  domain: example.com
  recordType: A
  expectedIp: "192.0.2.1"
  timeout: 5s
  schedule: "*/10 * * * *" # Every 10 minutes
```

#### SslCheck

For SSL certificate monitoring:

```yaml
apiVersion: moogie.io/v1
kind: SslCheck
metadata:
  name: ssl-certificate-check
spec:
  host: example.com
  port: 443
  daysBeforeExpiry: 30 # Alert when cert expires in 30 days
  schedule: "0 0 * * *" # Daily at midnight
```

#### PingCheck

For basic connectivity checks:

```yaml
apiVersion: moogie.io/v1
kind: PingCheck
metadata:
  name: server-ping-check
spec:
  host: server.example.com
  count: 3
  timeout: 5s
  schedule: "*/1 * * * *" # Every minute
```

### Adding New Checks

1. Create a new YAML file in the `config/checks/` directory
2. Follow the structure above for your desired check type
3. Commit the file to the repository
4. Moogie will automatically detect and start running the new check

### File Naming Convention

Use descriptive names for your check files:

- `api-health-check.yaml`
- `database-tcp-check.yaml`
- `ssl-certificate-monitor.yaml`

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Dashboard UI  │    │   API Server    │    │   PostgreSQL    │
│   (Svelte)      │◄──►│   (Go/Gin)      │◄──►│   Database      │
│   Port 3000     │    │   Port 8080     │    │   Port 5432     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              ▲
                              │ POST /api/v1/executions
                              ▼
                       ┌─────────────────┐
                       │ Runner Service  │
                       │ (K8s CronJobs)  │
                       └─────────────────┘
```

## Development

### Individual Services

- **UI Development**: See `ui/README.md`
- **API Development**: See `api/README.md`
- **Runner Development**: See `runner/README.md`

### Sample Data

The Docker setup includes sample monitoring jobs:

- API Health Check
- Database TCP Check
- DNS Resolution Check
- Ping Connectivity Check
- SSL Certificate Check

## Troubleshooting

### Common Issues

**Container won't start:**

```bash
# Check logs
docker-compose logs api
docker-compose logs ui
docker-compose logs postgres
```

**Database connection errors:**

```bash
# Ensure PostgreSQL is ready
docker-compose ps postgres

# Check if migrations ran
docker-compose exec api goose -dir migrations postgres "postgres://moogie:moogie@postgres:5432/moogie?sslmode=disable" status
```

**Port conflicts:**

```bash
# Check what's using the ports
lsof -i :3000
lsof -i :8080
lsof -i :5432

# Use different ports if needed
docker-compose up -e "API_PORT=8081" -e "UI_PORT=3001"
```

**Clean restart:**

```bash
# Stop everything and remove volumes
docker-compose down -v

# Remove images and rebuild
docker-compose build --no-cache
docker-compose up
```

## Deploying

This runs native in Kubernetes utilizing Cron Jobs to run your Synthetic Checks.
