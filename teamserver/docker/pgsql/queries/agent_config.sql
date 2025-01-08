CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.agent_config (
    config_id UUID PRIMARY KEY,
    listener_id UUID NOT NULL REFERENCES {POSTGRES_SCHEMA}.listeners(id),
    sleep INTEGER NOT NULL,
    jitter INTEGER NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    kill_date TIMESTAMPTZ NOT NULL,
    working_hours_start INTEGER NOT NULL,
    working_hours_end INTEGER NOT NULL
);
