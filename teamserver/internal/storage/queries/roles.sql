CREATE TABLE IF NOT EXISTS {DB_SCHEMA_NAME}.{TABLE_NAME} (
                                                             id SERIAL PRIMARY KEY,
                                                             role_name VARCHAR(50) UNIQUE NOT NULL,
    created_ts TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_ts TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
    );

INSERT INTO {POSTGRES_SCHEMA}.{TABLE_NAME} (role_name)
VALUES
    ('commander'), -- full access to manage and control the system
    ('operator'),  -- execute commands within defined limits
    ('auditor'),   -- read-only access for monitoring and compliance
    ('technician'),-- maintenance and support tasks
    ('observer')   -- limited view-only access