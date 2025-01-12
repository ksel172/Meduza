package dal

import (

	//standard
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	// internal
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
)

type IListenerDAL interface {
	CreateListener(context.Context, *models.Listener) error
	GetListenerById(context.Context, string) (models.Listener, error)
	GetAllListeners(context.Context) ([]models.Listener, error)
	DeleteListener(context.Context, string) error
	UpdateListener(context.Context, string, map[string]any) error
}

type ListenerDAL struct {
	db     *sql.DB
	schema string
}

func NewListenerDAL(db *sql.DB, schema string) IListenerDAL {
	return &ListenerDAL{db: db, schema: schema}
}

func (dal *ListenerDAL) CreateListener(ctx context.Context, listener *models.Listener) error {
	config, err := json.Marshal(listener.Config)
	if err != nil {
		logger.Error(layer, "Failed to marshal listener config: ", err)
		return fmt.Errorf("failed to marshal listener config: %w", err)
	}

	logging, err := json.Marshal(listener.Logging)
	if err != nil {
		logger.Error(layer, "Failed to marshal listener logging: ", err)
		return fmt.Errorf("failed to marshal listener logging: %w", err)
	}

	logger.Debug(layer, "Creating listener: "+listener.ID.String())
	query := fmt.Sprintf(`INSERT INTO %s.listeners (type, name, status, description, config, logging_enabled, logging, created_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8)`, dal.schema)
	_, err = dal.db.ExecContext(ctx, query, listener.Type, listener.Name, listener.Status, listener.Description, config, listener.LoggingEnabled, logging, time.Now().UTC())
	if err != nil {
		logger.Error(layer, "Failed to create listener: ", err)
		return fmt.Errorf("failed to create listener: %w", err)
	}
	return nil
}

func (dal *ListenerDAL) GetListenerById(ctx context.Context, lId string) (models.Listener, error) {
	logger.Debug(layer, "Getting listener by id: "+lId)
	query := fmt.Sprintf(`SELECT id, type, name, status, description, config, logging_enabled, logging, created_at, updated_at, started_at, stopped_at FROM %s.listeners WHERE id=$1`, dal.schema)
	row := dal.db.QueryRowContext(ctx, query, lId)

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
			logger.Error(layer, "Listener not found", err)
			return models.Listener{}, fmt.Errorf("unable to find the listener with id: %s", lId)
		}
		logger.Error(layer, "Error retrieving listener", err)
		return models.Listener{}, fmt.Errorf("failed to get listener")
	}

	if err := json.Unmarshal(rawConfig, &listener.Config); err != nil {
		logger.Error(layer, "Error unmarshalling listener config", err)
		return models.Listener{}, fmt.Errorf("failed to parse Config field")
	}

	if err := json.Unmarshal(rawLogging, &listener.Logging); err != nil {
		logger.Error(layer, "Error unmarshalling Logging", err)
		return models.Listener{}, fmt.Errorf("failed to parse Logging field")
	}

	return listener, nil
}

func (dal *ListenerDAL) GetAllListeners(ctx context.Context) ([]models.Listener, error) {
	logger.Debug(layer, "Getting all listeners")
	query := fmt.Sprintf(`SELECT id, type, name, status, description, config, logging_enabled, logging, created_at, updated_at, started_at, stopped_at FROM %s.listeners ORDER BY created_at DESC`, dal.schema)
	rows, err := dal.db.QueryContext(ctx, query)
	if err != nil {
		logger.Error(layer, "Failed to get listeners: ", err)
		return nil, fmt.Errorf("failed to get listeners")
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
			logger.Error(layer, "Failed to get the listener: ", err)
			return nil, fmt.Errorf("failed to get listener")
		}

		if err := json.Unmarshal(rawConfig, &listener.Config); err != nil {
			logger.Error(layer, "Failed to unmarshal config: ", err)
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		if err := json.Unmarshal(rawLogging, &listener.Logging); err != nil {
			logger.Error(layer, "Failed to unmarshal logging: ", err)
			return nil, fmt.Errorf("failed to unmarshal logging: %w", err)
		}
		lists = append(lists, listener)
	}
	return lists, nil
}

func (dal *ListenerDAL) DeleteListener(ctx context.Context, lid string) error {
	logger.Debug(layer, "Deleting listener: "+lid)
	query := fmt.Sprintf(`DELETE FROM %s.listeners WHERE id = $1`, dal.schema)
	_, err := dal.db.ExecContext(ctx, query, lid)
	if err != nil {
		logger.Error("Unable to Delete listener: ", err)
		return fmt.Errorf("failed to delete listener")
	}
	return nil
}

func (dal *ListenerDAL) UpdateListener(ctx context.Context, lid string, updates map[string]any) error {
	setClauses := []string{}
	args := []any{}
	count := 1

	for field, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, count))
		args = append(args, value)
		count++
	}

	query := fmt.Sprintf(`UPDATE %s.listeners
      SET %s
      WHERE id = $%d`, dal.schema, strings.Join(setClauses, ", "), count)
	args = append(args, lid)

	logger.Debug(layer, "Updating listener: "+lid)
	updated, err := dal.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Error("Failed to update listener: ", err)
		return fmt.Errorf("failed to update listener")
	}

	rowsAffected, err := updated.RowsAffected()
	if err != nil {
		logger.Error(layer, "Failed to retrieve affected rows: ", err)
		return fmt.Errorf("failed to retrieve affected rows")
	}
	if rowsAffected == 0 {
		logger.Warn(layer, "No rows were updated")
		return fmt.Errorf("no listener found with id: %s", lid)
	}

	logger.Debug(layer, "Rows Affected:", rowsAffected)
	return nil
}
