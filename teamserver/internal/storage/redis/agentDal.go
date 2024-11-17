package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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

// Get returns a single agent
func (dal *AgentDAL) GetAgent(ctx context.Context, agentID string) (models.Agent, error) {
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

func (dal *AgentDAL) UpdateAgent(ctx context.Context, agent models.Agent) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := dal.redis.JsonSet(ctx, agent.ID, agent); err != nil {
		return fmt.Errorf("failed to update agent: %w", err)
	}

	return nil
}

func (dal *AgentDAL) DeleteAgent(ctx context.Context, agentID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := dal.redis.JsonDelete(ctx, agentID); err != nil {
		return fmt.Errorf("failed to delete agent: %w", err)
	}

	return nil
}

func (dal *AgentDAL) CreateAgentTask(ctx context.Context, agentTask models.AgentTask) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := dal.redis.JsonSet(ctx, agentTask.ID, agentTask); err != nil {
		return fmt.Errorf("failed to create agent task: %w", err)
	}

	return nil
}
func (dal *AgentDAL) GetAgentTasks(ctx context.Context, agentID string) ([]models.AgentTask, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// TODO
	// Get all tasks for a single agent
	// Redis stores the taskID as the key, agentID as the value

	return nil, nil
}

func (dal *AgentDAL) DeleteAgentTask(ctx context.Context, taskID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := dal.redis.JsonDelete(ctx, taskID); err != nil {
		return fmt.Errorf("failed to delete agent task: %w", err)
	}

	return nil
}

func (dal *AgentDAL) DeleteAgentTasks(ctx context.Context, agentID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// TODO
	// Delete all tasks for a single agent
	// Redis stores the taskID as the key, agentID as the value

	return nil
}
