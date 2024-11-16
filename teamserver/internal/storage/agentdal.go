package storage

import "github.com/ksel172/Meduza/teamserver/models"

type AgentDAL struct {
	db     Database
	schema string
}

func NewAgentDAL(db Database, schema string) *AgentDAL {
	return &AgentDAL{db: db, schema: schema}
}

// RegisterAgent registers a new agent
func (dal *AgentDAL) RegisterAgent(agent models.Agent) error {
	return nil
}
