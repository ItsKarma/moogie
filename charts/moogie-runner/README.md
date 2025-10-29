# Moogie Runner Helm Chart

This Helm chart deploys Moogie synthetic check runners as Kubernetes CronJobs.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+
- Moogie API server deployed and accessible

## Installation

### Add the Helm repository (if published)

```bash
helm repo add moogie https://charts.moogie.io
helm repo update
```

### Install from local chart

```bash
# From the charts directory
helm install moogie-runner ./moogie-runner \
  --namespace moogie \
  --create-namespace \
  --set global.apiUrl=http://moogie-api.moogie.svc.cluster.local:8080
```

### Install with custom values

```bash
helm install moogie-runner ./moogie-runner \
  --namespace moogie \
  --create-namespace \
  --values my-values.yaml
```

## Configuration

### Global Settings

| Parameter                 | Description                      | Default                                           |
| ------------------------- | -------------------------------- | ------------------------------------------------- |
| `global.apiUrl`           | Moogie API server URL            | `http://moogie-api.moogie.svc.cluster.local:8080` |
| `global.image.repository` | Runner image repository          | `moogie-runner`                                   |
| `global.image.tag`        | Runner image tag                 | `latest`                                          |
| `global.image.pullPolicy` | Image pull policy                | `IfNotPresent`                                    |
| `global.imagePullSecrets` | Image pull secrets               | `[]`                                              |
| `global.resources`        | Default resource limits/requests | See `values.yaml`                                 |

### CronJob Settings

| Parameter                            | Description                       | Default     |
| ------------------------------------ | --------------------------------- | ----------- |
| `cronJob.successfulJobsHistoryLimit` | Number of successful jobs to keep | `3`         |
| `cronJob.failedJobsHistoryLimit`     | Number of failed jobs to keep     | `1`         |
| `cronJob.activeDeadlineSeconds`      | Job execution deadline            | `300`       |
| `cronJob.restartPolicy`              | Pod restart policy                | `OnFailure` |
| `cronJob.concurrencyPolicy`          | Concurrency policy                | `Forbid`    |

### Check Configuration

Checks are defined in the `checks` array. Each check requires:

- `name`: Unique name for the check (used as job name in Moogie)
- `enabled`: Whether to deploy this check
- `schedule`: Cron schedule expression
- `type`: Check type (`http`, `ssl`, `dns`, `tcp`, or `custom`)
- `config`: Type-specific configuration

#### HTTP Check Example

```yaml
checks:
  - name: api-health-check
    enabled: true
    schedule: "*/5 * * * *"
    type: http
    config:
      url: "https://api.example.com/health"
      method: "GET"
      expectedStatus: 200
      timeout: 30
      headers:
        Authorization: "Bearer token"
        Accept: "application/json"
```

#### SSL Check Example

```yaml
checks:
  - name: ssl-certificate
    enabled: true
    schedule: "0 8 * * *"
    type: ssl
    config:
      host: "www.example.com"
      port: 443
      daysWarning: 30
      timeout: 15
```

#### DNS Check Example

```yaml
checks:
  - name: dns-check
    enabled: true
    schedule: "*/15 * * * *"
    type: dns
    config:
      hostname: "www.example.com"
      expectedIps: "93.184.216.34,93.184.216.35"
      timeout: 10
      dnsServer: "8.8.8.8:53"
```

#### TCP Check Example

```yaml
checks:
  - name: database-check
    enabled: true
    schedule: "*/10 * * * *"
    type: tcp
    config:
      host: "db.example.com"
      port: 5432
      timeout: 10
```

#### Custom Container Example

```yaml
checks:
  - name: puppeteer-check
    enabled: true
    schedule: "*/30 * * * *"
    type: custom
    image:
      repository: your-registry/puppeteer-check
      tag: v1.0.0
    env:
      - name: TARGET_URL
        value: "https://app.example.com"
      - name: TIMEOUT
        value: "60000"
    resources:
      limits:
        cpu: 500m
        memory: 512Mi
```

## Managing Checks

### View Deployed CronJobs

```bash
kubectl get cronjobs -n moogie -l app.kubernetes.io/name=moogie-runner
```

### View Recent Job Runs

```bash
kubectl get jobs -n moogie -l app.kubernetes.io/name=moogie-runner
```

### View Logs

```bash
# View logs for a specific check
kubectl logs -n moogie -l moogie.io/check-name=api-health-check --tail=50

# Follow logs
kubectl logs -n moogie -l moogie.io/check-name=api-health-check -f
```

### Manually Trigger a Check

```bash
kubectl create job --from=cronjob/moogie-runner-api-health-check \
  -n moogie \
  api-health-check-manual-$(date +%s)
```

### Suspend/Resume a Check

```bash
# Suspend
kubectl patch cronjob moogie-runner-api-health-check \
  -n moogie \
  -p '{"spec":{"suspend":true}}'

# Resume
kubectl patch cronjob moogie-runner-api-health-check \
  -n moogie \
  -p '{"spec":{"suspend":false}}'
```

## Updating Checks

### Update a single check

Edit your `values.yaml` and upgrade:

```bash
helm upgrade moogie-runner ./moogie-runner \
  --namespace moogie \
  --values values.yaml
```

### Add a new check

Add the check to `values.yaml` and upgrade:

```yaml
checks:
  # ... existing checks ...
  - name: new-check
    enabled: true
    schedule: "*/5 * * * *"
    type: http
    config:
      url: "https://new-endpoint.com"
```

```bash
helm upgrade moogie-runner ./moogie-runner \
  --namespace moogie \
  --values values.yaml
```

### Remove a check

Set `enabled: false` or remove the check from `values.yaml`:

```bash
helm upgrade moogie-runner ./moogie-runner \
  --namespace moogie \
  --values values.yaml
```

## Uninstalling

```bash
helm uninstall moogie-runner --namespace moogie
```

## Troubleshooting

### Check CronJob status

```bash
kubectl describe cronjob moogie-runner-api-health-check -n moogie
```

### View failed job pods

```bash
kubectl get pods -n moogie -l app.kubernetes.io/name=moogie-runner --field-selector=status.phase=Failed
```

### Debug a failed job

```bash
# Get the pod name
kubectl get pods -n moogie -l moogie.io/check-name=api-health-check

# View logs
kubectl logs -n moogie <pod-name>

# Describe pod for events
kubectl describe pod -n moogie <pod-name>
```

### Common Issues

1. **Jobs not running**: Check CronJob schedule and ensure it's not suspended
2. **API connection failures**: Verify `global.apiUrl` is correct and API is accessible
3. **Image pull errors**: Ensure image exists and imagePullSecrets are configured
4. **Job name not found errors**: Ensure job names match between Helm chart and Moogie API

## Best Practices

1. **Use meaningful job names**: Match them with your Moogie job names
2. **Set appropriate schedules**: Balance monitoring needs with cluster resources
3. **Configure resource limits**: Prevent runaway jobs from consuming cluster resources
4. **Monitor job history**: Keep `successfulJobsHistoryLimit` and `failedJobsHistoryLimit` reasonable
5. **Use custom images for complex checks**: Don't try to fit everything into built-in check types
6. **Set activeDeadlineSeconds**: Ensure stuck jobs don't run indefinitely
7. **Use concurrencyPolicy: Forbid**: Prevent overlapping job executions

## Examples

See the `examples/` directory for complete example configurations.
