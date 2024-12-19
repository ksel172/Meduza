CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.agent_config (
    agent_id UUID PRIMARY KEY,
    callback_urls TEXT[] NOT NULL,
    rotation_type VARCHAR(50) NOT NULL,
    rotation_retries INTEGER NOT NULL,
    sleep INTEGER NOT NULL,
    jitter INTEGER NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    kill_date TIMESTAMPTZ NOT NULL,
    working_hours INTEGER[2] NOT NULL,
    headers JSONB NOT NULL,
    FOREIGN KEY (agent_id) REFERENCES {POSTGRES_SCHEMA}.agents(id) ON DELETE CASCADE
);