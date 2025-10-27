-- Moogie Database Initialization
-- Creates tables and inserts sample data

-- Drop existing tables if they exist
DROP TABLE IF EXISTS executions CASCADE;
DROP TABLE IF EXISTS jobs CASCADE;

-- Create jobs table
CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(100) NOT NULL,
    config JSONB NOT NULL,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create executions table
CREATE TABLE executions (
    id SERIAL PRIMARY KEY,
    job_id INTEGER NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    response_time BIGINT,
    details JSONB,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_executions_job_id ON executions(job_id);
CREATE INDEX idx_executions_timestamp ON executions(timestamp);
CREATE INDEX idx_executions_job_timestamp ON executions(job_id, timestamp DESC);

-- Insert sample jobs with variety across different services and environments
INSERT INTO jobs (name, type, config, enabled) VALUES 
-- API Health Checks
(
    'api-health-check-production',
    'api-health',
    '{"url": "https://api.example.com/health", "method": "GET", "expected_status": 200, "timeout": 30, "metadata": {"labels": {"service": "api", "environment": "production", "team": "backend"}}}'::jsonb,
    true
),
(
    'api-health-check-staging',
    'api-health',
    '{"url": "https://staging-api.example.com/health", "method": "GET", "expected_status": 200, "timeout": 30, "metadata": {"labels": {"service": "api", "environment": "staging", "team": "backend"}}}'::jsonb,
    true
),
(
    'api-users-endpoint',
    'api-health',
    '{"url": "https://api.example.com/users", "method": "GET", "expected_status": 200, "timeout": 30, "metadata": {"labels": {"service": "api", "environment": "production", "team": "backend"}}}'::jsonb,
    true
),
(
    'payment-gateway-health',
    'api-health',
    '{"url": "https://payments.example.com/health", "method": "GET", "expected_status": 200, "timeout": 30, "metadata": {"labels": {"service": "payments", "environment": "production", "team": "backend"}}}'::jsonb,
    true
),
-- Database Checks
(
    'database-primary-check',
    'tcp',
    '{"host": "db-primary.example.com", "port": 5432, "timeout": 10, "metadata": {"labels": {"service": "database", "environment": "production", "team": "infrastructure"}}}'::jsonb,
    true
),
(
    'database-replica-check',
    'tcp',
    '{"host": "db-replica.example.com", "port": 5432, "timeout": 10, "metadata": {"labels": {"service": "database", "environment": "production", "team": "infrastructure"}}}'::jsonb,
    true
),
(
    'redis-cache-check',
    'tcp',
    '{"host": "redis.example.com", "port": 6379, "timeout": 10, "metadata": {"labels": {"service": "cache", "environment": "production", "team": "infrastructure"}}}'::jsonb,
    true
),
-- DNS Checks
(
    'dns-primary-domain',
    'dns',
    '{"hostname": "example.com", "record_type": "A", "expected_ip": "93.184.216.34", "metadata": {"labels": {"service": "dns", "environment": "production", "team": "infrastructure"}}}'::jsonb,
    true
),
(
    'dns-api-subdomain',
    'dns',
    '{"hostname": "api.example.com", "record_type": "A", "expected_ip": "93.184.216.35", "metadata": {"labels": {"service": "dns", "environment": "production", "team": "infrastructure"}}}'::jsonb,
    true
),
-- Ping/Network Checks
(
    'ping-production-server',
    'ping',
    '{"host": "server1.example.com", "count": 4, "timeout": 5, "metadata": {"labels": {"service": "network", "environment": "production", "team": "infrastructure"}}}'::jsonb,
    true
),
(
    'ping-load-balancer',
    'ping',
    '{"host": "lb.example.com", "count": 4, "timeout": 5, "metadata": {"labels": {"service": "network", "environment": "production", "team": "infrastructure"}}}'::jsonb,
    true
),
-- SSL Certificate Checks
(
    'ssl-main-domain',
    'ssl',
    '{"hostname": "example.com", "port": 443, "warning_days": 30, "critical_days": 7, "metadata": {"labels": {"service": "ssl", "environment": "production", "team": "security"}}}'::jsonb,
    true
),
(
    'ssl-api-domain',
    'ssl',
    '{"hostname": "api.example.com", "port": 443, "warning_days": 30, "critical_days": 7, "metadata": {"labels": {"service": "ssl", "environment": "production", "team": "security"}}}'::jsonb,
    true
),
(
    'ssl-payment-gateway',
    'ssl',
    '{"hostname": "payments.example.com", "port": 443, "warning_days": 30, "critical_days": 7, "metadata": {"labels": {"service": "ssl", "environment": "production", "team": "security"}}}'::jsonb,
    true
),
-- Some disabled checks for variety
(
    'api-health-check-development',
    'api-health',
    '{"url": "https://dev-api.example.com/health", "method": "GET", "expected_status": 200, "timeout": 30, "metadata": {"labels": {"service": "api", "environment": "development", "team": "backend"}}}'::jsonb,
    false
);

-- Insert execution history for the last 180 days
-- This creates realistic test data with varied patterns
INSERT INTO executions (job_id, status, response_time, details, timestamp) 
SELECT 
    j.id,
    -- Success rate varies by service type with some realistic failure patterns
    CASE 
        -- Most checks succeed most of the time (90% success)
        WHEN random() < 0.90 THEN 'success'
        -- Occasional failures
        WHEN random() < 0.95 THEN 'failure'
        -- Very rare errors
        ELSE 'error'
    END as status,
    -- Response time varies by check type
    CASE j.type
        WHEN 'api-health' THEN (random() * 800 + 100)::int  -- 100-900ms
        WHEN 'tcp' THEN (random() * 100 + 10)::int          -- 10-110ms
        WHEN 'dns' THEN (random() * 50 + 5)::int            -- 5-55ms
        WHEN 'ping' THEN (random() * 200 + 20)::int         -- 20-220ms
        WHEN 'ssl' THEN (random() * 300 + 50)::int          -- 50-350ms
        ELSE (random() * 500 + 50)::int
    END as response_time,
    -- Details vary by check type
    CASE j.type
        WHEN 'api-health' THEN jsonb_build_object(
            'status_code', CASE WHEN random() < 0.90 THEN 200 ELSE 500 END,
            'body_size', (random() * 1000 + 100)::int
        )
        WHEN 'tcp' THEN jsonb_build_object(
            'connection_time', (random() * 100)::int,
            'port_open', random() < 0.95
        )
        WHEN 'dns' THEN jsonb_build_object(
            'resolved_ip', '93.184.216.' || (30 + (random() * 10)::int)::text,
            'query_time', (random() * 50)::int
        )
        WHEN 'ping' THEN jsonb_build_object(
            'packets_sent', 4,
            'packets_received', CASE WHEN random() < 0.90 THEN 4 ELSE (random() * 4)::int END,
            'packet_loss', CASE WHEN random() < 0.90 THEN 0 ELSE (random() * 25)::int END,
            'avg_rtt', (random() * 100 + 10)
        )
        WHEN 'ssl' THEN jsonb_build_object(
            'expires_in_days', (random() * 300 + 30)::int,
            'issuer', CASE WHEN random() < 0.5 THEN 'Let''s Encrypt' ELSE 'DigiCert' END,
            'valid', random() < 0.95
        )
    END as details,
    -- Distribute executions over the last 180 days
    CURRENT_TIMESTAMP - (random() * interval '180 days') as timestamp
FROM jobs j
CROSS JOIN generate_series(1, 1000) -- Generate 1000 executions per job over 180 days
WHERE j.enabled = true  -- Only generate executions for enabled jobs
ORDER BY random();

-- Confirm initialization complete
SELECT 
    (SELECT COUNT(*) FROM jobs) as total_jobs,
    (SELECT COUNT(*) FROM jobs WHERE enabled = true) as enabled_jobs,
    (SELECT COUNT(*) FROM executions) as total_executions;
