-- +goose Up
-- Create executions table
CREATE TABLE executions (
    id SERIAL PRIMARY KEY,
    job_id INTEGER NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL CHECK (status IN ('success', 'failure')),
    response_time INTEGER DEFAULT 0,
    details JSONB,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX idx_executions_job_id ON executions(job_id);
CREATE INDEX idx_executions_timestamp ON executions(timestamp);
CREATE INDEX idx_executions_status ON executions(status);
CREATE INDEX idx_executions_job_timestamp ON executions(job_id, timestamp);

-- +goose Down
DROP TABLE IF EXISTS executions;
