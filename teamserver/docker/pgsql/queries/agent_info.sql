CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.agent_info (
    agent_id UUID PRIMARY KEY,
    host_name VARCHAR(255) NOT NULL,
    ip_address VARCHAR(50) NOT NULL, 
    user_name VARCHAR(255) NOT NULL,
    system_info TEXT NOT NULL,
    os_info TEXT NOT NULL,
    FOREIGN KEY (agent_id) REFERENCES {POSTGRES_SCHEMA}.agents(id) ON DELETE CASCADE
);