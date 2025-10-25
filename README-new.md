# Moogie - Monitoring Dashboard

A comprehensive monitoring dashboard with real-time updates, built with Go, PostgreSQL, and Svelte.

## 🚀 Quick Start

### Option 1: Docker Compose (Recommended)

Start the entire stack with one command:

```bash
./start.sh
```

This starts:

- 📊 **PostgreSQL** database (port 5432)
- 🔗 **Go API** server (port 8080)
- 🎨 **Svelte UI** (port 3000)

**Access the dashboard:** http://localhost:3000

### Option 2: With Database Admin

```bash
./start.sh --admin
```

Includes **PgAdmin** at http://localhost:5050

- Email: `admin@moogie.local`
- Password: `admin`

### Option 3: Background Mode

```bash
./start.sh -d
```

Runs in background. View logs with:

```bash
docker-compose logs -f
```

Stop with:

```bash
docker-compose down
```

## 📁 Project Structure

```
moogie/
├── 🎨 ui/                    # Svelte frontend
├── 🔗 api/                   # Go backend
├── ⚙️  config/               # Job configurations
├── 🏃 runner/                # Job execution service
├── 🐳 docker-compose.yaml   # Full stack definition
├── 🚀 start.sh              # Easy startup script
└── 📖 README.md             # This file
```

## 🔧 Development

### Frontend (Svelte)

```bash
cd ui
npm run dev    # Development server at http://localhost:3000
npm run build  # Production build
```

### Backend (Go API)

```bash
cd api
./dev.sh       # Development server with hot reload
make build     # Build binary
make test      # Run tests
```

### Database

```bash
# Migrations
cd api
make migrate-up DB_URL="postgres://moogie:moogie@localhost:5432/moogie?sslmode=disable"
make migrate-down

# Create new migration
make migrate-create NAME="add_new_feature"
```

## 🌐 API Endpoints

- `GET /api/v1/jobs` - List monitoring jobs
- `GET /api/v1/jobs/:id` - Get job details
- `POST /api/v1/executions` - Create execution result
- `GET /api/v1/dashboard/summary` - Dashboard metrics
- `GET /ws` - WebSocket for real-time updates
- `GET /health` - Health check

**API Documentation:** http://localhost:8080/swagger/index.html

## 📊 Features

### Dashboard

- 📈 **Real-time metrics** - Success rates, response times
- 📅 **Date range filtering** - Last 7, 30, 90 days or custom
- 🔄 **Live updates** - WebSocket-powered real-time data
- 📱 **Responsive design** - Works on mobile and desktop

### Job Management

- ⚙️ **Multiple check types** - API health, ping, DNS, SSL, TCP
- 📋 **Execution history** - Detailed logs with filtering
- 🎯 **Flexible configuration** - JSON-based job configs
- 📊 **Performance metrics** - Response times, success rates

### Monitoring Types

- 🌐 **API Health Checks** - HTTP endpoint monitoring
- 🏓 **Ping Connectivity** - Network reachability tests
- 🔍 **DNS Resolution** - Domain name resolution checks
- 🔒 **SSL Certificate** - Certificate expiry monitoring
- 🔌 **TCP Port** - Port connectivity tests

## 🗄️ Database Schema

### Jobs Table

- Configuration storage for all monitoring jobs
- JSON-based flexible config per job type
- Enable/disable functionality

### Executions Table

- Historical execution results
- Performance metrics (response time)
- Flexible JSON details per check type

## 🔄 Real-time Updates

WebSocket integration provides live updates for:

- New execution results
- Job status changes
- Dashboard metric updates

## 🐳 Docker Configuration

### Services

- **postgres**: PostgreSQL 15 with sample data
- **api**: Go backend with auto-migrations
- **ui**: Svelte frontend with optimized build
- **pgadmin**: Database admin interface (optional)

### Volumes

- `postgres_data`: Persistent database storage
- `pgadmin_data`: PgAdmin configuration

### Networks

- `moogie-network`: Internal service communication

## 🔒 Environment Variables

### API Configuration

```env
APP_ENV=development
APP_PORT=8080
DB_HOST=postgres
DB_USER=moogie
DB_PASSWORD=moogie
DB_NAME=moogie
ALLOWED_ORIGINS=http://localhost:3000
```

### UI Configuration

```env
MOOGIE_API_URL=http://localhost:8080
```

## 🧪 Sample Data

The system includes sample monitoring jobs:

- API health check (example.com)
- Database TCP connectivity
- DNS resolution test
- Ping connectivity check
- SSL certificate monitoring

Sample execution history is generated for the last 7 days with realistic success/failure patterns.

## 🛠️ Troubleshooting

### Common Issues

**Port conflicts:**

```bash
# Check what's using the ports
lsof -i :3000,8080,5432,5050

# Or use different ports
docker-compose up --scale ui=0  # Skip UI if port 3000 is busy
```

**Database connection issues:**

```bash
# Check database health
docker-compose logs postgres

# Reset database
docker-compose down -v  # Removes volumes
docker-compose up
```

**Build failures:**

```bash
# Clean rebuild
docker-compose down
docker-compose build --no-cache
docker-compose up
```

### Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api
docker-compose logs -f ui
docker-compose logs -f postgres
```

## 🚦 Health Checks

All services include health checks:

- **API**: `GET /health`
- **Database**: `pg_isready`
- **UI**: HTTP response check

Monitor with:

```bash
docker-compose ps
```

## 📈 Production Deployment

For production:

1. **Update environment variables**
2. **Configure proper secrets**
3. **Set up reverse proxy** (nginx/traefik)
4. **Configure SSL certificates**
5. **Set up backups** for PostgreSQL
6. **Configure monitoring** and alerting

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

---

**Ready to monitor?** Run `./start.sh` and visit http://localhost:3000
