CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.{TABLE_NAME} (
    id UUID PRIMARY KEY,
    payload_id UUID NOT NULL REFERENCES {POSTGRES_SCHEMA}.payloads(payload_id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    architecture VARCHAR(50) NOT NULL,
    parameters JSONB NOT NULL,
    start_time TIMESTAMPTZ,
    end_time TIMESTAMPTZ,
    output_path TEXT,
    error_message TEXT,
    build_log TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
);