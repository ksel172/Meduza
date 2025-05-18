CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.{TABLE_NAME} (
    job_id UUID PRIMARY KEY,
    payload_id UUID NOT NULL REFERENCES {POSTGRES_SCHEMA}.payload_manifests(manifest_id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    architecture VARCHAR(50) NOT NULL,
    parameters JSONB NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    output_path TEXT,
    error_message TEXT,
    build_log TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);