CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.{TABLE_NAME} (
    "id" UUID PRIMARY KEY,
    "type" VARCHAR(10) NOT NULL CHECK ("type" IN ('cert', 'key')),
    "path" VARCHAR(255) NOT NULL,
    "filename" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL,
    "updated_at" TIMESTAMPTZ NOT NULL
);