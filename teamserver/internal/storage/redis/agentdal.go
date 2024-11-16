package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/models"
)

// var (
// 	ErrNotFound = errors.New("Resource not found")
// )

type AgentDAL struct {
	redis Service
}

func NewAgentDAL(redisService *Service) *AgentDAL {
	return &AgentDAL{redis: *redisService}
}

// Register registers a new agent
func (dal *AgentDAL) Register(agent models.Agent) error {
	if _, err := dal.redis.JsonSet(context.Background(), agent.ID, agent); err != nil {
		return fmt.Errorf("Failed to register agent: %w", err)
	}
	return nil
}

// Get returns a single agent
func (dal *AgentDAL) Get(agentID string) (models.Agent, error) {
	agentJSON, err := dal.redis.JsonGet(context.Background(), agentID)
	if err != nil {
		return models.Agent{}, fmt.Errorf("failed to get agent: %w", err)
	}

	// Check if empty
	if agentJSON == "" {
		return models.Agent{}, fmt.Errorf("agent not found")
	}

	// Unmarshall
	var agent models.Agent
	json.Unmarshal([]byte(agentJSON), &agent)

	return agent, nil
}
