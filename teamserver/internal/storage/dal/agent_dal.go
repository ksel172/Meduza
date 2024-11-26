package dal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ksel172/Meduza/teamserver/internal/models"
	redis2 "github.com/ksel172/Meduza/teamserver/internal/storage/repos"
	"github.com/ksel172/Meduza/teamserver/utils"
)

// var (
// 	ErrNotFound = errors.New("Resource not found")
// )

type AgentDAL struct {
	redis redis2.Service
}

func NewAgentDAL(redisService *redis2.Service) *AgentDAL {
	return &AgentDAL{redis: *redisService}
}

// Get returns a single agent
func (dal *AgentDAL) GetAgent(agentID string) (models.Agent, error) {
	agentJSON, err := dal.redis.JsonGet(context.Background(), "agents:"+agentID)
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

	// Try to get agent, if not succesful, return error
	// later: update err.Error() into some custom error handling to ensure the
	// returned error is of type ErrNotFound
	if _, err := dal.GetAgent(agent.ID); err != nil {
		if err.Error() == "agent not found" {
			return fmt.Errorf("cannot update non-existing agent: %w", err)
		} else {
			return fmt.Errorf("unexpected error: %w", err)
		}
	}

	if err := dal.redis.JsonSet(ctx, agent.RedisID(), agent); err != nil {
		return fmt.Errorf("failed to update agent: %w", err)
	}

	return nil
}

func (dal *AgentDAL) DeleteAgent(ctx context.Context, agentID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := dal.redis.JsonDelete(ctx, "agents:"+agentID); err != nil {
		return fmt.Errorf("failed to delete agent: %w", err)
	}

	return nil
}

func (dal *AgentDAL) CreateAgentTask(ctx context.Context, agentTask models.AgentTask) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := dal.redis.JsonSet(ctx, agentTask.RedisID(), agentTask); err != nil {
		return fmt.Errorf("failed to create agent task: %w", err)
	}

	return nil
}
func (dal *AgentDAL) GetAgentTasks(ctx context.Context, agentID string) ([]models.AgentTask, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Get all tasks by partial key
	tasks, err := dal.redis.GetAllByPartial(ctx, "tasks:"+agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent tasks: %w", err)
	}

	tasksModel := make([]models.AgentTask, 0, len(tasks))
	for _, task := range tasks {
		var agentTask models.AgentTask

		if err := json.Unmarshal([]byte(task.(string)), &agentTask); err != nil {
			return nil, fmt.Errorf("failed to unmarshal agent task: %w", err)
		}
		tasksModel = append(tasksModel, agentTask)
	}

	return tasksModel, nil
}

func (dal *AgentDAL) DeleteAgentTask(ctx context.Context, agentID, taskID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := dal.redis.JsonDelete(ctx, "tasks:"+agentID+":"+taskID); err != nil {
		return fmt.Errorf("failed to delete agent task: %w", err)
	}

	return nil
}

func (dal *AgentDAL) DeleteAgentTasks(ctx context.Context, agentID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Get all tasks by partial key
	if err := dal.redis.DeleteAllByPartial(ctx, "tasks:"+agentID); err != nil {
		return fmt.Errorf("failed to delete agent tasks: %w", err)
	}

	return nil
}
