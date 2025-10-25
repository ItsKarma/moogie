---
sidebar_position: 1
---

# Docker Setup

This guide walks you through setting up Moogie using Docker Compose, which is the recommended way to run Moogie in development and production.

## Prerequisites

Before you begin, ensure you have the following installed:

- **Docker** (version 20.10 or later)
- **Docker Compose** (version 2.0 or later)
- **Git** (for cloning the repository)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/ItsKarma/moogie.git
cd moogie
```

### 2. Start All Services

```bash
docker compose up --build
```

This command will:

- Build all Docker images
- Start PostgreSQL database
- Run database migrations
- Start the Go API server
- Start the Svelte web dashboard
- Start the documentation site

### 3. Access the Applications

Once all containers are running, you can access:

| Service           | URL                   | Description               |
| ----------------- | --------------------- | ------------------------- |
| **Dashboard**     | http://localhost:3000 | Main monitoring interface |
| **API**           | http://localhost:8080 | REST API and health check |
| **Documentation** | http://localhost:3000 | This documentation site   |
| **Database**      | localhost:5432        | PostgreSQL (internal)     |

## Services Overview

### PostgreSQL Database

- **Container**: `moogie-postgres`
- **Port**: 5432 (internal)
- **Credentials**: moogie/moogie
- **Volume**: Persistent data storage

### API Server

- **Container**: `moogie-api`
- **Port**: 8080
- **Features**: REST API, WebSocket, auto-migrations
- **Health Check**: http://localhost:8080/health

### Web Dashboard

- **Container**: `moogie-ui`
- **Port**: 3000
- **Features**: Responsive UI, real-time updates
- **Framework**: Svelte + Vite

### Documentation

- **Container**: `moogie-docs`
- **Port**: 3000
- **Features**: Interactive docs, API reference
- **Framework**: Docusaurus

## Environment Configuration

### Default Environment Variables

The following environment variables are configured by default:

```yaml
# API Configuration
APP_ENV: development
APP_PORT: 8080
DB_HOST: postgres
DB_PORT: 5432
DB_NAME: moogie
DB_USER: moogie
DB_PASSWORD: moogie

# CORS Configuration
ALLOWED_ORIGINS: http://localhost:3000
```

### Custom Configuration

To customize the environment, create a `.env` file in the project root:

```bash
# .env
APP_ENV=production
DB_PASSWORD=your_secure_password
ALLOWED_ORIGINS=https://your-dashboard.com,https://your-docs.com
```

## Development Workflow

### Starting Services

```bash
# Start all services
docker compose up

# Start in detached mode
docker compose up -d

# Build and start
docker compose up --build

# Start specific service
docker compose up postgres api
```

### Stopping Services

```bash
# Stop all services
docker compose down

# Stop and remove volumes
docker compose down -v

# Stop and remove images
docker compose down --rmi all
```

### Viewing Logs

```bash
# View all logs
docker compose logs

# View specific service logs
docker compose logs api
docker compose logs ui
docker compose logs docs

# Follow logs in real-time
docker compose logs -f api
```

### Database Operations

```bash
# Access PostgreSQL shell
docker compose exec postgres psql -U moogie -d moogie

# Run migrations manually
docker compose exec api goose -dir migrations postgres "host=postgres user=moogie password=moogie dbname=moogie sslmode=disable" up

# Create new migration
docker compose exec api goose -dir migrations create new_migration sql
```

## Troubleshooting

### Common Issues

#### Port Conflicts

If you encounter port conflicts, modify the ports in `docker-compose.yaml`:

```yaml
services:
  ui:
    ports:
      - "3000:3000" # UI accessible on port 3000

  docs:
    ports:
      - "3001:3000" # Docs accessible on port 3001
```

#### Database Connection Issues

Ensure PostgreSQL is ready before API starts:

```bash
# Check database health
docker compose exec postgres pg_isready -U moogie

# View API logs for connection errors
docker compose logs api
```

#### Build Failures

Clean Docker cache and rebuild:

```bash
# Remove all containers and volumes
docker compose down -v

# Remove images
docker compose down --rmi all

# Clean build cache
docker system prune -a

# Rebuild from scratch
docker compose up --build
```

### Health Checks

All services include health checks. Monitor service health:

```bash
# Check service status
docker compose ps

# View health check details
docker inspect moogie-api | grep -A 10 Health
```

## Production Deployment

### Security Considerations

For production deployment:

1. **Change default passwords**:

   ```yaml
   environment:
     DB_PASSWORD: ${DB_PASSWORD} # Use secrets
   ```

2. **Use external database**:

   ```yaml
   # Remove postgres service
   # Update API environment
   environment:
     DB_HOST: your-postgres-host.com
   ```

3. **Configure reverse proxy**:

   ```nginx
   # Nginx configuration
   server {
     listen 80;
     server_name your-domain.com;

     location / {
       proxy_pass http://localhost:3000;
     }

     location /api/ {
       proxy_pass http://localhost:8080;
     }
   }
   ```

4. **Use production builds**:
   ```dockerfile
   # In ui/Dockerfile
   ENV NODE_ENV=production
   RUN npm run build
   ```

### Scaling

Scale services based on load:

```bash
# Scale API instances
docker compose up --scale api=3

# Use load balancer for multiple instances
```

## Next Steps

- [API Reference](../api/overview) - Explore the REST API
- [Main Documentation](/) - Learn about the platform

Need help? [Open an issue](https://github.com/ItsKarma/moogie/issues) on GitHub.
