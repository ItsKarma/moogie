-- +goose Up
-- Create jobs table
CREATE TABLE jobs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    type VARCHAR(100) NOT NULL,
    config JSONB NOT NULL,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create index on job name for faster lookups
CREATE INDEX idx_jobs_name ON jobs(name);
CREATE INDEX idx_jobs_type ON jobs(type);
CREATE INDEX idx_jobs_enabled ON jobs(enabled);

-- +goose Down
DROP TABLE IF EXISTS jobs;
