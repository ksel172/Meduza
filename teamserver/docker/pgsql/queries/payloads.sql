CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.{TABLE_NAME} (
    payload_id UUID PRIMARY KEY,
    listener_id UUID,
    config_id UUID,
    manifest_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    architecture VARCHAR(50) NOT NULL,
    public_key BYTEA,
    private_key BYTEA,
    token VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (manifest_id) REFERENCES {POSTGRES_SCHEMA}.payload_manifests(manifest_id) ON DELETE CASCADE
);