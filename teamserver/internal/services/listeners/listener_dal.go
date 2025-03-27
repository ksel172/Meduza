package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type IListenerDAL interface {
	CreateListener(context.Context, *Listener) error
	GetListenerById(context.Context, string) (*Listener, error)
	GetAllListeners(context.Context) ([]*Listener, error)
	DeleteListener(context.Context, string) error
	UpdateListener(context.Context, string, map[string]any) error
	GetActiveListeners(context.Context) ([]*Listener, error)
	GetListenerByName(context.Context, string) (*Listener, error)
}

var (
	logLevel          = "[DAL]"
	logDetailListener = "[Listener]"
	ErrUnimplemented  = errors.New("method not implemented")
)

type ListenerDAL struct {
	db     *sql.DB
	schema string
}

func NewListenerDAL(db *sql.DB, schema string) IListenerDAL {
	return &ListenerDAL{db: db, schema: schema}
}

func (dal *ListenerDAL) CreateListener(ctx context.Context, listener *Listener) error {
	query := fmt.Sprintf(`
        INSERT INTO %s.listeners (type, name, status, description, config) 
        VALUES ($1, $2, $3, $4, $5)`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		config, err := json.Marshal(listener.Config)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to marshal listener config: ", err)
			return fmt.Errorf("failed to marshal listener config: %w", err)
		}

		logger.Debug(logLevel, logDetailListener, fmt.Sprintf("Creating listener: %s", listener.ID))

		_, err = stmt.ExecContext(ctx, listener.Type, listener.Name, listener.Status, listener.Description, config)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to create listener: ", err)
			return fmt.Errorf("failed to create listener: %w", err)
		}
		return nil
	})
}

func (dal *ListenerDAL) GetListenerById(ctx context.Context, listenerID string) (*Listener, error) {
	query := fmt.Sprintf(`
		SELECT id, type, name, status, description, config, created_at, updated_at, started_at, stopped_at 
		FROM %s.listeners WHERE id = $1`, dal.schema)

	var listener Listener
	err := utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		row := stmt.QueryRowContext(ctx, listenerID)

		var config []byte
		err := row.Scan(&listener.ID, &listener.Type, &listener.Name, &listener.Status, &listener.Description, &config, &listener.CreatedAt, &listener.UpdatedAt, &listener.StartedAt, &listener.StoppedAt)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to get listener: ", err)
			return fmt.Errorf("failed to get listener: %w", err)
		}

		err = json.Unmarshal(config, &listener.Config)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to unmarshal listener config: ", err)
			return fmt.Errorf("failed to unmarshal listener config: %w", err)
		}

		return nil
	})

	return &listener, err
}

func (dal *ListenerDAL) GetAllListeners(ctx context.Context) ([]*Listener, error) {
	query := fmt.Sprintf(`
        SELECT id, type, name, status, description, config, created_at, updated_at, started_at, stopped_at 
        FROM %s.listeners`, dal.schema)

	var listeners []*Listener
	err := utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		rows, err := stmt.QueryContext(ctx)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to get listeners: ", err)
			return fmt.Errorf("failed to get listeners: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var listener Listener
			var config []byte
			err := rows.Scan(&listener.ID, &listener.Type, &listener.Name, &listener.Status, &listener.Description, &config, &listener.CreatedAt, &listener.UpdatedAt, &listener.StartedAt, &listener.StoppedAt)
			if err != nil {
				logger.Error(logLevel, logDetailListener, "Failed to scan listener: ", err)
				return fmt.Errorf("failed to scan listener: %w", err)
			}

			err = json.Unmarshal(config, &listener.Config)
			if err != nil {
				logger.Error(logLevel, logDetailListener, "Failed to unmarshal listener config: ", err)
				return fmt.Errorf("failed to unmarshal listener config: %w", err)
			}

			listeners = append(listeners, &listener)
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
	query := fmt.Sprintf(`
		UPDATE %s.listeners SET type = $1, name = $2, status = $3, description = $4, config = $5, 
		 WHERE id = $6`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		config, err := json.Marshal(updates["config"])
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to marshal listener config: ", err)
			return fmt.Errorf("failed to marshal listener config: %w", err)
		}

		logger.Debug(logLevel, logDetailListener, fmt.Sprintf("Updating listener: %s", listenerID))

		_, err = stmt.ExecContext(ctx, updates["type"], updates["name"], updates["status"], updates["description"], config, listenerID)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to update listener: ", err)
			return fmt.Errorf("failed to update listener: %w", err)
		}

		return nil
	})
}

func (dal *ListenerDAL) GetActiveListeners(ctx context.Context) ([]*Listener, error) {
	query := fmt.Sprintf(`
        SELECT id, type, name, status, description, config, created_at, updated_at, started_at, stopped_at 
        FROM %s.listeners WHERE status = 'running'`, dal.schema)

	var listeners []*Listener
	err := utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		rows, err := stmt.QueryContext(ctx)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to get active listeners: ", err)
			return fmt.Errorf("failed to get active listeners: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var listener Listener
			var config []byte
			err := rows.Scan(&listener.ID, &listener.Type, &listener.Name, &listener.Status, &listener.Description, &config, &listener.CreatedAt, &listener.UpdatedAt, &listener.StartedAt, &listener.StoppedAt)
			if err != nil {
				logger.Error(logLevel, logDetailListener, "Failed to scan listener: ", err)
				return fmt.Errorf("failed to scan listener: %w", err)
			}

			err = json.Unmarshal(config, &listener.Config)
			if err != nil {
				logger.Error(logLevel, logDetailListener, "Failed to unmarshal listener config: ", err)
				return fmt.Errorf("failed to unmarshal listener config: %w", err)
			}

			listeners = append(listeners, &listener)
		}

		return nil
	})

	return listeners, err
}

func (dal *ListenerDAL) GetListenerByName(ctx context.Context, name string) (*Listener, error) {
	query := fmt.Sprintf(`
        SELECT id, type, name, status, description, config, created_at, updated_at, started_at, stopped_at 
        FROM %s.listeners WHERE name = $1`, dal.schema)

	var listener Listener
	err := utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		row := stmt.QueryRowContext(ctx, name)

		var config []byte
		err := row.Scan(&listener.ID, &listener.Type, &listener.Name, &listener.Status, &listener.Description, &config, &listener.CreatedAt, &listener.UpdatedAt, &listener.StartedAt, &listener.StoppedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				logger.Debug(logLevel, logDetailListener, fmt.Sprintf("No listener found with name: %s", name))
				return nil // No rows is not an error for this function
			}
			logger.Error(logLevel, logDetailListener, "Failed to get listener by name: ", err)
			return fmt.Errorf("failed to get listener by name: %w", err)
		}

		err = json.Unmarshal(config, &listener.Config)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to unmarshal listener config: ", err)
			return fmt.Errorf("failed to unmarshal listener config: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if listener.ID == "" {
		return nil, nil
	}

	return &listener, nil
}
