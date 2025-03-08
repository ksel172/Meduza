package dal

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type IAgentDAL interface {
	GetAgent(ctx context.Context, agentID string) (models.Agent, error)
	GetAgents(ctx context.Context) ([]models.Agent, error)
	UpdateAgent(ctx context.Context, agent models.UpdateAgentRequest) (models.Agent, error)
	DeleteAgent(ctx context.Context, agentID string) error
	CreateAgentTask(ctx context.Context, task models.AgentTask) error
	UpdateAgentTask(ctx context.Context, task models.AgentTask) error
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

func (dal *AgentDAL) GetAgent(ctx context.Context, agentID string) (models.Agent, error) {
	query := fmt.Sprintf(`
			SELECT a.id, a.name, a.note, a.status, a.first_callback, a.last_callback, a.modified_at
			FROM %s.agents a
			WHERE a.id = $1`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) (models.Agent, error) {
		logger.Debug(logLevel, logDetailAgent, "Querying for agentID: "+agentID)
		var agent models.Agent
		if err := stmt.QueryRow(agentID).Scan(
			&agent.AgentID, &agent.Name, &agent.Note, &agent.Status,
			&agent.FirstCallback, &agent.LastCallback, &agent.ModifiedAt,
		); err != nil {
			if err == sql.ErrNoRows {
				return models.Agent{}, fmt.Errorf("agent not found")
			}
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("failed to get agent: %v", err))
			return models.Agent{}, fmt.Errorf("failed to get agent: %w", err)
		}

		return agent, nil
	})
}

func (dal *AgentDAL) GetAgents(ctx context.Context) ([]models.Agent, error) {
	query := fmt.Sprintf(`
			SELECT a.id, a.name, a.note, a.status, a.first_callback, a.last_callback, a.modified_at
			FROM %s.agents a`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) ([]models.Agent, error) {
		logger.Debug(logLevel, logDetailAgent, "Getting all agents")
		rows, err := stmt.QueryContext(ctx)
		if err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("failed to get agents: %v", err))
			return nil, fmt.Errorf("failed to get agents: %w", err)
		}
		defer rows.Close()

		var agents []models.Agent
		for rows.Next() {
			var agent models.Agent
			if err := rows.Scan(
				&agent.AgentID, &agent.Name, &agent.Note, &agent.Status,
				&agent.FirstCallback, &agent.LastCallback, &agent.ModifiedAt,
			); err != nil {
				logger.Error(logLevel, logDetailAgent, fmt.Sprintf("failed to scan agent row: %v", err))
				return nil, fmt.Errorf("failed to scan agent row: %w", err)
			}
			agents = append(agents, agent)
		}

		if err = rows.Err(); err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("failed to scan agents: %v", err))
			return nil, fmt.Errorf("failed to scan agents: %w", err)
		}

		return agents, nil
	})
}

func (dal *AgentDAL) UpdateAgent(ctx context.Context, agent models.UpdateAgentRequest) (models.Agent, error) {
	query := fmt.Sprintf(`
			UPDATE %s.agents
			SET name = $1, note = $2, status = $3, modified_at = $4
			WHERE id = $5
			RETURNING id, name, note, status, first_callback, last_callback, modified_at`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) (models.Agent, error) {
		logger.Debug(logLevel, logDetailAgent, fmt.Sprintf("Updating agent: %v", agent.AgentID))

		var updatedAgent models.Agent
		if err := stmt.QueryRowContext(ctx, agent.Name, agent.Note, agent.Status, agent.ModifiedAt, agent.AgentID).Scan(
			&updatedAgent.AgentID,
			&updatedAgent.Name,
			&updatedAgent.Note,
			&updatedAgent.Status,
			&updatedAgent.FirstCallback,
			&updatedAgent.LastCallback,
			&updatedAgent.ModifiedAt,
		); err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("failed to update agent: %v", err))
			return models.Agent{}, fmt.Errorf("failed to update agent: %w", err)
		}

		return updatedAgent, nil
	})
}

func (dal *AgentDAL) DeleteAgent(ctx context.Context, agentID string) error {
	query := fmt.Sprintf(`DELETE FROM %s.agents WHERE id = $1`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		logger.Debug(logLevel, logDetailAgent, fmt.Sprintf("Deleting agent: %v", agentID))

		result, err := stmt.ExecContext(ctx, agentID)
		if err != nil {
			return fmt.Errorf("failed to delete agent: %w", err)
		}

		rows, err := result.RowsAffected()
		if err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("failed to get rows affected: %v", err))
			return fmt.Errorf("failed to get rows affected: %w", err)
		}
		if rows == 0 {
			logger.Warn(logLevel, logDetailAgent, "Agent not found")
			return fmt.Errorf("agent not found")
		}

		return nil
	})
}

func (dal *AgentDAL) CreateAgentTask(ctx context.Context, task models.AgentTask) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.agent_task (task_id, agent_id, type, status, module, command, created_at, started_at, finished_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		commandJSON, err := json.Marshal(task.Command)
		if err != nil {
			return fmt.Errorf("failed to marshal command to JSON: %w", err)
		}

		logger.Debug(logLevel, logDetailAgent, fmt.Sprintf("Creating agent task: %s", task.TaskID))
		_, err = stmt.ExecContext(ctx, task.TaskID, task.AgentID, task.Type, task.Status,
			task.Module, commandJSON, task.Created, task.Started, task.Finished)
		if err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("failed to create agent task: %v", err))
			return fmt.Errorf("failed to create agent task: %w", err)
		}
		return nil
	})
}

func (dal *AgentDAL) UpdateAgentTask(ctx context.Context, task models.AgentTask) error {
	query := fmt.Sprintf(`
		UPDATE %s.agent_task
		SET type = $1, status = $2, module = $3, command = $4, started_at = $5, finished_at = $6
		WHERE task_id = $7 AND agent_id = $8`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		commandJSON, err := json.Marshal(task.Command)
		if err != nil {
			return fmt.Errorf("failed to marshal command to JSON: %w", err)
		}

		_, err = stmt.ExecContext(ctx, task.Type, task.Status, task.Module, commandJSON,
			task.Started, task.Finished, task.TaskID, task.AgentID)
		if err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("Failed to update agent task: %v", err))
			return fmt.Errorf("failed to update agent task: %w", err)
		}
		return nil
	})
}

func (dal *AgentDAL) GetAgentTasks(ctx context.Context, agentID string) ([]models.AgentTask, error) {
	query := fmt.Sprintf(`
			SELECT * FROM %s.agent_task 
			WHERE agent_id = $1 
			ORDER BY created_at DESC`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) ([]models.AgentTask, error) {
		logger.Debug(logLevel, logDetailAgent, "Getting agent tasks for agent: "+agentID)
		rows, err := stmt.QueryContext(ctx, agentID)
		if err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("Failed to get agent tasks: %v", err))
			return nil, fmt.Errorf("failed to get agent tasks: %w", err)
		}
		defer rows.Close()

		var tasks []models.AgentTask
		for rows.Next() {
			var task models.AgentTask
			var commandJSON []byte

			// Nullable fields
			var nullableModule sql.NullString
			var nullableStarted sql.NullTime
			var nullableFinished sql.NullTime

			err := rows.Scan(
				&task.TaskID,
				&task.AgentID,
				&task.Type,
				&task.Status,
				&nullableModule,
				&commandJSON,
				&task.Created,
				&nullableStarted,
				&nullableFinished,
			)
			if err != nil {
				logger.Error(logLevel, logDetailAgent, fmt.Sprintf("Failed to scan row: %v", err))
				return nil, fmt.Errorf("failed to scan row: %w", err)
			}

			// Handle nullable fields
			task.Module = nullableModule.String
			if !nullableModule.Valid {
				task.Module = "" // Default value if null
			}

			if nullableStarted.Valid {
				task.Started = nullableStarted.Time
			}

			if nullableFinished.Valid {
				task.Finished = nullableFinished.Time
			}

			// Convert JSON to task.Command
			err = json.Unmarshal(commandJSON, &task.Command)
			if err != nil {
				logger.Error(logLevel, logDetailAgent, fmt.Sprintf("Failed to scan task row: %v", err))
				return nil, fmt.Errorf("failed to scan task row: %w", err)
			}

			tasks = append(tasks, task)
		}

		if err = rows.Err(); err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("Failed to scan agent tasks: %v", err))
			return nil, fmt.Errorf("rows iteration error: %w", err)
		}

		return tasks, nil
	})
}

func (dal *AgentDAL) DeleteAgentTask(ctx context.Context, agentID, taskID string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.agent_task 
		WHERE agent_id = $1 AND task_id = $2`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		logger.Debug(logLevel, logDetailAgent, fmt.Sprintf("Deleting agent task: %s", taskID))
		result, err := stmt.ExecContext(ctx, agentID, taskID)
		if err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("failed to delete agent task: %v", err))
			return fmt.Errorf("failed to delete agent task: %w", err)
		}

		rows, err := result.RowsAffected()
		if err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("failed to get rows affected: %v", err))
			return fmt.Errorf("failed to get rows affected: %w", err)
		}
		if rows == 0 {
			logger.Warn(logLevel, logDetailAgent, "Task not found")
			return fmt.Errorf("task not found")
		}
		return nil
	})
}

func (dal *AgentDAL) DeleteAgentTasks(ctx context.Context, agentID string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.agent_task 
		WHERE agent_id = $1`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		logger.Debug(logLevel, logDetailAgent, fmt.Sprintf("Deleting agent tasks for agent: %s", agentID))
		_, err := stmt.ExecContext(ctx, agentID)
		if err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("Failed to delete agent tasks: %v", err))
			return fmt.Errorf("failed to delete agent tasks: %w", err)
		}
		return nil
	})
}

func (dal *AgentDAL) CreateAgentConfig(ctx context.Context, agentConfig models.AgentConfig) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.agent_config (config_id, listener_id, sleep, jitter, start_date, kill_date, working_hours_start, working_hours_end)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx,
			agentConfig.ConfigID, agentConfig.ListenerID,
			agentConfig.Sleep, agentConfig.Jitter, agentConfig.StartDate,
			agentConfig.KillDate, agentConfig.WorkingHoursStart, agentConfig.WorkingHoursEnd)
		if err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("Failed to create agent config: %v", err))
			return fmt.Errorf("failed to create agent config: %w", err)
		}
		return nil
	})
}

func (dal *AgentDAL) GetAgentConfig(ctx context.Context, agentID string) (models.AgentConfig, error) {
	query := fmt.Sprintf(`
		SELECT *
		FROM %s.agent_config
		WHERE agent_id = $1`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) (models.AgentConfig, error) {
		var agentConfig models.AgentConfig
		err := stmt.QueryRowContext(ctx, agentID).Scan(
			&agentConfig.ConfigID, &agentConfig.ListenerID, &agentConfig.Sleep,
			&agentConfig.Jitter, &agentConfig.StartDate, &agentConfig.KillDate,
			&agentConfig.WorkingHoursStart, &agentConfig.WorkingHoursEnd)
		if err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("Failed to get agent config: %v", err))
			return models.AgentConfig{}, fmt.Errorf("failed to get agent config: %w", err)
		}
		return agentConfig, nil
	})
}

func (dal *AgentDAL) UpdateAgentConfig(ctx context.Context, agentID string, agentConfig models.AgentConfig) error {
	query := fmt.Sprintf(`
		UPDATE %s.agent_config
		SET config_id = $1, listener_id = $2, sleep = $3, jitter = $4, start_date = $5, kill_date = $6,
			working_hours_start = $7, working_hours_end = $8
		WHERE agent_id = $9`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx, agentConfig.ConfigID, agentConfig.ListenerID, agentConfig.Sleep, agentConfig.Jitter,
			agentConfig.StartDate, agentConfig.KillDate, agentConfig.WorkingHoursStart, agentConfig.WorkingHoursEnd, agentID)
		if err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("Failed to update agent config: %v", err))
			return fmt.Errorf("failed to update agent config: %w", err)
		}
		return nil
	})
}

func (dal *AgentDAL) DeleteAgentConfig(ctx context.Context, agentID string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.agent_config
		WHERE id = $1`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx, agentID)
		if err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("Failed to delete agent config: %v", err))
			return fmt.Errorf("failed to delete agent config: %w", err)
		}
		return nil
	})
}

func (dal *AgentDAL) CreateAgentInfo(ctx context.Context, agent models.AgentInfo) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.agent_info (agent_id, host_name, ip_address, user_name, system_info, os_info)
		VALUES ($1, $2, $3, $4, $5, $6)`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx, agent.AgentID, agent.HostName, agent.IPAddress, agent.Username, agent.SystemInfo, agent.OSInfo)
		if err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("Failed to create agent info: %v", err))
			return fmt.Errorf("failed to create agent info: %w", err)
		}
		return nil
	})
}

func (dal *AgentDAL) UpdateAgentInfo(ctx context.Context, agent models.AgentInfo) error {
	query := fmt.Sprintf(`
		UPDATE %s.agent_info
		SET host_name = $1, ip_address = $2, user_name = $3, system_info = $4, os_info = $5
		WHERE agent_id = $6`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx, agent.HostName, agent.IPAddress, agent.Username, agent.SystemInfo, agent.OSInfo, agent.AgentID)
		if err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("Failed to update agent info: %v", err))
			return fmt.Errorf("failed to update agent info: %w", err)
		}
		return nil
	})
}

func (dal *AgentDAL) GetAgentInfo(ctx context.Context, agentID string) (models.AgentInfo, error) {
	query := fmt.Sprintf(`
		SELECT agent_id, host_name, ip_address, user_name, system_info, os_info
		FROM %s.agent_info
		WHERE agent_id = $1`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) (models.AgentInfo, error) {
		var agentInfo models.AgentInfo
		err := stmt.QueryRowContext(ctx, agentID).Scan(
			&agentInfo.AgentID, &agentInfo.HostName, &agentInfo.IPAddress, &agentInfo.Username, &agentInfo.SystemInfo, &agentInfo.OSInfo)
		if err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("failed to get agent info: %v", err))
			return models.AgentInfo{}, fmt.Errorf("failed to get agent info: %w", err)
		}
		return agentInfo, nil
	})
}

func (dal *AgentDAL) DeleteAgentInfo(ctx context.Context, agentID string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.agent_info
		WHERE agent_id = $1`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx, agentID)
		if err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("failed to delete agent info: %v", err))
			return fmt.Errorf("failed to delete agent info: %w", err)
		}
		return nil
	})
}

func (dal *AgentDAL) UpdateAgentLastCallback(ctx context.Context, agentID string, lastCallback string) error {
	query := fmt.Sprintf(`
		UPDATE %s.agents
		SET last_callback = $1
		WHERE id = $2`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx, lastCallback, agentID)
		if err != nil {
			logger.Error(logLevel, logDetailAgent, fmt.Sprintf("failed to update agent last callback: %v", err))
			return fmt.Errorf("failed to update last callback: %w", err)
		}
		return nil
	})
}
