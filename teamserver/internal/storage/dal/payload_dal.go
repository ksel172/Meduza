package dal

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

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

type IPayloadDAL interface {
	CreatePayloadManifest(ctx context.Context, payload *models.PayloadManifestV1) error
	GetPayloadManifest(ctx context.Context, payloadID string) (*models.PayloadManifestV1, error)
	GetAllPayloadManifests(ctx context.Context) ([]*models.PayloadManifestV1, error)
	DeletePayloadManifest(ctx context.Context, payloadID string) error
	DeleteAllPayloadManifests(ctx context.Context) error
}

func (dal *PayloadDAL) CreatePayloadManifest(ctx context.Context, payload *models.PayloadManifestV1) error {
	query := fmt.Sprintf(`
        INSERT INTO %s.payload_manifests (
            manifest_id, body, created_at
        ) VALUES ($1, $2, $3)`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		// Convert entire payload to JSON blob
		manifestJSON, err := json.Marshal(payload)
		if err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to marshal payload manifest: %v", err))
			return fmt.Errorf("failed to marshal payload manifest: %w", err)
		}

		// Execute insert query with the manifest as a JSON blob
		_, err = stmt.ExecContext(ctx,
			payload.ID,
			manifestJSON,
			time.Now())

		if err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to create payload manifest: %v", err))
			return fmt.Errorf("failed to create payload manifest: %w", err)
		}

		return nil
	})
}

func (dal *PayloadDAL) GetPayloadManifest(ctx context.Context, payloadID string) (*models.PayloadManifestV1, error) {
	query := fmt.Sprintf(`
        SELECT body FROM %s.payload_manifests WHERE manifest_id = $1`, dal.schema)

	var manifestJSON []byte
	err := utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		row := stmt.QueryRowContext(ctx, payloadID)
		return row.Scan(&manifestJSON)
	})

	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("payload manifest not found: %s", payloadID))
			return nil, fmt.Errorf("payload manifest not found: %s", payloadID)
		}
		logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to get payload manifest: %v", err))
		return nil, fmt.Errorf("failed to get payload manifest: %w", err)
	}

	var manifest models.PayloadManifestV1
	if err := json.Unmarshal(manifestJSON, &manifest); err != nil {
		logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to unmarshal payload manifest: %v", err))
		return nil, fmt.Errorf("failed to unmarshal payload manifest: %w", err)
	}

	return &manifest, nil
}

func (dal *PayloadDAL) GetAllPayloadManifests(ctx context.Context) ([]*models.PayloadManifestV1, error) {
	query := fmt.Sprintf(`
        SELECT body FROM %s.payload_manifests`, dal.schema)

	var manifests []*models.PayloadManifestV1

	err := utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		rows, err := stmt.QueryContext(ctx)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var manifestJSON []byte
			if err := rows.Scan(&manifestJSON); err != nil {
				return err
			}

			var manifest models.PayloadManifestV1
			if err := json.Unmarshal(manifestJSON, &manifest); err != nil {
				return err
			}

			manifests = append(manifests, &manifest)
		}

		return rows.Err()
	})

	if err != nil {
		logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to get all payload manifests: %v", err))
		return nil, fmt.Errorf("failed to get all payload manifests: %w", err)
	}

	return manifests, nil
}

func (dal *PayloadDAL) DeletePayloadManifest(ctx context.Context, payloadID string) error {
	query := fmt.Sprintf(`
        DELETE FROM %s.payload_manifests WHERE manifest_id = $1`, dal.schema)

	err := utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		result, err := stmt.ExecContext(ctx, payloadID)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return fmt.Errorf("payload manifest not found: %s", payloadID)
		}

		return nil
	})

	if err != nil {
		logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to delete payload manifest: %v", err))
		return fmt.Errorf("failed to delete payload manifest: %w", err)
	}

	return nil
}

func (dal *PayloadDAL) DeleteAllPayloadManifests(ctx context.Context) error {
	query := fmt.Sprintf(`
        DELETE FROM %s.payload_manifests`, dal.schema)

	err := utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx)
		return err
	})

	if err != nil {
		logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to delete all payload manifests: %v", err))
		return fmt.Errorf("failed to delete all payload manifests: %w", err)
	}

	return nil
}
