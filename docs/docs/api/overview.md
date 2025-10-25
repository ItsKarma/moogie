---
sidebar_position: 1
---

# API Overview

The Moogie API provides a comprehensive REST interface for managing monitoring jobs, executions, and retrieving dashboard data. Built with Go and Gin, it offers high performance and reliability.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Currently, the API does not require authentication. This may be added in future versions for production deployments.

## Response Format

All API responses follow a consistent JSON format:

### Success Response

```json
{
  "data": { ... },
  "message": "Success"
}
```

### Error Response

```json
{
  "error": "Error message",
  "code": "ERROR_CODE"
}
```

## HTTP Status Codes

| Code | Description                        |
| ---- | ---------------------------------- |
| 200  | Success                            |
| 201  | Created                            |
| 400  | Bad Request - Invalid input        |
| 404  | Not Found - Resource doesn't exist |
| 500  | Internal Server Error              |

## Real-time Updates

The API supports WebSocket connections for real-time updates:

```
ws://localhost:8080/ws
```

When connected, you'll receive live updates for:

- New job executions
- Status changes
- Dashboard metrics updates

## Rate Limiting

Currently, no rate limiting is implemented. Consider implementing rate limiting for production deployments.

## CORS Policy

The API is configured to accept requests from:

- `http://localhost:3000` (UI Dashboard)
- `http://localhost:3000` (Documentation)

Additional origins can be configured via the `ALLOWED_ORIGINS` environment variable.

## Data Models

### Job Model

```json
{
  "id": 1,
  "name": "api-health-check",
  "type": "api-health",
  "config": {
    "url": "https://api.example.com/health",
    "method": "GET",
    "timeout": 30,
    "expected_status": 200
  },
  "enabled": true,
  "created_at": "2025-01-01T00:00:00Z",
  "updated_at": "2025-01-01T00:00:00Z",
  "success_rate": 98.5,
  "last_execution": "2025-01-15T10:30:00Z",
  "avg_response_time": 120.5,
  "executions": []
}
```

### Execution Model

```json
{
  "id": 123,
  "job_id": 1,
  "status": "success",
  "response_time": 120,
  "details": {
    "status_code": 200,
    "response_body": "OK",
    "headers": {
      "content-type": "application/json"
    }
  },
  "timestamp": "2025-01-15T10:30:00Z",
  "job": {
    "id": 1,
    "name": "api-health-check"
  }
}
```

## Check Types

Moogie supports multiple monitoring check types:

### API Health Check (`api-health`)

```json
{
  "type": "api-health",
  "config": {
    "url": "https://api.example.com/health",
    "method": "GET",
    "timeout": 30,
    "expected_status": 200,
    "headers": {
      "Authorization": "Bearer token"
    }
  }
}
```

### SSL Certificate Check (`ssl`)

```json
{
  "type": "ssl",
  "config": {
    "host": "example.com",
    "port": 443,
    "days_warning": 30
  }
}
```

### DNS Resolution Check (`dns`)

```json
{
  "type": "dns",
  "config": {
    "hostname": "example.com",
    "nameserver": "8.8.8.8",
    "record_type": "A"
  }
}
```

### Ping Connectivity Check (`ping`)

```json
{
  "type": "ping",
  "config": {
    "host": "example.com",
    "count": 4,
    "timeout": 10
  }
}
```

## Date Range Filtering

Many endpoints support date range filtering with query parameters:

| Parameter | Format     | Description            |
| --------- | ---------- | ---------------------- |
| `from`    | YYYY-MM-DD | Start date (inclusive) |
| `to`      | YYYY-MM-DD | End date (inclusive)   |

Example:

```
GET /api/v1/jobs?from=2025-01-01&to=2025-01-31
```

## Error Handling

The API provides detailed error messages to help with debugging:

### Validation Errors

```json
{
  "error": "Validation failed",
  "details": {
    "job_name": "Field is required",
    "status": "Must be 'success' or 'failure'"
  }
}
```

### Not Found Errors

```json
{
  "error": "Job not found",
  "job_id": 123
}
```

### Database Errors

```json
{
  "error": "Database connection failed",
  "code": "DB_CONNECTION_ERROR"
}
```

## Health Check

A simple health check endpoint is available:

```http
GET /health
```

Response:

```json
{
  "status": "healthy",
  "timestamp": "2025-01-15T10:30:00Z",
  "database": "connected",
  "version": "1.0.0"
}
```

## Next Steps

Explore the API in more detail:

- Check out the [Getting Started guide](../getting-started/docker-setup)
- View the main documentation at the [home page](/)

## SDK and Examples

### cURL Examples

All endpoint documentation includes cURL examples for easy testing.

### Future SDKs

We plan to provide official SDKs for:

- JavaScript/TypeScript
- Python
- Go

## Feedback

Found an issue with the API? Have a feature request?

- [Open an Issue](https://github.com/ItsKarma/moogie/issues)
- [Join Discussions](https://github.com/ItsKarma/moogie/discussions)
