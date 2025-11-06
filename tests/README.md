# Moogie API Tests

This directory contains Bruno API test collections for the Moogie monitoring application.

## Prerequisites

1. **Install Bruno**: Download and install Bruno from [https://www.usebruno.com/](https://www.usebruno.com/)
   - macOS: `brew install bruno`
   - Or download the desktop app from the website

2. **Start the Moogie API**: Make sure your API server is running
   ```bash
   cd api
   make dev
   # or
   go run cmd/server/main.go
   ```

## Getting Started

1. **Open Bruno** and click "Open Collection"
2. **Navigate** to this directory (`tests/api`) and open it as a collection
3. **Select Environment**: Choose the appropriate environment from the dropdown:
   - `local` - for local development (default: localhost:8080)
   - `development` - for development server
   - `staging` - for staging environment
   - `production` - for production environment

## Test Organization

The tests are organized into logical groups:

### 📋 Health
- `health-check.bru` - Basic health endpoint test

### 👔 Jobs  
- `get-all-jobs.bru` - Get all jobs
- `get-jobs-with-date-range.bru` - Get jobs with date filtering
- `get-job-by-id.bru` - Get specific job by ID

### ⚡ Executions
- `create-execution-success.bru` - Create successful execution
- `create-execution-failure.bru` - Create failed execution  
- `create-execution-invalid.bru` - Test validation errors

### 📊 Dashboard
- `get-summary.bru` - Get dashboard summary metrics

## Running Tests

### Individual Tests
- Click on any test file and hit the "Send" button
- View the response and test results in the right panel

### Collection Runner
1. Right-click on the collection name in the sidebar
2. Select "Run Collection"
3. Choose which tests to run
4. View aggregated results

### Command Line (Optional)
Bruno also supports CLI execution:
```bash
# Install Bruno CLI
npm install -g @usebruno/cli

# Run all tests
bru run tests/api --env local

# Run specific folder
bru run tests/api/jobs --env local
```

## Environment Variables

The tests use these variables (defined in `environments.bru`):

- `base_url` - Base URL for the API server
- `api_base` - Base URL for API endpoints (includes /api/v1)

### Customizing for Your Setup

If your API runs on a different port or host, either:

1. **Modify the environment file**: Edit `environments.bru`
2. **Override in Bruno UI**: Use the environment panel in Bruno to temporarily override variables

## Test Data Requirements

Some tests expect certain data to exist:

### Jobs
- The `get-job-by-id.bru` test looks for job ID `1`
- Update the `job_id` variable if needed

### Executions
- Tests create executions for jobs named `test-api-health` and `test-ssl-check`
- These jobs should exist in your database, or update the `job_name` in test bodies

## Adding New Tests

1. **Create a new .bru file** in the appropriate folder
2. **Use the Bruno syntax**:
   ```
   meta {
     name: Test Name
     type: http
     seq: 1
   }

   get {
     url: {{api_base}}/endpoint
     body: none
     auth: none
   }

   tests {
     test("description", function() {
       expect(res.getStatus()).to.equal(200);
     });
   }
   ```

3. **Follow naming conventions**:
   - Use kebab-case for file names
   - Group related tests in folders
   - Use descriptive names

## Test Patterns

### Common Assertions
```javascript
// Status codes
expect(res.getStatus()).to.equal(200);
expect([200, 404]).to.include(res.getStatus());

// Response body
expect(res.getBody()).to.be.an('array');
expect(res.getBody()).to.have.property('id');

// Performance
expect(res.getResponseTime()).to.be.below(1000);

// Headers
expect(res.getHeader('content-type')).to.include('application/json');
```

### Environment-Specific Tests
```javascript
// Skip certain checks in production
if (bru.getEnvVar('base_url').includes('localhost')) {
  test("development only test", function() {
    // test code
  });
}
```

## Integration with CI/CD

To run these tests in your CI/CD pipeline:

```yaml
# Example GitHub Actions step
- name: Run API Tests
  run: |
    # Start the API server
    cd api && make dev &
    
    # Wait for server to be ready
    sleep 10
    
    # Run Bruno tests
    npx @usebruno/cli run tests/api --env local
    
    # Stop the server
    pkill -f "go run"
```

## Troubleshooting

### Connection Issues
- Verify API server is running on expected port
- Check firewall/network settings
- Confirm environment variables are correct

### Test Failures  
- Check API server logs for errors
- Verify test data exists in database
- Update test expectations if API responses changed

### Bruno Issues
- Ensure Bruno desktop app is up to date
- Clear Bruno cache if having issues
- Check Bruno logs in the app's console

## Contributing

When adding new endpoints to the API:

1. **Add corresponding tests** in the appropriate folder
2. **Update environment files** if new variables are needed
3. **Update this README** with any new setup requirements
4. **Test against all environments** before merging