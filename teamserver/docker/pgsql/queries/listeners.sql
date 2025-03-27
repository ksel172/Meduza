CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS  {POSTGRES_SCHEMA}.{TABLE_NAME} (
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   type VARCHAR(255) NOT NULL,
   name VARCHAR(255) NOT NULL,
   status VARCHAR(255) NOT NULL,
   description TEXT,
   config JSONB,
   created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
   started_at TIMESTAMPTZ,
   stopped_at TIMESTAMPTZ
);
