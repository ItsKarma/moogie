// Mock data for jobs - will be replaced with API calls later

export const mockJobs = [
  {
    id: "api-health-check",
    status: "running",
    lastRun: "2024-10-20 14:30:00",
    executions: 1247,
    successRate: 98.2,
    config: {
      apiVersion: "moogie.io/v1",
      kind: "HttpCheck",
      metadata: {
        id: "api-health-check",
        displayName: "API Health Check",
        description:
          "Monitors the health endpoint of our main API service to ensure availability and proper response codes.",
        labels: {
          environment: "production",
          service: "api",
          team: "backend",
        },
      },
      spec: {
        url: "https://api.example.com/health",
        method: "GET",
        timeout: "30s",
        expectedStatusCode: 200,
        schedule: "*/5 * * * *",
        retries: 3,
        headers: {
          "User-Agent": "Moogie/1.0",
          Accept: "application/json",
        },
        alerts: {
          onFailure: true,
          email: "alerts@example.com",
        },
      },
    },
  },
  {
    id: "database-connection-check",
    status: "success",
    lastRun: "2024-10-20 14:15:00",
    executions: 892,
    successRate: 99.1,
    config: {
      apiVersion: "moogie.io/v1",
      kind: "TcpCheck",
      metadata: {
        id: "database-connection-check",
        displayName: "Database Connection",
        description:
          "Verifies TCP connectivity to our primary PostgreSQL database server.",
        labels: {
          environment: "production",
          service: "database",
          team: "infrastructure",
        },
      },
      spec: {
        host: "db.example.com",
        port: 5432,
        timeout: "10s",
        schedule: "*/2 * * * *",
        retries: 2,
        alerts: {
          onFailure: true,
          email: "dba@example.com",
        },
      },
    },
  },
  {
    id: "domain-dns-resolution",
    status: "failed",
    lastRun: "2024-10-20 13:30:00",
    executions: 456,
    successRate: 95.8,
    config: {
      apiVersion: "moogie.io/v1",
      kind: "DnsCheck",
      metadata: {
        id: "domain-dns-resolution",
        displayName: "DNS Resolution Check",
        description:
          "Ensures our domain resolves correctly to the expected IP address for users worldwide.",
        labels: {
          environment: "production",
          service: "dns",
          team: "infrastructure",
        },
      },
      spec: {
        domain: "example.com",
        recordType: "A",
        expectedIp: "192.168.1.1",
        timeout: "15s",
        schedule: "0 */6 * * *",
        retries: 3,
        alerts: {
          onFailure: true,
          email: "network@example.com",
        },
      },
    },
  },
  {
    id: "ssl-certificate-check",
    status: "warning",
    lastRun: "2024-10-20 12:00:00",
    executions: 234,
    successRate: 94.2,
    config: {
      apiVersion: "moogie.io/v1",
      kind: "SslCheck",
      metadata: {
        id: "ssl-certificate-check",
        displayName: "SSL Certificate Validation",
        description:
          "Validates SSL certificate expiration and ensures secure connections are properly configured.",
        labels: {
          environment: "production",
          service: "security",
          team: "infrastructure",
        },
      },
      spec: {
        host: "example.com",
        port: 443,
        timeout: "20s",
        schedule: "0 0 * * *",
        retries: 1,
        warningDays: 30,
        alerts: {
          onFailure: true,
          onWarning: true,
          email: "security@example.com",
        },
      },
    },
  },
  {
    id: "ping-connectivity-check",
    status: "success",
    lastRun: "2024-10-20 14:25:00",
    executions: 567,
    successRate: 99.8,
    config: {
      apiVersion: "moogie.io/v1",
      kind: "PingCheck",
      metadata: {
        id: "ping-connectivity-check",
        displayName: "Network Connectivity",
        description:
          "Tests basic network connectivity to external services to ensure internet access.",
        labels: {
          environment: "production",
          service: "network",
          team: "infrastructure",
        },
      },
      spec: {
        target: "8.8.8.8",
        timeout: "5s",
        schedule: "*/1 * * * *",
        retries: 2,
        alerts: {
          onFailure: true,
          email: "network@example.com",
        },
      },
    },
  },
];

export const mockDashboardJobs = [
  {
    id: 1,
    name: "Website Health Check",
    status: "running",
    lastRun: "2 minutes ago",
    nextRun: "in 28 minutes",
    successRate: 98.5,
  },
  {
    id: 2,
    name: "API Response Time",
    status: "success",
    lastRun: "15 minutes ago",
    nextRun: "in 15 minutes",
    successRate: 99.2,
  },
  {
    id: 3,
    name: "Database Connection",
    status: "failed",
    lastRun: "1 hour ago",
    nextRun: "in 2 hours",
    successRate: 87.3,
  },
  {
    id: 4,
    name: "SSL Certificate Check",
    status: "warning",
    lastRun: "3 hours ago",
    nextRun: "in 21 hours",
    successRate: 94.1,
  },
];

export const mockExecutionHistory = [
  {
    timestamp: "2024-10-20 14:30:00",
    status: "success",
    responseTime: 145,
    statusCode: 200,
    message: "API health check completed successfully",
    logs: [
      "[2024-10-20 14:30:00] Starting health check...",
      "[2024-10-20 14:30:00] Connecting to https://api.example.com/health",
      "[2024-10-20 14:30:00] Response received: 200 OK",
      "[2024-10-20 14:30:00] Response time: 145ms",
      "[2024-10-20 14:30:00] Health check passed all validations",
    ],
  },
  {
    timestamp: "2024-10-20 14:25:00",
    status: "success",
    responseTime: 132,
    statusCode: 200,
    message: "API health check completed successfully",
    logs: [
      "[2024-10-20 14:25:00] Starting health check...",
      "[2024-10-20 14:25:00] Connecting to https://api.example.com/health",
      "[2024-10-20 14:25:00] Response received: 200 OK",
      "[2024-10-20 14:25:00] Response time: 132ms",
      "[2024-10-20 14:25:00] Health check passed all validations",
    ],
  },
  {
    timestamp: "2024-10-20 14:20:00",
    status: "failed",
    responseTime: 5000,
    statusCode: 500,
    message: "Internal server error - health check failed",
    logs: [
      "[2024-10-20 14:20:00] Starting health check...",
      "[2024-10-20 14:20:00] Connecting to https://api.example.com/health",
      "[2024-10-20 14:20:00] Response received: 500 Internal Server Error",
      "[2024-10-20 14:20:00] Response time: 5000ms",
      "[2024-10-20 14:20:00] ERROR: Server returned error status",
      "[2024-10-20 14:20:00] Health check failed validation",
    ],
  },
  {
    timestamp: "2024-10-20 14:15:00",
    status: "warning",
    responseTime: 2500,
    statusCode: 200,
    message: "API responded but slower than expected",
    logs: [
      "[2024-10-20 14:15:00] Starting health check...",
      "[2024-10-20 14:15:00] Connecting to https://api.example.com/health",
      "[2024-10-20 14:15:00] Response received: 200 OK",
      "[2024-10-20 14:15:00] Response time: 2500ms",
      "[2024-10-20 14:15:00] WARNING: Response time exceeded threshold (2000ms)",
      "[2024-10-20 14:15:00] Health check completed with warnings",
    ],
  },
];
