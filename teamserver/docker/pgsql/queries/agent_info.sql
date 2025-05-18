CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.{TABLE_NAME} (
    agent_id UUID PRIMARY KEY REFERENCES {POSTGRES_SCHEMA}.agents(id) ON DELETE CASCADE,
    hostname VARCHAR(255),
    ip_address VARCHAR(50), 
    user_name VARCHAR(255),
    system_info TEXT,
    os_info TEXT
);