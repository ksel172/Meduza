CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.agent_command(
   id BIGSERIAL PRIMARY KEY,
   agent_id UUID NOT NULL,
   name VARCHAR(255) NOT NULL,
   started TIMESTAMPTZ DEFAULT NULL,
   completed TIMESTAMPTZ DEFAULT NULL,
   parameters TEXT[] NOT NULL,
   output TEXT NOT NULL,

   CONSTRAINT fk_agent FOREIGN KEY(agent_id) REFERENCES {POSTGRES_SCHEMA}.{TABLE_NAME}(id) ON DELETE CASCADE
);