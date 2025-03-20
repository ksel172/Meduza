CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.{TABLE_NAME} (
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   controller_id UUID,
   type VARCHAR(255) NOT NULL,
   name VARCHAR(255) NOT NULL,
   status INTEGER NOT NULL,
   description TEXT,
   config JSONB,
   logging_enabled BOOLEAN NOT NULL DEFAULT FALSE,
   logging JSONB,
   created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
   started_at TIMESTAMPTZ,
   stopped_at TIMESTAMPTZ,
   deleted_at TIMESTAMPTZ,

   FOREIGN KEY (controller_id) REFERENCES {POSTGRES_SCHEMA}.controllers(id) ON DELETE SET NULL
);