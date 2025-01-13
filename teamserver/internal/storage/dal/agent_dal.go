package dal

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/models"
)

type IAgentDAL interface {
	GetAgent(agentID string) (models.Agent, error)
	UpdateAgent(ctx context.Context, agent models.UpdateAgentRequest) (models.Agent, error)
	DeleteAgent(ctx context.Context, agentID string) error
	CreateAgentTask(ctx context.Context, task models.AgentTask) error
	GetAgentTasks(ctx context.Context, agentID string) ([]models.AgentTask, error)
	DeleteAgentTask(ctx context.Context, agentID string, taskID string) error
	DeleteAgentTasks(ctx context.Context, agentID string) error
	CreateAgentConfig(ctx context.Context, agentConfig models.AgentConfig) error
	GetAgentConfig(ctx context.Context, agentID string) (models.AgentConfig, error)
	UpdateAgentConfig(ctx context.Context, agentID string, agentConfig models.AgentConfig) error
	DeleteAgentConfig(ctx context.Context, agentID string) error
	CreateAgentInfo(ctx context.Context, agent models.AgentInfo) error
	UpdateAgentInfo(ctx context.Context, agent models.AgentInfo) error
	GetAgentInfo(ctx context.Context, agentID string) (models.AgentInfo, error)
	DeleteAgentInfo(ctx context.Context, agentID string) error
	UpdateAgentLastCallback(ctx context.Context, agentID string, lastCallback string) error
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
		&agent.AgentID, &agent.Name, &agent.Note, &agent.Status,
		&agent.FirstCallback, &agent.LastCallback, &agent.ModifiedAt)

	if err == sql.ErrNoRows {
		return models.Agent{}, fmt.Errorf("agent not found")
	}
	if err != nil {
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

	agentQuery := fmt.Sprintf(`
        UPDATE %s.agents 
        SET name = $1, note = $2, status = $3, modified_at = $4
        WHERE id = $5
		RETURNING id, name, note, status, first_callback, last_callback, modified_at`, dal.schema)

	var updatedAgent models.Agent
	if err = tx.QueryRowContext(ctx, agentQuery,
		agent.Name, agent.Note, agent.Status, agent.ModifiedAt, agent.AgentID,
	).Scan(
		&updatedAgent.AgentID,
		&updatedAgent.Name,
		&updatedAgent.Note,
		&updatedAgent.Status,
		&updatedAgent.FirstCallback,
		&updatedAgent.LastCallback,
		&updatedAgent.ModifiedAt,
	); err != nil {
		return models.Agent{}, fmt.Errorf("failed to update agent: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return models.Agent{}, fmt.Errorf("failed to execte transaction: %w", err)
	}

	return updatedAgent, nil
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
	// Convert task.Command to JSON
	commandJSON, err := json.Marshal(task.Command)
	if err != nil {
		return fmt.Errorf("failed to marshal command to JSON: %w", err)
	}

	query := fmt.Sprintf(`
        INSERT INTO %s.agent_task (task_id, agent_id, type, status, module, command, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`, dal.schema)

	_, err = dal.db.ExecContext(ctx, query,
		task.TaskID, task.AgentID, task.Type, task.Status, task.Module,
		commandJSON, task.Created)
	if err != nil {
		return fmt.Errorf("failed to create agent task: %w", err)
	}
	return nil
}

func (dal *AgentDAL) GetAgentTasks(ctx context.Context, agentID string) ([]models.AgentTask, error) {
	query := fmt.Sprintf(`
        SELECT task_id, agent_id, type, status, module, command, 
               created_at, started_at, finished_at
        FROM %s.agent_task 
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
		var commandJSON []byte

		err := rows.Scan(&task.TaskID, &task.AgentID, &task.Type, &task.Status, &task.Module,
			&commandJSON, &task.Created, &task.Started, &task.Finished)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert JSON to task.Command
		err = json.Unmarshal(commandJSON, &task.Command)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal command from JSON: %w", err)
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return tasks, nil
}

func (dal *AgentDAL) DeleteAgentTask(ctx context.Context, agentID, taskID string) error {
	query := fmt.Sprintf(`
        DELETE FROM %s.agent_task 
        WHERE agent_id = $1 AND task_id = $2`, dal.schema)

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
        DELETE FROM %s.agent_task 
        WHERE agent_id = $1`, dal.schema)

	_, err := dal.db.ExecContext(ctx, query, agentID)
	if err != nil {
		return fmt.Errorf("failed to delete agent tasks: %w", err)
	}
	return nil
}

func (dal *AgentDAL) CreateAgentConfig(ctx context.Context, agentConfig models.AgentConfig) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.agent_config (config_id, listener_id, sleep, jitter, start_date, kill_date, working_hours_start, working_hours_end)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, dal.schema)

	_, err := dal.db.ExecContext(ctx, query,
		agentConfig.ConfigID, agentConfig.ListenerID,
		agentConfig.Sleep, agentConfig.Jitter, agentConfig.StartDate,
		agentConfig.KillDate, agentConfig.WorkingHoursStart, agentConfig.WorkingHoursEnd)
	if err != nil {
		return fmt.Errorf("failed to create agent config: %w", err)
	}
	return nil
}

func (dal *AgentDAL) GetAgentConfig(ctx context.Context, agentID string) (models.AgentConfig, error) {
	query := fmt.Sprintf(`
		SELECT *
		FROM %s.agent_config
		WHERE agent_id = $1`, dal.schema)

	var agentConfig models.AgentConfig
	err := dal.db.QueryRowContext(ctx, query, agentID).Scan(
		&agentConfig.ConfigID, &agentConfig.ListenerID, &agentConfig.Sleep,
		&agentConfig.Jitter, &agentConfig.StartDate, &agentConfig.KillDate,
		&agentConfig.WorkingHoursStart, &agentConfig.WorkingHoursEnd)
	if err != nil {
		return models.AgentConfig{}, fmt.Errorf("failed to get agent config: %w", err)
	}
	return agentConfig, nil
}

func (dal *AgentDAL) UpdateAgentConfig(ctx context.Context, agentID string, agentConfig models.AgentConfig) error {
	query := fmt.Sprintf(`
        UPDATE %s.agent_config
        SET config_id = $1, listener_id = $2, sleep = $3, jitter = $4, start_date = $5, kill_date = $6,
            working_hours_start = $7, working_hours_end = $8
        WHERE agent_id = $9`, dal.schema)

	_, err := dal.db.ExecContext(ctx, query, agentConfig.ConfigID, agentConfig.ListenerID, agentConfig.Sleep, agentConfig.Jitter,
		agentConfig.StartDate, agentConfig.KillDate, agentConfig.WorkingHoursStart, agentConfig.WorkingHoursEnd, agentID)
	if err != nil {
		return fmt.Errorf("failed to update agent config: %w", err)
	}
	return nil
}

func (dal *AgentDAL) DeleteAgentConfig(ctx context.Context, agentID string) error {
	query := fmt.Sprintf(`
        DELETE FROM %s.agent_config
        WHERE id = $1`, dal.schema)

	_, err := dal.db.ExecContext(ctx, query, agentID)
	if err != nil {
		return fmt.Errorf("failed to delete agent config: %w", err)
	}
	return nil
}

func (dal *AgentDAL) CreateAgentInfo(ctx context.Context, agent models.AgentInfo) error {
	query := fmt.Sprintf(`
        INSERT INTO %s.agent_info (agent_id, host_name, ip_address, user_name, system_info, os_info)
        VALUES ($1, $2, $3, $4, $5, $6)`, dal.schema)

	_, err := dal.db.ExecContext(ctx, query, agent.AgentID, agent.HostName, agent.IPAddress, agent.Username, agent.SystemInfo, agent.OSInfo)
	if err != nil {
		return fmt.Errorf("failed to set agent info: %w", err)
	}
	return nil
}

func (dal *AgentDAL) UpdateAgentInfo(ctx context.Context, agent models.AgentInfo) error {
	query := fmt.Sprintf(`
        UPDATE %s.agent_info
        SET host_name = $1, ip_address = $2, user_name = $3, system_info = $4, os_info = $5
        WHERE agent_id = $6`, dal.schema)

	_, err := dal.db.ExecContext(ctx, query, agent.HostName, agent.IPAddress, agent.Username, agent.SystemInfo, agent.OSInfo, agent.AgentID)
	if err != nil {
		return fmt.Errorf("failed to update agent info: %w", err)
	}
	return nil
}

func (dal *AgentDAL) GetAgentInfo(ctx context.Context, agentID string) (models.AgentInfo, error) {
	query := fmt.Sprintf(`
        SELECT agent_id, host_name, ip_address, user_name, system_info, os_info
        FROM %s.agent_info
        WHERE agent_id = $1`, dal.schema)

	var agentInfo models.AgentInfo
	err := dal.db.QueryRowContext(ctx, query, agentID).Scan(
		&agentInfo.AgentID, &agentInfo.HostName, &agentInfo.IPAddress, &agentInfo.Username, &agentInfo.SystemInfo, &agentInfo.OSInfo)
	if err != nil {
		return models.AgentInfo{}, fmt.Errorf("failed to get agent info: %w", err)
	}
	return agentInfo, nil
}

func (dal *AgentDAL) DeleteAgentInfo(ctx context.Context, agentID string) error {
	query := fmt.Sprintf(`
        DELETE FROM %s.agent_info
        WHERE agent_id = $1`, dal.schema)

	_, err := dal.db.ExecContext(ctx, query, agentID)
	if err != nil {
		return fmt.Errorf("failed to delete agent info: %w", err)
	}
	return nil
}

func (dal *AgentDAL) UpdateAgentLastCallback(ctx context.Context, agentID string, lastCallback string) error {
	query := fmt.Sprintf(`
		UPDATE %s.agents
		SET last_callback = $1
		WHERE id = $2`, dal.schema)

	_, err := dal.db.ExecContext(ctx, query, lastCallback, agentID)
	if err != nil {
		return fmt.Errorf("failed to update last callback: %w", err)
	}
	return nil
}
