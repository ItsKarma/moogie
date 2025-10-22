# Moogie

Synthetics Runner

## Overview

Moogie is a synthetics monitoring system that runs native in Kubernetes utilizing Cron Jobs to execute your synthetic checks. Instead of configuring checks through a UI, Moogie uses YAML configuration files that follow Kubernetes-style resource definitions.

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

## Deploying

This runs native in Kubernetes utilizing Cron Jobs to run your Synthetic Checks.
