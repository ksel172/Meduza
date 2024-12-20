package dal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/models"
)

type IAgentDAL interface {
	GetAgent(agentID string) (models.Agent, error)
	UpdateAgent(ctx context.Context, agent models.Agent) error
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

	var agent models.Agent
	err := dal.db.QueryRow(query, agentID).Scan(
		&agent.ID, &agent.Name, &agent.Note, &agent.Status,
		&agent.FirstCallback, &agent.LastCallback, &agent.ModifiedAt)

	if err == sql.ErrNoRows {
		return models.Agent{}, fmt.Errorf("agent not found")
	}
	if err != nil {
		return models.Agent{}, fmt.Errorf("failed to get agent: %w", err)
	}

	return agent, nil
}

func (dal *AgentDAL) UpdateAgent(ctx context.Context, agent models.Agent) error {
	tx, err := dal.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	agentQuery := fmt.Sprintf(`
        UPDATE %s.agents 
        SET name = $1, note = $2, status = $3, modified_at = $4
        WHERE id = $5`, dal.schema)

	_, err = tx.ExecContext(ctx, agentQuery,
		agent.Name, agent.Note, agent.Status, agent.ModifiedAt, agent.ID)
	if err != nil {
		return fmt.Errorf("failed to update agent: %w", err)
	}

	return tx.Commit()
}

func (dal *AgentDAL) DeleteAgent(ctx context.Context, agentID string) error {
	tx, err := dal.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	agentQuery := fmt.Sprintf(`DELETE FROM %s.agents WHERE id = $1`, dal.schema)
	result, err := tx.ExecContext(ctx, agentQuery, agentID)
	if err != nil {
		return fmt.Errorf("failed to delete agent: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("agent not found")
	}

	return tx.Commit()
}

func (dal *AgentDAL) CreateAgentTask(ctx context.Context, task models.AgentTask) error {
	query := fmt.Sprintf(`
        INSERT INTO %s.agent_tasks (id, agent_id, type, status, module, command, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`, dal.schema)

	_, err := dal.db.ExecContext(ctx, query,
		task.ID, task.AgentID, task.Type, task.Status, task.Module,
		task.Command, task.Created)
	if err != nil {
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

	result, err := dal.db.ExecContext(ctx, query, agentID, taskID)
	if err != nil {
		return fmt.Errorf("failed to delete agent task: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func (dal *AgentDAL) DeleteAgentTasks(ctx context.Context, agentID string) error {
	query := fmt.Sprintf(`
        DELETE FROM %s.agent_tasks 
        WHERE agent_id = $1`, dal.schema)

	_, err := dal.db.ExecContext(ctx, query, agentID)
	if err != nil {
		return fmt.Errorf("failed to delete agent tasks: %w", err)
	}
	return nil
}
