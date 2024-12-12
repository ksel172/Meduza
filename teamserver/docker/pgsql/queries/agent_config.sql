CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.agent_config(
   id BIGSERIAL PRIMARY KEY,
   agent_id UUID NOT NULL,
   callback_urls TEXT[] NOT NULL,
   rotation_type VARCHAR(255) NOT NULL,
   rotation_retries INTEGER NOT NULL,
   sleep INTEGER NOT NULL,
   jitter INTEGER NOT NULL,
   start_date TIMESTAMPTZ NOT NULL,
   kill_date TIMESTAMPTZ NOT NULL,
   working_hours INTEGER[] NOT NULL,
   headers TEXT[] NOT NULL,

   CONSTRAINT fk_agent FOREIGN KEY(agent_id) REFERENCES {POSTGRES_SCHEMA}.{TABLE_NAME}(id) ON DELETE CASCADE
);