-- +goose Up
-- Insert sample jobs based on your config directory structure
INSERT INTO jobs (name, type, config, enabled) VALUES 
(
    'api-health-check',
    'api-health',
    '{"url": "https://api.example.com/health", "method": "GET", "expected_status": 200, "timeout": 30}',
    true
),
(
    'database-tcp-check',
    'tcp',
    '{"host": "db.example.com", "port": 5432, "timeout": 10}',
    true
),
(
    'dns-resolution-check',
    'dns',
    '{"hostname": "example.com", "record_type": "A", "expected_ip": "93.184.216.34"}',
    true
),
(
    'ping-connectivity-check',
    'ping',
    '{"host": "example.com", "count": 4, "timeout": 5}',
    true
),
(
    'ssl-certificate-check',
    'ssl',
    '{"hostname": "example.com", "port": 443, "warning_days": 30, "critical_days": 7}',
    true
);

-- Insert some sample execution history for the last 7 days
-- This creates realistic test data
INSERT INTO executions (job_id, status, response_time, details, timestamp) 
SELECT 
    j.id,
    CASE 
        WHEN random() < 0.85 THEN 'success'
        ELSE 'failure'
    END as status,
    (random() * 1000 + 50)::int as response_time,
    CASE j.type
        WHEN 'api-health' THEN '{"status_code": 200, "body_size": 234}'
        WHEN 'tcp' THEN '{"connection_time": 45, "port_open": true}'
        WHEN 'dns' THEN '{"resolved_ip": "93.184.216.34", "query_time": 23}'
        WHEN 'ping' THEN '{"packets_sent": 4, "packets_received": 4, "packet_loss": 0, "avg_rtt": 25.4}'
        WHEN 'ssl' THEN '{"expires_in_days": 92, "issuer": "Let''s Encrypt", "valid": true}'
    END::jsonb as details,
    NOW() - (random() * interval '7 days') as timestamp
FROM jobs j
CROSS JOIN generate_series(1, 50) -- Generate 50 executions per job
ORDER BY random();

-- +goose Down
DELETE FROM executions;
DELETE FROM jobs;
