CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.{TABLE_NAME}(
    payload_id UUID PRIMARY KEY NOT NULL,
    payload_name VARCHAR(255) NOT NULL,
    config_id UUID NOT NULL,
    listener_id UUID NOT NULL,
    public_key BYTEA NOT NULL,
    private_key BYTEA NOT NULL,
    payload_token UUID NOT NULL,
    arch VARCHAR(50) NOT NULL,
    listener_config JSONB NOT NULL,
    sleep INTEGER NOT NULL,
    jitter INTEGER NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    kill_date TIMESTAMPTZ NOT NULL,
    working_hours_start INTEGER NOT NULL,
    working_hours_end INTEGER NOT NULL
);