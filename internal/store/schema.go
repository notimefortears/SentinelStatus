package store

const Schema = `
CREATE TABLE IF NOT EXISTS targets (
    id SERIAL PRIMARY KEY,
    url TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS monitor_results (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    status_code INTEGER,
    latency_ms INTEGER,
    checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`