package dal

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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

var (
	ErrUnimplemented = errors.New("method not implemented")
)

type ListenerDAL struct {
	db     *sql.DB
	schema string
}

func NewListenerDAL(db *sql.DB, schema string) IListenerDAL {
	return &ListenerDAL{db: db, schema: schema}
}

func (dal *ListenerDAL) CreateListener(ctx context.Context, listener *models.Listener) error {
	query := fmt.Sprintf(`
        INSERT INTO %s.listeners (kind, status, name, description, is_external, host, port, heartbeat, config) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		logger.Debug(logLevel, logDetailListener, fmt.Sprintf("Creating listener: %s", listener.ID))

		_, err := stmt.ExecContext(ctx, listener.Kind, listener.Status, listener.Name, listener.Description,
			listener.External, listener.Host, listener.Port, listener.Heartbeat, listener.RawConfig)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to create listener: ", err)
			return fmt.Errorf("failed to create listener: %w", err)
		}
		return nil
	})
}

func (dal *ListenerDAL) GetListenerById(ctx context.Context, listenerID string) (models.Listener, error) {
	query := fmt.Sprintf(`
        SELECT id, kind, status, name, description, is_external, host, port, heartbeat, config, created_at, updated_at, started_at, stopped_at 
        FROM %s.listeners WHERE id = $1`, dal.schema)

	var listener models.Listener
	err := utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		row := stmt.QueryRowContext(ctx, listenerID)

		var startedAt, stoppedAt sql.NullTime

		err := row.Scan(
			&listener.ID,
			&listener.Kind,
			&listener.Status,
			&listener.Name,
			&listener.Description,
			&listener.External,
			&listener.Host,
			&listener.Port,
			&listener.Heartbeat,
			&listener.RawConfig,
			&listener.CreatedAt,
			&listener.UpdatedAt,
			&startedAt,
			&stoppedAt,
		)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to get listener: ", err)
			return fmt.Errorf("failed to get listener: %w", err)
		}

		if startedAt.Valid {
			listener.StartedAt = startedAt.Time
		}
		if stoppedAt.Valid {
			listener.StoppedAt = stoppedAt.Time
		}

		return nil
	})

	return listener, err
}

func (dal *ListenerDAL) GetAllListeners(ctx context.Context) ([]models.Listener, error) {
	query := fmt.Sprintf(`
        SELECT id, kind, status, name, description, is_external, host, port, heartbeat, config, created_at, updated_at, started_at, stopped_at 
        FROM %s.listeners`, dal.schema)

	var listeners []models.Listener
	err := utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		rows, err := stmt.QueryContext(ctx)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to get listeners: ", err)
			return fmt.Errorf("failed to get listeners: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var listener models.Listener
			var startedAt, stoppedAt sql.NullTime

			err := rows.Scan(
				&listener.ID,
				&listener.Kind,
				&listener.Status,
				&listener.Name,
				&listener.Description,
				&listener.External,
				&listener.Host,
				&listener.Port,
				&listener.Heartbeat,
				&listener.RawConfig,
				&listener.CreatedAt,
				&listener.UpdatedAt,
				&startedAt,
				&stoppedAt,
			)
			if err != nil {
				logger.Error(logLevel, logDetailListener, "Failed to scan listener: ", err)
				return fmt.Errorf("failed to scan listener: %w", err)
			}

			if startedAt.Valid {
				listener.StartedAt = startedAt.Time
			}
			if stoppedAt.Valid {
				listener.StoppedAt = stoppedAt.Time
			}
			listeners = append(listeners, listener)
		}

		return nil
	})

	return listeners, err
}

func (dal *ListenerDAL) DeleteListener(ctx context.Context, listenerID string) error {
	query := fmt.Sprintf(`
		DELETE FROM %s.listeners WHERE id = $1`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		logger.Debug(logLevel, logDetailListener, fmt.Sprintf("Deleting listener: %s", listenerID))

		_, err := stmt.ExecContext(ctx, listenerID)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to delete listener: ", err)
			return fmt.Errorf("failed to delete listener: %w", err)
		}

		return nil
	})
}

func (dal *ListenerDAL) UpdateListener(ctx context.Context, listenerID string, updates map[string]any) error {
	if len(updates) == 0 {
		return errors.New("no updates provided")
	}

	var setClauses []string
	var values []any
	paramIndex := 1

	for key, value := range updates {
		if key == "config" { // Convert config to JSON before updating
			configJSON, err := json.Marshal(value)
			if err != nil {
				logger.Error(logLevel, logDetailListener, "Failed to marshal listener config: ", err)
				return fmt.Errorf("failed to marshal listener config: %w", err)
			}
			value = configJSON
		}

		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", key, paramIndex))
		values = append(values, value)
		paramIndex++
	}

	// Add listenerID as the last parameter
	values = append(values, listenerID)

	query := fmt.Sprintf(`UPDATE %s.listeners SET %s WHERE id = $%d`,
		dal.schema, strings.Join(setClauses, ", "), paramIndex)

	logger.Debug(logLevel, logDetailListener, fmt.Sprintf("Updating listener: %s with query: %s", listenerID, query))

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx, values...)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to update listener: ", err)
			return fmt.Errorf("failed to update listener: %w", err)
		}

		return nil
	})
}

func (dal *ListenerDAL) GetActiveListeners(ctx context.Context) ([]models.Listener, error) {
	query := fmt.Sprintf(`
        SELECT id, kind, status, name, description, is_external, host, port, heartbeat, config, created_at, updated_at, started_at, stopped_at 
        FROM %s.listeners WHERE status = '%s'`, dal.schema, models.StatusRunning)

	var listeners []models.Listener
	err := utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		rows, err := stmt.QueryContext(ctx)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to get active listeners: ", err)
			return fmt.Errorf("failed to get active listeners: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var listener models.Listener
			var startedAt, stoppedAt sql.NullTime

			err := rows.Scan(
				&listener.ID,
				&listener.Kind,
				&listener.Status,
				&listener.Name,
				&listener.Description,
				&listener.External,
				&listener.Host,
				&listener.Port,
				&listener.Heartbeat,
				&listener.RawConfig,
				&listener.CreatedAt,
				&listener.UpdatedAt,
				&startedAt,
				&stoppedAt,
			)
			if err != nil {
				logger.Error(logLevel, logDetailListener, "Failed to scan listener: ", err)
				return fmt.Errorf("failed to scan listener: %w", err)
			}

			if startedAt.Valid {
				listener.StartedAt = startedAt.Time
			}
			if stoppedAt.Valid {
				listener.StoppedAt = stoppedAt.Time
			}
			listeners = append(listeners, listener)
		}

		return nil
	})

	return listeners, err
}

func (dal *ListenerDAL) GetListenerByName(ctx context.Context, name string) (models.Listener, error) {
	query := fmt.Sprintf(`
        SELECT id, kind, status, name, description, is_external, host, port, heartbeat, config, created_at, updated_at, started_at, stopped_at 
        FROM %s.listeners WHERE name = $1`, dal.schema)

	var listener models.Listener
	var startedAt, stoppedAt sql.NullTime
	err := utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		row := stmt.QueryRowContext(ctx, name)

		err := row.Scan(
			&listener.ID,
			&listener.Kind,
			&listener.Status,
			&listener.Name,
			&listener.Description,
			&listener.External,
			&listener.Host,
			&listener.Port,
			&listener.Heartbeat,
			&listener.RawConfig,
			&listener.CreatedAt,
			&listener.UpdatedAt,
			&startedAt,
			&stoppedAt,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				logger.Debug(logLevel, logDetailListener, fmt.Sprintf("No listener found with name: %s", name))
				return nil
			}

			logger.Error(logLevel, logDetailListener, "Failed to get listener by name: ", err)
			return fmt.Errorf("failed to get listener by name: %w", err)
		}

		if startedAt.Valid {
			listener.StartedAt = startedAt.Time
		}
		if stoppedAt.Valid {
			listener.StoppedAt = stoppedAt.Time
		}

		return nil
	})

	if err != nil {
		return models.Listener{}, err
	}

	return listener, nil
}
