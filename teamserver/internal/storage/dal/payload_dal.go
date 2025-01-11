package dal

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/models"
)

type IPayloadDAL interface {
	CreatePayload(ctx context.Context, config models.PayloadConfig) error
	GetAllPayloads(ctx context.Context) ([]models.PayloadConfig, error)
	DeletePayload(ctx context.Context, payloadID string) error
	DeleteAllPayloads(ctx context.Context) error
}

type PayloadDAL struct {
	db     *sql.DB
	schema string
}

func NewPayloadDAL(db *sql.DB, schema string) *PayloadDAL {
	return &PayloadDAL{
		db:     db,
		schema: schema,
	}
}
func (dal *PayloadDAL) CreatePayload(ctx context.Context, config models.PayloadConfig) error {
	tx, err := dal.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	listenerConfigJSON, err := json.Marshal(config.ListenerConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal listener config to JSON: %w", err)
	}

	// Insert agent
	payloadQuery := fmt.Sprintf(`
        INSERT INTO %s.payloads (payload_id, payload_name, config_id, listener_id, arch, listener_config, sleep, jitter, start_date, kill_date, working_hours_start, working_hours_end)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`, dal.schema)

	_, err = tx.Exec(payloadQuery, config.PayloadID, config.PayloadName, config.ConfigID, config.ListenerID, config.Arch, listenerConfigJSON, config.Sleep, config.Jitter, config.StartDate, config.KillDate, config.WorkingHoursStart, config.WorkingHoursEnd)
	if err != nil {
		return fmt.Errorf("failed to insert agent: %w", err)
	}

	return tx.Commit()
}

func (dal *PayloadDAL) GetAllPayloads(ctx context.Context) ([]models.PayloadConfig, error) {
	query := fmt.Sprintf(`
		SELECT *
		FROM %s.payloads`, dal.schema)

	rows, err := dal.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all agents: %w", err)
	}
	defer rows.Close()

	var configs []models.PayloadConfig
	for rows.Next() {
		var config models.PayloadConfig
		err := rows.Scan(&config.PayloadID, &config.PayloadName, &config.ConfigID, &config.ListenerID, &config.Arch, &config.ListenerConfig, &config.Sleep, &config.Jitter, &config.StartDate, &config.KillDate, &config.WorkingHoursStart, &config.WorkingHoursEnd)
		if err != nil {
			return nil, fmt.Errorf("failed to scan agent: %w", err)
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func (dal *PayloadDAL) DeletePayload(ctx context.Context, payloadID string) error {
	query := "DELETE FROM payloads WHERE payload_id = $1"
	_, err := dal.db.ExecContext(ctx, query, payloadID)
	if err != nil {
		return fmt.Errorf("failed to delete payload: %w", err)
	}
	return nil
}

func (dal *PayloadDAL) DeleteAllPayloads(ctx context.Context) error {
	query := "DELETE FROM payloads"
	_, err := dal.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete all payloads: %w", err)
	}
	return nil
}
