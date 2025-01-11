package dal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
)

const layer = " [DAL] "

type IAgentDAL interface {
	GetAgent(agentID string) (models.Agent, error)
	UpdateAgent(ctx context.Context, agent models.UpdateAgentRequest) (models.Agent, error)
	DeleteAgent(ctx context.Context, agentID string) error
	CreateAgentTask(ctx context.Context, task models.AgentTask) error
	GetAgentTasks(ctx context.Context, agentID string) ([]models.AgentTask, error)
	DeleteAgentTask(ctx context.Context, agentID string, taskID string) error
	DeleteAgentTasks(ctx context.Context, agentID string) error
}

type AgentDAL struct {
	db     *sql.DB
	schema string
}

func NewAgentDAL(db *sql.DB, schema string) *AgentDAL {
	return &AgentDAL{
		db:     db,
		schema: schema,
	}
}

func (dal *AgentDAL) GetAgent(agentID string) (models.Agent, error) {
	query := fmt.Sprintf(`
        SELECT a.id, a.name, a.note, a.status, a.first_callback, a.last_callback, a.modified_at
        FROM %s.agents a
        WHERE a.id = $1`, dal.schema)

	logger.Debug(layer, "Querying for agentID: "+agentID)
	var agent models.Agent
	if err := dal.db.QueryRow(query, agentID).Scan(
		&agent.ID, &agent.Name, &agent.Note, &agent.Status,
		&agent.FirstCallback, &agent.LastCallback, &agent.ModifiedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return models.Agent{}, fmt.Errorf("agent not found")
		}
		logger.Error(layer, fmt.Sprintf("failed to get agent: %v", err))
		return models.Agent{}, fmt.Errorf("failed to get agent: %w", err)
	}

	return agent, nil
}

func (dal *AgentDAL) UpdateAgent(ctx context.Context, agent models.UpdateAgentRequest) (models.Agent, error) {
	tx, err := dal.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Agent{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	logger.Debug(layer, "Updating agent: "+agent.ID)
	agentQuery := fmt.Sprintf(`
        UPDATE %s.agents
        SET name = $1, note = $2, status = $3, modified_at = $4
        WHERE id = $5
		RETURNING id, name, note, status, first_callback, last_callback, modified_at`, dal.schema)

	var updatedAgent models.Agent
	if err = tx.QueryRowContext(ctx, agentQuery,
		agent.Name, agent.Note, agent.Status, agent.ModifiedAt, agent.ID,
	).Scan(
		&updatedAgent.ID,
		&updatedAgent.Name,
		&updatedAgent.Note,
		&updatedAgent.Status,
		&updatedAgent.FirstCallback,
		&updatedAgent.LastCallback,
		&updatedAgent.ModifiedAt,
	); err != nil {
		logger.Error(layer, fmt.Sprintf("failed to update agent: %v", err))
		return models.Agent{}, fmt.Errorf("failed to update agent: %w", err)
	}

	if err = tx.Commit(); err != nil {
		logger.Error(layer, fmt.Sprintf("failed to execute transaction: %v", err))
		return models.Agent{}, fmt.Errorf("failed to execute transaction: %w", err)
	}

	return updatedAgent, nil
}

func (dal *AgentDAL) DeleteAgent(ctx context.Context, agentID string) error {
	tx, err := dal.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	logger.Debug(layer, "Deleting agent: "+agentID)
	agentQuery := fmt.Sprintf(`DELETE FROM %s.agents WHERE id = $1`, dal.schema)
	result, err := tx.ExecContext(ctx, agentQuery, agentID)
	if err != nil {
		return fmt.Errorf("failed to delete agent: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.Error(layer, fmt.Sprintf("failed to get rows affected: %v", err))
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		logger.Warn(layer, "Agent not found")
		return fmt.Errorf("agent not found")
	}

	return tx.Commit()
}

func (dal *AgentDAL) CreateAgentTask(ctx context.Context, task models.AgentTask) error {
	query := fmt.Sprintf(`
        INSERT INTO %s.agent_tasks (id, agent_id, type, status, module, command, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`, dal.schema)

	logger.Debug(layer, "Creating agent task: "+task.ID)
	_, err := dal.db.ExecContext(ctx, query,
		task.ID, task.AgentID, task.Type, task.Status, task.Module,
		task.Command, task.Created)
	if err != nil {
		logger.Error(layer, fmt.Sprintf("failed to create agent task: %v", err))
		return fmt.Errorf("failed to create agent task: %w", err)
	}
	return nil
}

func (dal *AgentDAL) GetAgentTasks(ctx context.Context, agentID string) ([]models.AgentTask, error) {
	query := fmt.Sprintf(`
        SELECT id, agent_id, type, status, module, command, 
               created_at, started_at, finished_at
        FROM %s.agent_tasks 
        WHERE agent_id = $1 
        ORDER BY created_at DESC`, dal.schema)

	logger.Debug(layer, "Getting agent tasks for agent: "+agentID)
	rows, err := dal.db.QueryContext(ctx, query, agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get agent tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.AgentTask
	for rows.Next() {
		var task models.AgentTask
		err := rows.Scan(
			&task.ID, &task.AgentID, &task.Type, &task.Status,
			&task.Module, &task.Command, &task.Created,
			&task.Started, &task.Finished)
		if err != nil {
			logger.Error(layer, fmt.Sprintf("failed to scan task row: %v", err))
			return nil, fmt.Errorf("failed to scan task row: %w", err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (dal *AgentDAL) DeleteAgentTask(ctx context.Context, agentID, taskID string) error {
	query := fmt.Sprintf(`
        DELETE FROM %s.agent_tasks 
        WHERE agent_id = $1 AND id = $2`, dal.schema)

	logger.Debug(layer, "Deleting agent task: "+taskID)
	result, err := dal.db.ExecContext(ctx, query, agentID, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete agent task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.Error(layer, fmt.Sprintf("failed to get rows affected: %v", err))
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		logger.Warn(layer, "Task not found")
		return fmt.Errorf("task not found")
	}
	return nil
}

func (dal *AgentDAL) DeleteAgentTasks(ctx context.Context, agentID string) error {
	query := fmt.Sprintf(`
        DELETE FROM %s.agent_tasks 
        WHERE agent_id = $1`, dal.schema)

	logger.Debug(layer, "Deleting agent tasks for agent: "+agentID)
	_, err := dal.db.ExecContext(ctx, query, agentID)
	if err != nil {
		logger.Error(layer, fmt.Sprintf("failed to delete agent tasks: %v", err))
		return fmt.Errorf("failed to delete agent tasks: %w", err)
	}
	return nil
}
