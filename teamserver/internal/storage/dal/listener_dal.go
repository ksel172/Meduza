package dal

import (

	//standard
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	// internal
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type IListenerDAL interface {
	CreateListener(context.Context, *models.Listener) error
	GetListenerById(context.Context, string) (models.Listener, error)
	GetAllListeners(context.Context) ([]models.Listener, error)
	DeleteListener(context.Context, string) error
	UpdateListener(context.Context, string, map[string]any) error
	GetActiveListeners(context.Context) ([]models.Listener, error)
	GetListenerByName(context.Context, string) (models.Listener, error)
}

type ListenerDAL struct {
	db     *sql.DB
	schema string
}

func NewListenerDAL(db *sql.DB, schema string) IListenerDAL {
	return &ListenerDAL{db: db, schema: schema}
}

func (dal *ListenerDAL) CreateListener(ctx context.Context, listener *models.Listener) error {
	query := fmt.Sprintf(
		`INSERT INTO %s.listeners (type, name, status, description, config, logging_enabled, logging)
		VALUES($1, $2, $3, $4, $5, $6, $7)`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		config, err := json.Marshal(listener.Config)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to marshal listener config: ", err)
			return fmt.Errorf("failed to marshal listener config: %w", err)
		}

		logging, err := json.Marshal(listener.Logging)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to marshal listener logging: ", err)
			return fmt.Errorf("failed to marshal listener logging: %w", err)
		}
		logger.Debug(logLevel, logDetailListener, fmt.Sprintf("Creating listener: %s", listener.ID.String()))

		_, err = stmt.ExecContext(ctx, listener.Type, listener.Name, listener.Status, listener.Description, config, listener.LoggingEnabled, logging)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to create listener: ", err)
			return fmt.Errorf("failed to create listener: %w", err)
		}
		return nil
	})
}

func (dal *ListenerDAL) GetListenerById(ctx context.Context, listenerID string) (models.Listener, error) {
	query := fmt.Sprintf(`
		SELECT id, type, name, status, description, config, logging_enabled, logging,
		created_at, updated_at, started_at, stopped_at FROM %s.listeners WHERE id=$1`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) (models.Listener, error) {
		logger.Debug(logLevel, logDetailListener, fmt.Sprintf("Getting listener by id: %s", listenerID))
		row := stmt.QueryRowContext(ctx, listenerID)

		var (
			rawConfig  json.RawMessage
			rawLogging json.RawMessage
			listener   models.Listener
		)

		if err := row.Scan(
			&listener.ID,
			&listener.Type,
			&listener.Name,
			&listener.Status,
			&listener.Description,
			&rawConfig,
			&listener.LoggingEnabled,
			&rawLogging,
			&listener.CreatedAt,
			&listener.UpdatedAt,
			&listener.StartedAt,
			&listener.StoppedAt,
		); err != nil {
			if err == sql.ErrNoRows {
				logger.Error(logLevel, logDetailListener, "Listener not found", err)
				return models.Listener{}, fmt.Errorf("unable to find the listener with id: %s", listenerID)
			}
			logger.Error(logLevel, fmt.Sprintf("Error retrieving listener: %v", err))
			return models.Listener{}, fmt.Errorf("failed to get listener")
		}

		if err := json.Unmarshal(rawConfig, &listener.Config); err != nil {
			logger.Error(logLevel, logDetailListener, fmt.Sprintf("Error unmarshalling listener config: %v", err))
			return models.Listener{}, fmt.Errorf("failed to parse Config field")
		}

		if err := json.Unmarshal(rawLogging, &listener.Logging); err != nil {
			logger.Error(logLevel, logDetailListener, fmt.Sprintf("Error unmarshalling Logging: %v", err))
			return models.Listener{}, fmt.Errorf("failed to parse Logging field")
		}

		return listener, nil
	})

}

func (dal *ListenerDAL) GetAllListeners(ctx context.Context) ([]models.Listener, error) {
	query := fmt.Sprintf(`
		SELECT id, type, name, status, description, config, logging_enabled, logging, created_at,
		updated_at, started_at, stopped_at FROM %s.listeners ORDER BY created_at DESC`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) ([]models.Listener, error) {
		logger.Debug(logLevel, logDetailListener, "Getting all listeners")
		rows, err := stmt.QueryContext(ctx, query)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to get listeners: %v", err)
			return nil, fmt.Errorf("failed to get listeners: %w", err)
		}
		defer rows.Close()

		var lists []models.Listener
		for rows.Next() {
			var listener models.Listener
			var rawConfig json.RawMessage
			var rawLogging json.RawMessage
			if err := rows.Scan(&listener.ID, &listener.Type, &listener.Name, &listener.Status,
				&listener.Description, &rawConfig, &listener.LoggingEnabled, &rawLogging,
				&listener.CreatedAt, &listener.UpdatedAt, &listener.StartedAt, &listener.StoppedAt,
			); err != nil {
				logger.Error(logLevel, "Failed to get the listener: ", err)
				return nil, fmt.Errorf("failed to get listener: %w", err)
			}

			if err := json.Unmarshal(rawConfig, &listener.Config); err != nil {
				logger.Error(logLevel, "Failed to unmarshal config: ", err)
				return nil, fmt.Errorf("failed to unmarshal config: %w", err)
			}

			if err := json.Unmarshal(rawLogging, &listener.Logging); err != nil {
				logger.Error(logLevel, "Failed to unmarshal logging: ", err)
				return nil, fmt.Errorf("failed to unmarshal logging: %w", err)
			}
			lists = append(lists, listener)
		}
		return lists, nil
	})

}

func (dal *ListenerDAL) DeleteListener(ctx context.Context, listenerID string) error {
	query := fmt.Sprintf(`DELETE FROM %s.listeners WHERE id = $1`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		logger.Debug(logLevel, logDetailListener, fmt.Sprintf("Deleting listener: %v", listenerID))

		_, err := stmt.ExecContext(ctx, listenerID)
		if err != nil {
			logger.Error(logLevel, logDetailListener, fmt.Sprintf("Unable to Delete listener: %v", err))
			return fmt.Errorf("failed to delete listener: %w", err)
		}
		return nil
	})

}

func (dal *ListenerDAL) UpdateListener(ctx context.Context, listenerID string, updates map[string]any) error {
	setClauses := []string{}
	args := []any{}
	count := 1

	for field, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, count))
		args = append(args, value)
		count++
	}

	query := fmt.Sprintf(`UPDATE %s.listeners SET %s WHERE id = $%d`, dal.schema, strings.Join(setClauses, ", "), count)
	args = append(args, listenerID)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		logger.Debug(logLevel, logDetailListener, fmt.Sprintf("Updating listener: %v", listenerID))
		updated, err := stmt.ExecContext(ctx, args...)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to update listener: ", err)
			return fmt.Errorf("failed to update listener: %w", err)
		}

		rowsAffected, err := updated.RowsAffected()
		if err != nil {
			logger.Error(logLevel, logDetailListener, fmt.Sprintf("Failed to retrieve affected rows: %v", err))
			return fmt.Errorf("failed to retrieve affected rows")
		}
		if rowsAffected == 0 {
			logger.Warn(logLevel, logDetailListener, "No listener rows were updated")
			return fmt.Errorf("no listener found with id: %s", listenerID)
		}

		logger.Debug(logLevel, logDetailListener, fmt.Sprintf("Rows Affected: %v", rowsAffected))
		return nil
	})

}

func (dal *ListenerDAL) GetActiveListeners(ctx context.Context) ([]models.Listener, error) {
	query := fmt.Sprintf(`
		SELECT id, type, name, status, description, config, logging_enabled, logging, created_at, updated_at, started_at, stopped_at 
		FROM %s.listeners 
		WHERE status = $1 
		ORDER BY created_at DESC`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) ([]models.Listener, error) {
		rows, err := stmt.QueryContext(ctx, 1) // Pass the `status=1` parameter
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to get listeners\n", err)
			return nil, fmt.Errorf("failed to get listeners")
		}
		defer rows.Close()
		var lists []models.Listener
		for rows.Next() {
			var listener models.Listener
			var rawConfig json.RawMessage
			var rawLogging json.RawMessage
			if err := rows.Scan(&listener.ID, &listener.Type, &listener.Name, &listener.Status, &listener.Description, &rawConfig,
				&listener.LoggingEnabled, &rawLogging, &listener.CreatedAt, &listener.UpdatedAt, &listener.StartedAt, &listener.StoppedAt,
			); err != nil {
				logger.Error(logLevel, logDetailListener, fmt.Sprintf("Failed to get the listener: %v", err))
				return nil, fmt.Errorf("failed to get listener: %w", err)
			}

			if err := json.Unmarshal(rawConfig, &listener.Config); err != nil {
				logger.Error(logLevel, logDetailListener, fmt.Sprintf("Failed to unmarshal config: %v", err))
				return nil, fmt.Errorf("failed to unmarshal config: %w", err)
			}

			if err := json.Unmarshal(rawLogging, &listener.Logging); err != nil {
				logger.Error(logLevel, logDetailListener, fmt.Sprintf("Failed to unmarshal listener logging: %v", err))
				return nil, fmt.Errorf("failed to unmarshal logging: %w", err)
			}
			lists = append(lists, listener)
		}
		return lists, nil
	})

}

func (dal *ListenerDAL) GetListenerByName(ctx context.Context, name string) (models.Listener, error) {
	query := fmt.Sprintf(
		`SELECT id, type, name, status, description, config, logging_enabled, logging, created_at,
		updated_at, started_at, stopped_at FROM %s.listeners WHERE name=$1`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) (models.Listener, error) {
		row := stmt.QueryRowContext(ctx, name)

		var (
			rawConfig  json.RawMessage
			rawLogging json.RawMessage
			listener   models.Listener
		)

		if err := row.Scan(
			&listener.ID,
			&listener.Type,
			&listener.Name,
			&listener.Status,
			&listener.Description,
			&rawConfig,
			&listener.LoggingEnabled,
			&rawLogging,
			&listener.CreatedAt,
			&listener.UpdatedAt,
			&listener.StartedAt,
			&listener.StoppedAt,
		); err != nil {
			if err == sql.ErrNoRows {
				logger.Info(logLevel, logDetailListener, fmt.Sprintf("Listener not found by name: %v", err))
				return models.Listener{}, nil
			}
			return models.Listener{}, err
		}

		if err := json.Unmarshal(rawConfig, &listener.Config); err != nil {
			logger.Error(logLevel, logDetailListener, fmt.Sprintf("Failed to unmarshal listener config: %v", err))
			return models.Listener{}, fmt.Errorf("failed to parse listener config field: %w", err)
		}

		if err := json.Unmarshal(rawLogging, &listener.Logging); err != nil {
			logger.Error(logLevel, logDetailListener, fmt.Sprintf("Failed to unmarshal listener: %v", err))
			return models.Listener{}, fmt.Errorf("failed to parse listener logging field: %w", err)
		}

		return listener, nil
	})

}
