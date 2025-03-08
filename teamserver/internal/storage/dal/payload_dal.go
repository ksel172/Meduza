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

type IPayloadDAL interface {
	CreatePayload(ctx context.Context, config models.PayloadConfig) error
	GetAllPayloads(ctx context.Context) ([]models.PayloadConfig, error)
	DeletePayload(ctx context.Context, payloadID string) error
	DeleteAllPayloads(ctx context.Context) error
	GetKeys(ctx context.Context, authToken string) ([]byte, []byte, error)
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
	query := fmt.Sprintf(`INSERT INTO %s.payloads (
		payload_id, payload_name, config_id, listener_id, private_key, public_key, payload_token, arch,
		listener_config, sleep, jitter, start_date, kill_date, working_hours_start, working_hours_end, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		listenerConfigJSON, err := json.Marshal(config.ListenerConfig)
		if err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to marshal listener config: %v", err))
			return fmt.Errorf("failed to marshal listener config to JSON: %w", err)
		}

		_, err = stmt.ExecContext(ctx, config.PayloadID, config.PayloadName, config.ConfigID,
			config.ListenerID, config.PrivateKey, config.PublicKey, config.Token, config.Arch, listenerConfigJSON,
			config.Sleep, config.Jitter, config.StartDate, config.KillDate, config.WorkingHoursStart,
			config.WorkingHoursEnd, config.CreatedAt)
		if err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to create payload: %v", err))
			return fmt.Errorf("failed to create payload: %w", err)
		}

		return nil
	})
}

func (dal *PayloadDAL) GetAllPayloads(ctx context.Context) ([]models.PayloadConfig, error) {
	query := fmt.Sprintf(`SELECT payload_id, payload_name, config_id, listener_id, arch,
		listener_config, sleep, jitter, start_date, kill_date, working_hours_start, 
		working_hours_end, created_at FROM %s.payloads`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) ([]models.PayloadConfig, error) {
		rows, err := stmt.QueryContext(ctx)
		if err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to get all payloads: %v", err))
			return nil, fmt.Errorf("failed to get all payloads: %w", err)
		}
		defer rows.Close()

		var configs []models.PayloadConfig
		for rows.Next() {
			var config models.PayloadConfig
			err := rows.Scan(&config.PayloadID, &config.PayloadName, &config.ConfigID, &config.ListenerID, &config.Arch, &config.ListenerConfig, &config.Sleep, &config.Jitter, &config.StartDate, &config.KillDate, &config.WorkingHoursStart, &config.WorkingHoursEnd, &config.CreatedAt)
			if err != nil {
				logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to scan payload: %v", err))
				return nil, fmt.Errorf("failed to scan payload: %w", err)
			}
			configs = append(configs, config)
		}

		if err := rows.Err(); err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("rows iteration error: %v", err))
			return nil, fmt.Errorf("rows iteration error: %w", err)
		}

		return configs, nil
	})
}

func (dal *PayloadDAL) DeletePayload(ctx context.Context, payloadID string) error {
	query := fmt.Sprintf(`DELETE FROM %s.payloads WHERE payload_id = $1`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx, payloadID)
		if err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to delete payload: %v", err))
			return fmt.Errorf("failed to delete payload: %w", err)
		}
		return nil
	})
}

func (dal *PayloadDAL) DeleteAllPayloads(ctx context.Context) error {
	query := fmt.Sprintf(`DELETE FROM %s.payloads`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx)
		if err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to delete all payloads: %v", err))
			return fmt.Errorf("failed to delete all payloads: %w", err)
		}
		return nil
	})
}

func (dal *PayloadDAL) GetKeys(ctx context.Context, authToken string) ([]byte, []byte, error) {
	query := fmt.Sprintf(`
		SELECT private_key, public_key FROM %s.payloads
		WHERE payload_token = $1`, dal.schema)

	stmt, err := dal.db.PrepareContext(ctx, query)
	if err != nil {
		logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to prepare statement: %v", err))
		return nil, nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	var publicKey []byte
	var privateKey []byte
	if err := stmt.QueryRowContext(ctx, authToken).Scan(&privateKey, &publicKey); err != nil {
		logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to scan keys: %v", err))
		return nil, nil, fmt.Errorf("failed to scan keys: %w", err)
	}

	return privateKey, publicKey, nil
}

func (dal *PayloadDAL) GetToken(ctx context.Context, configID string) (string, error) {
	query := fmt.Sprintf(`
		SELECT payload_token
		FROM %s.payloads
		WHERE config_id = $1`,
		dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) (string, error) {
		var payloadToken string
		if err := stmt.QueryRowContext(ctx, configID).Scan(&payloadToken); err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to get payload token for configID '%s': %v", configID, err))
			return "", fmt.Errorf("failed to get payload token: %w", err)
		}

		return payloadToken, nil
	})
}
