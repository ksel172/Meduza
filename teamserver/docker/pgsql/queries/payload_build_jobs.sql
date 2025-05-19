CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.{TABLE_NAME} (
    job_id UUID PRIMARY KEY,
    payload_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL,
    architecture VARCHAR(50) NOT NULL,
    parameters JSONB,
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    output_path TEXT,
    error_message TEXT,
    build_log TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    FOREIGN KEY (payload_id) REFERENCES {POSTGRES_SCHEMA}.payloads(payload_id) ON DELETE CASCADE
);