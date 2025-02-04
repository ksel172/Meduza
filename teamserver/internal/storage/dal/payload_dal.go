package dal

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ksel172/Meduza/teamserver/models"
)

type IPayloadDAL interface {
	CreatePayload(ctx context.Context, config models.PayloadConfig) error
	GetAllPayloads(ctx context.Context) ([]models.PayloadConfig, error)
	DeletePayload(ctx context.Context, payloadID string) error
	DeleteAllPayloads(ctx context.Context) error
	GetKeys(ctx context.Context, configID string) ([]byte, []byte, error)
	GetToken(ctx context.Context, configID string) (string, error)
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
	tx, err := dal.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	listenerConfigJSON, err := json.Marshal(config.ListenerConfig)
	if err != nil {
		return fmt.Errorf("failed to marshal listener config to JSON: %w", err)
	}

	query := fmt.Sprintf(`INSERT INTO %s.payloads (
		payload_id, payload_name, config_id, listener_id, private_key, public_key, payload_token, arch,
		listener_config, sleep, jitter, start_date, kill_date, working_hours_start, working_hours_end)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`, dal.schema)
	_, err = tx.ExecContext(ctx, query, config.PayloadID, config.PayloadName, config.ConfigID,
		config.ListenerID, config.PrivateKey, config.PublicKey, config.Token, config.Arch, listenerConfigJSON,
		config.Sleep, config.Jitter, config.StartDate, config.KillDate, config.WorkingHoursStart,
		config.WorkingHoursEnd)
	if err != nil {
		return fmt.Errorf("failed to insert payload: %w", err)
	}

	return tx.Commit()
}

func (dal *PayloadDAL) GetAllPayloads(ctx context.Context) ([]models.PayloadConfig, error) {
	query := fmt.Sprintf(`SELECT payload_id, payload_name, config_id, listener_id, arch, listener_config, sleep, jitter, start_date, kill_date, working_hours_start, working_hours_end
        FROM %s.payloads`, dal.schema)

	rows, err := dal.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all payloads: %w", err)
	}
	defer rows.Close()

	var configs []models.PayloadConfig
	for rows.Next() {
		var config models.PayloadConfig
		err := rows.Scan(&config.PayloadID, &config.PayloadName, &config.ConfigID, &config.ListenerID, &config.Arch, &config.ListenerConfig, &config.Sleep, &config.Jitter, &config.StartDate, &config.KillDate, &config.WorkingHoursStart, &config.WorkingHoursEnd)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payload: %w", err)
		}
		configs = append(configs, config)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return configs, nil
}

func (dal *PayloadDAL) DeletePayload(ctx context.Context, payloadID string) error {
	query := fmt.Sprintf(`DELETE FROM %s.payloads WHERE payload_id = $1`, dal.schema)
	_, err := dal.db.ExecContext(ctx, query, payloadID)
	if err != nil {
		return fmt.Errorf("failed to delete payload: %w", err)
	}
	return nil
}

func (dal *PayloadDAL) DeleteAllPayloads(ctx context.Context) error {
	query := fmt.Sprintf(`DELETE FROM %s.payloads`, dal.schema)
	_, err := dal.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete all payloads: %w", err)
	}
	return nil
}

func (dal *PayloadDAL) GetKeys(ctx context.Context, configID string) ([]byte, []byte, error) {
	query := fmt.Sprintf(`
		SELECT private_key, public_key
		FROM %s.payloads
		WHERE config_id = $1`,
		dal.schema)
	stmt, err := dal.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return nil, nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	var publicKey []byte
	var privateKey []byte
	if err := stmt.QueryRowContext(ctx, configID).Scan(&privateKey, &publicKey); err != nil {
		log.Printf("Failed to scan public key: %v", err)
		return nil, nil, fmt.Errorf("failed to scan public key: %w", err)
	}

	return privateKey, publicKey, nil
}

func (dal *PayloadDAL) GetToken(ctx context.Context, configID string) (string, error) {
	query := fmt.Sprintf(`
		SELECT payload_token
		FROM %s.payloads
		WHERE config_id = $1`,
		dal.schema)
	stmt, err := dal.db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return "", fmt.Errorf("failed to prepare statement: %w", err)
	}

	var payloadToken string
	if err := stmt.QueryRowContext(ctx, configID).Scan(&payloadToken); err != nil {
		log.Printf("Failed to scan payload token: %v", err)
		return "", fmt.Errorf("failed to scan payload token: %w", err)
	}

	return payloadToken, nil
}
