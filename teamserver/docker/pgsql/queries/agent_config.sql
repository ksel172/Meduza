CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.{TABLE_NAME} (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    sleep INTEGER NOT NULL,
    jitter INTEGER NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    kill_date TIMESTAMPTZ NOT NULL,
    working_hours_start INTEGER NOT NULL,
    working_hours_end INTEGER NOT NULL
);
