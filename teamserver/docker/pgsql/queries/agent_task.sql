CREATE TABLE IF NOT EXISTS {POSTGRES_SCHEMA}.agent_task (
    id UUID PRIMARY KEY,
    agent_id UUID NOT NULL REFERENCES meduza.agents(id),
    type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    module VARCHAR(255),
    command TEXT NOT NULL,
    created TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started TIMESTAMPTZ,
    finished TIMESTAMPTZ,
    FOREIGN KEY (agent_id) REFERENCES {POSTGRES_SCHEMA}.agents(id) ON DELETE CASCADE
<<<<<<< HEAD
);
=======
);
>>>>>>> 98d79f3 (changes in listener)
