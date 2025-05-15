package dal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type IPayloadDAL interface {
	CreatePayload(ctx context.Context, payload models.PayloadConfig) error
	GetPayloadByToken(ctx context.Context, payloadToken string) (models.PayloadConfig, error)
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

func (dal *PayloadDAL) CreatePayload(ctx context.Context, payload models.PayloadConfig) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.payloads 
			id, listener_id, config_id, name, arch, public_key, private_key, token
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx, payload.ID, payload.ListenerID, payload.ConfigID, payload.Name,
			payload.Arch, payload.PrivateKey, payload.PublicKey, payload.Token)
		if err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to create payload: %v", err))
			return fmt.Errorf("failed to create payload: %w", err)
		}

		return nil
	})
}

// Agent only know the payload token, not its ID, this retrieves the payload using the token
func (dal *PayloadDAL) GetPayloadByToken(ctx context.Context, payloadToken string) (models.PayloadConfig, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, listener_id, config_id, name, arch, created_at
		FROM %s.payloads
		WHERE token = $1`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) (models.PayloadConfig, error) {
		var payload models.PayloadConfig
		if err := stmt.QueryRowContext(ctx).Scan(&payload.ID, &payload.ListenerID, &payload.ConfigID,
			&payload.Name, &payload.Arch, &payload.CreatedAt,
		); err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to scan payload: %v", err))
			return models.PayloadConfig{}, fmt.Errorf("failed to scan payload: %w", err)
		}
		return payload, nil
	})
}

func (dal *PayloadDAL) GetAllPayloads(ctx context.Context) ([]models.PayloadConfig, error) {
	query := fmt.Sprintf(`
		SELECT 
			id, listener_id, config_id, name, arch, created_at
		FROM %s.payloads`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) ([]models.PayloadConfig, error) {
		rows, err := stmt.QueryContext(ctx)
		if err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to get all payloads: %v", err))
			return nil, fmt.Errorf("failed to get all payloads: %w", err)
		}
		defer rows.Close()

		var payloads []models.PayloadConfig
		for rows.Next() {
			var payload models.PayloadConfig
			if err := rows.Scan(&payload.ID, &payload.ListenerID, &payload.ConfigID, &payload.Name,
				&payload.Arch, &payload.CreatedAt
			); err != nil {
				logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to scan payload: %v", err))
				return nil, fmt.Errorf("failed to scan payload: %w", err)
			}
			payloads = append(payloads, payload)
		}

		if err := rows.Err(); err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("rows iteration error: %v", err))
			return nil, fmt.Errorf("rows iteration error: %w", err)
		}

		return payloads, nil
	})
}

func (dal *PayloadDAL) DeletePayload(ctx context.Context, payloadID string) error {
	query := fmt.Sprintf(`DELETE FROM %s.payloads WHERE id = $1`, dal.schema)

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
		SELECT private_key, public_key
		FROM %s.payloads
		WHERE token = $1`, dal.schema)

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
		SELECT token
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
