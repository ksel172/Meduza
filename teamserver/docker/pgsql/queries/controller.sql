CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.{TABLE_NAME} (
    id UUID PRIMARY KEY,
    endpoint VARCHAR(255) NOT NULL,
    public_key BYTEA NOT NULL,
    private_key BYTEA NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    heartbeat_timestamp TIMESTAMPTZ,
);