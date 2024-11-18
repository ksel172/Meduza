CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE role_enum AS ENUM ('admin','moderator','client','visitor');
CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.{TABLE_NAME} (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    pw_hash TEXT NOT NULL,
    role role_enum NOT NULL DEFAULT 'visitor',
    created_by UUID DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ DEFAULT NULL
);
