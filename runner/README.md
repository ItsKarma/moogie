# Moogie Runner

The Moogie Runner is a lightweight Go-based container that executes synthetic checks and reports results to the Moogie API server. It's designed to run as Kubernetes CronJobs for scheduled monitoring.

## Built-in Check Types

### HTTP Check

Performs HTTP/HTTPS requests and validates responses.

**Environment Variables:**

- `CHECK_TYPE=http` (required)
- `HTTP_URL` - Target URL (required)
- `HTTP_METHOD` - HTTP method (default: GET)
- `HTTP_TIMEOUT` - Timeout in seconds (default: 30)
- `HTTP_EXPECTED_STATUS` - Expected status code (default: 200)
- `HTTP_HEADERS` - Comma-separated headers (e.g., `Authorization:Bearer token,Accept:application/json`)
- `HTTP_BODY` - Request body for POST/PUT

**Example:**

```yaml
env:
  - name: CHECK_TYPE
    value: "http"
  - name: HTTP_URL
    value: "https://api.example.com/health"
  - name: HTTP_METHOD
    value: "GET"
  - name: HTTP_EXPECTED_STATUS
    value: "200"
  - name: HTTP_TIMEOUT
    value: "30"
```

### SSL Certificate Check

Validates SSL/TLS certificates and checks expiry.

**Environment Variables:**

- `CHECK_TYPE=ssl` (required)
- `SSL_HOST` - Target host (required)
- `SSL_PORT` - Target port (default: 443)
- `SSL_TIMEOUT` - Timeout in seconds (default: 15)
- `SSL_DAYS_WARNING` - Days before expiry to warn (default: 30)

**Example:**

```yaml
env:
  - name: CHECK_TYPE
    value: "ssl"
  - name: SSL_HOST
    value: "www.example.com"
  - name: SSL_PORT
    value: "443"
  - name: SSL_DAYS_WARNING
    value: "30"
```

### DNS Check

Verifies DNS resolution and validates IPs.

**Environment Variables:**

- `CHECK_TYPE=dns` (required)
- `DNS_HOSTNAME` - Hostname to resolve (required)
- `DNS_EXPECTED_IPS` - Comma-separated expected IPs (optional)
- `DNS_TIMEOUT` - Timeout in seconds (default: 10)
- `DNS_SERVER` - Custom DNS server (optional, e.g., `8.8.8.8:53`)

**Example:**

```yaml
env:
  - name: CHECK_TYPE
    value: "dns"
  - name: DNS_HOSTNAME
    value: "www.example.com"
  - name: DNS_EXPECTED_IPS
    value: "93.184.216.34"
  - name: DNS_TIMEOUT
    value: "10"
```

### TCP Check

Tests TCP connectivity to a host and port.

**Environment Variables:**

- `CHECK_TYPE=tcp` (required)
- `TCP_HOST` - Target host (required)
- `TCP_PORT` - Target port (required)
- `TCP_TIMEOUT` - Timeout in seconds (default: 10)

**Example:**

```yaml
env:
  - name: CHECK_TYPE
    value: "tcp"
  - name: TCP_HOST
    value: "db.example.com"
  - name: TCP_PORT
    value: "5432"
  - name: TCP_TIMEOUT
    value: "10"
```

## Required Environment Variables (All Checks)

- `MOOGIE_API_URL` - Moogie API server URL (e.g., `http://moogie-api:8080`)
- `JOB_NAME` - Job name from Moogie (used to associate execution results with the correct job)

## Building

```bash
docker build -t moogie-runner:latest .
```

## Running Locally

```bash
# HTTP check example
docker run --rm \
  --network moogie_moogie-network \
  -e CHECK_TYPE=http \
  -e HTTP_URL=https://httpbin.org/status/200 \
  -e HTTP_METHOD=GET \
  -e HTTP_EXPECTED_STATUS=200 \
  -e MOOGIE_API_URL=http://moogie-api:8080 \
  -e JOB_NAME=api-health-check-production \
  moogie-runner:latest

# SSL check example
docker run --rm \
  --network moogie_moogie-network \
  -e CHECK_TYPE=ssl \
  -e SSL_HOST=www.google.com \
  -e SSL_PORT=443 \
  -e SSL_DAYS_WARNING=30 \
  -e MOOGIE_API_URL=http://moogie-api:8080 \
  -e JOB_NAME=ssl-certificate-check \
  moogie-runner:latest

# DNS check example
docker run --rm \
  --network moogie_moogie-network \
  -e CHECK_TYPE=dns \
  -e DNS_HOSTNAME=www.google.com \
  -e MOOGIE_API_URL=http://moogie-api:8080 \
  -e JOB_NAME=dns-resolution-check \
  moogie-runner:latest

# TCP check example
docker run --rm \
  --network moogie_moogie-network \
  -e CHECK_TYPE=tcp \
  -e TCP_HOST=google.com \
  -e TCP_PORT=443 \
  -e MOOGIE_API_URL=http://moogie-api:8080 \
  -e JOB_NAME=tcp-connectivity-check \
  moogie-runner:latest
```

## Kubernetes CronJob Example

````yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: api-health-check
  namespace: moogie
spec:
  schedule: "*/5 * * * *" # Every 5 minutes
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: runner
              image: moogie-runner:latest
              env:
                - name: CHECK_TYPE
                  value: "http"
                - name: HTTP_URL
                  value: "https://api.example.com/health"
                - name: HTTP_METHOD
                  value: "GET"
                - name: HTTP_EXPECTED_STATUS
                  value: "200"
            - name: MOOGIE_API_URL
              value: "http://moogie-api.moogie.svc.cluster.local:8080"
            - name: JOB_NAME
              value: "api-health-check-production"
          restartPolicy: OnFailure
```## Custom Check Containers

You can create custom check containers for specialized monitoring (e.g., Puppeteer for browser automation). Your container must:

1. Perform the check/test
2. POST results to the Moogie API

### Custom Container Requirements

Your container should POST an execution result to:

````

POST {MOOGIE_API_URL}/api/v1/executions

````

**Payload:**

```json
**Payload:**
```json
{
  "job_name": "your-job-name",
  "status": "success",  // or "failure"
  "response_time": 1234,  // milliseconds
  "timestamp": "2025-10-27T12:00:00Z",
  "details": {
    "custom_field": "value",
    "error": "optional error message if status is failure"
  }
}
````

````

### Puppeteer Example

```dockerfile
# examples/puppeteer-check/Dockerfile
FROM node:18-slim

# Install Puppeteer dependencies
RUN apt-get update && apt-get install -y \
    chromium \
    && rm -rf /var/lib/apt/lists/*

ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true
ENV PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY check.js ./

CMD ["node", "check.js"]
````

```javascript
// examples/puppeteer-check/check.js
const puppeteer = require("puppeteer");
const axios = require("axios");

async function runCheck() {
  const startTime = Date.now();
  let status = "success";
  let errorMessage = "";

  try {
    const browser = await puppeteer.launch({
      headless: "new",
      args: ["--no-sandbox", "--disable-setuid-sandbox"],
    });

    const page = await browser.newPage();
    await page.goto(process.env.TARGET_URL, { waitUntil: "networkidle2" });

    // Perform your custom checks
    const title = await page.title();
    console.log("Page title:", title);

    // Check for specific elements, test functionality, etc.
    const hasLoginButton = (await page.$(".login-button")) !== null;
    if (!hasLoginButton) {
      throw new Error("Login button not found");
    }

    await browser.close();
  } catch (error) {
    status = "error";
    errorMessage = error.message;
    console.error("Check failed:", error);
  }

  const responseTime = Date.now() - startTime;


  // Report to Moogie API
  await axios.post(`${process.env.MOOGIE_API_URL}/api/v1/executions`, {
    job_name: process.env.JOB_NAME,
    status: status,
    response_time: responseTime,
    timestamp: new Date().toISOString(),
    details: {
      check_type: 'puppeteer',
      url: process.env.TARGET_URL,
      error: errorMessage
    }
  });
}
}

runCheck().catch(console.error);
```

### Custom Container CronJob

````yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: custom-puppeteer-check
  namespace: moogie
spec:
  schedule: "*/15 * * * *" # Every 15 minutes
  jobTemplate:
    spec:
      template:
        spec:
          containers:
            - name: puppeteer-check
              image: your-registry/puppeteer-check:latest
              env:
                - name: TARGET_URL
                  value: "https://app.example.com"
            - name: MOOGIE_API_URL
              value: "http://moogie-api.moogie.svc.cluster.local:8080"
            - name: JOB_NAME
              value: "custom-puppeteer-check"
          restartPolicy: OnFailure
```## Development

### Running Tests

```bash
go test ./...
````

### Building Locally

```bash
go build -o runner .
```

### Running Without Docker

```bash
export CHECK_TYPE=http
export HTTP_URL=https://httpbin.org/status/200
export MOOGIE_API_URL=http://moogie-api:8080
export JOB_NAME=api-health-check-production
./runner
```

## Architecture

The runner is designed to be:

- **Stateless**: No persistent storage required
- **Fast**: Small image (~10MB), quick startup
- **Configurable**: All settings via environment variables
- **Extensible**: Support for custom check containers
- **Kubernetes-native**: Works seamlessly with CronJobs

## Next Steps

See the main Moogie documentation for:

- Helm chart deployment
- Creating jobs via API
- Managing check configurations
- Viewing results in the dashboard
