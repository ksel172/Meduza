package dal

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ksel172/Meduza/teamserver/pkg/listeners"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
)

type ListenersDAL struct {
	db     *sql.DB
	schema string
}

func NewListenersDAL(db *sql.DB, schema string) *ListenersDAL {
	return &ListenersDAL{db: db, schema: schema}
}

func (dal *ListenersDAL) CreateListeners(ctx context.Context, listener *listeners.Listener) error {
	config, err := json.Marshal(listener.Config)
	if err != nil {
		logger.Error("Error in Listener Dal:", err)
	}
	rr, err := json.Marshal(listener.ResponseRules)
	if err != nil {
		logger.Error("Error in Listener Dal:", err)
	}
	logging, err := json.Marshal(listener.Logging)
	if err != nil {
		logger.Error("Error in Listener Dal:", err)
	}
	query := fmt.Sprintf(`INSERT INTO %s.listeners (type, name, status, description, config, response_rules, logging_enabled, logging, created_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`, dal.schema)
	_, err = dal.db.ExecContext(ctx, query, listener.Type, listener.Name, listener.Status, listener.Description, config, rr, listener.LoggingEnabled, logging, time.Now().UTC())
	return err
}

func (dal *ListenersDAL) GetListenerById(ctx context.Context, lId string) (listeners.Listener, error) {
	// Fixing the query: Removed extra `)` after `stopped_at`
	query := fmt.Sprintf(`SELECT id, type, name, status, description, config, response_rules, logging_enabled, logging, created_at, updated_at, started_at, stopped_at FROM %s.listeners WHERE id=$1`, dal.schema)
	row := dal.db.QueryRowContext(ctx, query, lId)

	var (
		rawConfig        json.RawMessage
		rawResponseRules json.RawMessage
		rawLogging       json.RawMessage
		listener         listeners.Listener
	)

	// Scan raw JSON fields into json.RawMessage for later unmarshalling
	if err := row.Scan(
		&listener.ID,
		&listener.Type,
		&listener.Name,
		&listener.Status,
		&listener.Description,
		&rawConfig,
		&rawResponseRules,
		&listener.LoggingEnabled,
		&rawLogging,
		&listener.CreatedAt,
		&listener.UpdatedAt,
		&listener.StartedAt,
		&listener.StoppedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Listener not found", err)
			return listeners.Listener{}, fmt.Errorf("unable to find the listener with id: %s", lId)
		}
		logger.Error("Error retrieving listener", err)
		return listeners.Listener{}, fmt.Errorf("failed to get listener")
	}

	// Unmarshal JSON fields into their respective structs
	if err := json.Unmarshal(rawConfig, &listener.Config); err != nil {
		logger.Error("Error unmarshalling Config", err)
		return listeners.Listener{}, fmt.Errorf("failed to parse Config field")
	}

	if err := json.Unmarshal(rawResponseRules, &listener.ResponseRules); err != nil {
		logger.Error("Error unmarshalling ResponseRules", err)
		return listeners.Listener{}, fmt.Errorf("failed to parse ResponseRules field")
	}

	if err := json.Unmarshal(rawLogging, &listener.Logging); err != nil {
		logger.Error("Error unmarshalling Logging", err)
		return listeners.Listener{}, fmt.Errorf("failed to parse Logging field")
	}

	return listener, nil
}

func (dal *ListenersDAL) GetAllListener(ctx context.Context) ([]listeners.Listener, error) {
	query := fmt.Sprintf(`SELECT id, type, name, status, description, config, response_rules,logging_enabled, logging, created_at, updated_at, started_at, stopped_at FROM %s.listeners ORDER BY created_at DESC`, dal.schema)
	rows, err := dal.db.QueryContext(ctx, query)
	if err != nil {
		logger.Error("Failed to get listeners\n", err)
		return nil, fmt.Errorf("failed to get listeners")
	}
	defer rows.Close()
	var lists []listeners.Listener
	for rows.Next() {
		var listener listeners.Listener
		var rawConfig json.RawMessage
		var rawResponseRules json.RawMessage
		var rawLogging json.RawMessage
		if err := rows.Scan(&listener.ID, &listener.Type, &listener.Name, &listener.Status, &listener.Description, &rawConfig, &rawResponseRules, &listener.LoggingEnabled, &rawLogging, &listener.CreatedAt, &listener.UpdatedAt, &listener.StartedAt, &listener.StoppedAt); err != nil {
			logger.Error("Failed to get the listener\n", err)
			return nil, fmt.Errorf("Failed to get listener")
		}

		// Unmarshal Config
		if err := json.Unmarshal(rawConfig, &listener.Config); err != nil {
			logger.Error("Failed to unmarshal config\n", err)
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		// Unmarshal ResponseRules
		if err := json.Unmarshal(rawResponseRules, &listener.ResponseRules); err != nil {
			logger.Error("Failed to unmarshal response rules\n", err)
			return nil, fmt.Errorf("failed to unmarshal response rules: %w", err)
		}

		// Unmarshal Logging
		if err := json.Unmarshal(rawLogging, &listener.Logging); err != nil {
			logger.Error("Failed to unmarshal logging\n", err)
			return nil, fmt.Errorf("failed to unmarshal logging: %w", err)
		}
		lists = append(lists, listener)
	}
	return lists, nil
}

func (dal *ListenersDAL) DeleteListener(ctx context.Context, lid string) error {
	query := fmt.Sprintf(`DELETE FROM %s.listeners WHERE id = $1`, dal.schema)
	_, err := dal.db.ExecContext(ctx, query, lid)
	if err != nil {
		logger.Error("Unable to Delete listener: ", err)
	}
	return nil
}

func (dal *ListenersDAL) UpdateListener(ctx context.Context, lid string, updates map[string]any) error {

	// setClauses dynamically builds the SET part of the UPDATE query.
	setClauses := []string{}

	// args stores the actual values for the placeholders like $1, $2, etc.
	args := []any{}

	// count is used to track and increment the placeholder index.
	count := 1

	// Iterates over the updates to dynamically build the SET clause and collect the values for the placeholders in args.
	for field, value := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", field, count))
		args = append(args, value)
		count++
	}

	query := fmt.Sprintf(`UPDATE %s.listeners
      SET %s
      WHERE id = $%d`, dal.schema, strings.Join(setClauses, ", "), count)
	args = append(args, lid)

	updated, err := dal.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Error("Failed to update listener\n", err)
		return fmt.Errorf("Failed to update listener")
	}

	// checks how many rows are affected or not.
	rowsAffected, err := updated.RowsAffected()
	if err != nil {
		logger.Error("Failed to retrieve affected rows\n", err)
		return fmt.Errorf("Failed to retrieve affected rows")
	}
	if rowsAffected == 0 {
		logger.Warn("No rows were updated.")
		return fmt.Errorf("No listener found with id: %s", lid)
	}

	logger.Debug("Rows Affected:", rowsAffected)
	return nil
}
