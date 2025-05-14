CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.{TABLE_NAME}(
   id UUID PRIMARY KEY NOT NULL,
   config_id UUID NOT NULL REFERENCES {POSTGRES_SCHEMA}.agent_config(id),
   name VARCHAR(255) NOT NULL,
   note TEXT NOT NULL,
   status VARCHAR(50) NOT NULL,
   first_callback TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
   last_callback TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
   modified_at TIMESTAMPTZ DEFAULT NULL
);

-- Canonical Agents Table
CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.agents (
   id BIGSERIAL PRIMARY KEY,
   agent_id UUID NOT NULL, -- assigned by the agent itself, not unique
   first_seen TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
   last_seen TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
   current_config_id UUID,
   status VARCHAR(50) NOT NULL,
   FOREIGN KEY (current_config_id) REFERENCES {POSTGRES_SCHEMA}.agent_configs(config_id)
);

-- History of Agent Details (IP, hostname, etc.)
CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.agent_info_history (
   id BIGSERIAL PRIMARY KEY,
   agent_fk BIGINT NOT NULL REFERENCES {POSTGRES_SCHEMA}.agents(id) ON DELETE CASCADE,
   host_name VARCHAR(255) NOT NULL,
   ip_address VARCHAR(50) NOT NULL, 
   user_name VARCHAR(255) NOT NULL,
   system_info TEXT NOT NULL,
   os_info TEXT NOT NULL,
   recorded_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Configurations
CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.agent_configs (
   config_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   listener_id UUID NOT NULL REFERENCES {POSTGRES_SCHEMA}.listeners(id),
   sleep INTEGER NOT NULL,
   jitter INTEGER NOT NULL,
   start_date TIMESTAMPTZ NOT NULL,
   kill_date TIMESTAMPTZ NOT NULL,
   working_hours_start INTEGER NOT NULL,
   working_hours_end INTEGER NOT NULL
);

-- Track config changes over time (optional, if you want history of config usage)
CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.agent_config_history (
    id BIGSERIAL PRIMARY KEY,
    agent_fk BIGINT NOT NULL REFERENCES {POSTGRES_SCHEMA}.agents(id) ON DELETE CASCADE,
    config_id UUID NOT NULL REFERENCES {POSTGRES_SCHEMA}.agent_configs(config_id),
    applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
