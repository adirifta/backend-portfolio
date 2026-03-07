-- Migration: Create visitors table for tracking website visitors
CREATE TABLE IF NOT EXISTS visitors (
    id BIGSERIAL PRIMARY KEY,
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT,
    path VARCHAR(255) NOT NULL DEFAULT '/',
    referer VARCHAR(512),
    country VARCHAR(100),
    city VARCHAR(100),
    visited_at TIMESTAMP
    WITH
        TIME ZONE NOT NULL DEFAULT NOW(),
        created_at TIMESTAMP
    WITH
        TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes for fast aggregation queries
CREATE INDEX IF NOT EXISTS idx_visitors_visited_at ON visitors (visited_at);

CREATE INDEX IF NOT EXISTS idx_visitors_ip_address ON visitors (ip_address);

CREATE INDEX IF NOT EXISTS idx_visitors_path ON visitors (path);

CREATE INDEX IF NOT EXISTS idx_visitors_ip_visited ON visitors (ip_address, visited_at);