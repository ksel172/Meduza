package services

import (

	//standard
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
	GetListenerById(context.Context, string) (Listener, error)
	GetAllListeners(context.Context) ([]Listener, error)
	DeleteListener(context.Context, string) error
	UpdateListener(context.Context, string, map[string]any) error
	GetActiveListeners(context.Context) ([]Listener, error)
	// GetListenerByName(context.Context, string) (*Listener, error)
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

// TODO: Also have to change the .sql for these functions
// and fix pointer issues

func (dal *ListenerDAL) CreateListener(ctx context.Context, listener *Listener) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.listeners (id, config) VALUES ($1, $2)`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		config, err := json.Marshal(listener.Config)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to marshal listener config: ", err)
			return fmt.Errorf("failed to marshal listener config: %w", err)
		}

		logger.Debug(logLevel, logDetailListener, fmt.Sprintf("Creating listener: %s", listener.ID))

		_, err = stmt.ExecContext(ctx, config)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to create listener: ", err)
			return fmt.Errorf("failed to create listener: %w", err)
		}
		return nil
	})
}

func (dal *ListenerDAL) GetListenerById(ctx context.Context, listenerID string) (Listener, error) {
	query := fmt.Sprintf(`
		SELECT id, config FROM %s.listeners WHERE id = $1`, dal.schema)

	var listener Listener
	err := utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		row := stmt.QueryRowContext(ctx, listenerID)

		var config []byte
		err := row.Scan(&listener.ID, &config)
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

	return listener, err
}

func (dal *ListenerDAL) GetAllListeners(ctx context.Context) ([]Listener, error) {
	query := fmt.Sprintf(`
        SELECT id, config FROM %s.listeners`, dal.schema)

	var listeners []Listener
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
			err := rows.Scan(&listener.ID, &config)
			if err != nil {
				logger.Error(logLevel, logDetailListener, "Failed to scan listener: ", err)
				return fmt.Errorf("failed to scan listener: %w", err)
			}

			err = json.Unmarshal(config, &listener.Config) // Fixed: Added & to pass as pointer
			if err != nil {
				logger.Error(logLevel, logDetailListener, "Failed to unmarshal listener config: ", err)
				return fmt.Errorf("failed to unmarshal listener config: %w", err)
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
	query := fmt.Sprintf(`
		UPDATE %s.listeners SET config = $1 WHERE id = $2`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		config, err := json.Marshal(updates["config"])
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to marshal listener config: ", err)
			return fmt.Errorf("failed to marshal listener config: %w", err)
		}

		logger.Debug(logLevel, logDetailListener, fmt.Sprintf("Updating listener: %s", listenerID))

		_, err = stmt.ExecContext(ctx, config, listenerID)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to update listener: ", err)
			return fmt.Errorf("failed to update listener: %w", err)
		}

		return nil
	})
}

func (dal *ListenerDAL) GetActiveListeners(ctx context.Context) ([]Listener, error) {
	query := fmt.Sprintf(`
        SELECT id, config FROM %s.listeners WHERE config->>'status' = 'running'`, dal.schema)

	var listeners []Listener
	err := utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		rows, err := stmt.QueryContext(ctx)
		if err != nil {
			logger.Error(logLevel, logDetailListener, "Failed to get active listeners: ", err)
			return fmt.Errorf("failed to get active listeners: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			listener := Listener{} // Initialize a new pointer to avoid nil pointer dereference
			var config []byte
			err := rows.Scan(&listener.ID, &config)
			if err != nil {
				logger.Error(logLevel, logDetailListener, "Failed to scan listener: ", err)
				return fmt.Errorf("failed to scan listener: %w", err)
			}

			err = json.Unmarshal(config, &listener.Config)
			if err != nil {
				logger.Error(logLevel, logDetailListener, "Failed to unmarshal listener config: ", err)
				return fmt.Errorf("failed to unmarshal listener config: %w", err)
			}

			listeners = append(listeners, listener)
		}

		return nil
	})

	return listeners, err
}

// func (dal *ListenerDAL) GetListenerByName(ctx context.Context, name string) (*Listener, error) {
// 	return Listener{}, ErrUnimplemented
// }
