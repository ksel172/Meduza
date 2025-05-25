package dal

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type IPayloadDAL interface {
	CreatePayloadManifest(ctx context.Context, payload *models.PayloadManifestV1) error
	GetPayloadManifest(ctx context.Context, payloadID string) (*models.PayloadManifestV1, error)
	GetAllPayloadManifests(ctx context.Context) ([]*models.PayloadManifestV1, error)
	DeletePayloadManifest(ctx context.Context, payloadID string) error
	DeleteAllPayloadManifests(ctx context.Context) error

	// Build job methods
	CreateBuildJob(ctx context.Context, job *models.PayloadJob) error
	UpdateBuildJob(ctx context.Context, job *models.PayloadJob) error
	GetBuildJob(ctx context.Context, jobID string) (*models.PayloadJob, error)
	GetPayloadBuildJobs(ctx context.Context, payloadID string) ([]*models.PayloadJob, error)
	DeleteBuildJob(ctx context.Context, jobID string) error

	// Payload methods
	CreatePayload(ctx context.Context, payload *models.Payload) error
	GetPayload(ctx context.Context, payloadID string) (*models.Payload, error)
	GetPayloads(ctx context.Context) ([]*models.Payload, error)
	DeletePayload(ctx context.Context, payloadID string) error
	GetPayloadByToken(ctx context.Context, payloadToken string) (models.Payload, error)
	GetKeys(ctx context.Context, authToken string) ([]byte, []byte, error)
	GetToken(ctx context.Context, configID string) (string, error)
	
	// Download payload build
	DownloadPayloadBuild(ctx context.Context, jobID string) ([]byte, string, error)
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

// <<<<<<< dev
// func (dal *PayloadDAL) CreatePayload(ctx context.Context, payload models.PayloadConfig) error {
// 	query := fmt.Sprintf(`
// 		INSERT INTO %s.payloads
// 			(id, listener_id, config_id, name, arch, public_key, private_key, token)
//         VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, dal.schema)

//	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
//		_, err := stmt.ExecContext(ctx, payload.ID, payload.ListenerID, payload.ConfigID, payload.Name,
//			payload.Arch, payload.PublicKey, payload.PrivateKey, payload.Token)
//
// =======

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

func (dal *PayloadDAL) GetPayloadByToken(ctx context.Context, payloadToken string) (models.Payload, error) {
	query := fmt.Sprintf(`
		SELECT
			id, listener_id, config_id, name, arch, created_at
		FROM %s.payloads
		WHERE token = $1`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) (models.Payload, error) {
		var payload models.Payload
		if err := stmt.QueryRowContext(ctx, payloadToken).Scan(&payload.ID, &payload.ListenerID, &payload.ConfigID,
			&payload.Name, &payload.Arch, &payload.CreatedAt,
		); err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to scan payload: %v", err))
			return models.Payload{}, fmt.Errorf("failed to scan payload: %w", err)
		}
		return payload, nil
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

		// <<<<<<< dev
		// 		var payloads []models.PayloadConfig
		// 		for rows.Next() {
		// 			var payload models.PayloadConfig
		// 			if err := rows.Scan(&payload.ID, &payload.ListenerID, &payload.ConfigID, &payload.Name,
		// 				&payload.Arch, &payload.CreatedAt,
		// 			); err != nil {
		// 				logger.Error(logLevel, logDetailPayload, fmt.Sprintf("failed to scan payload: %v", err))
		// 				return nil, fmt.Errorf("failed to scan payload: %w", err)
		// 			}
		// 			payloads = append(payloads, payload)
		// =======
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

		// <<<<<<< dev
		// 		return payloads, nil
		// =======
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

// <<<<<<< dev
// func (dal *PayloadDAL) DeletePayload(ctx context.Context, payloadID string) error {
// 	query := fmt.Sprintf(`DELETE FROM %s.payloads WHERE id = $1`, dal.schema)
// =======
// Build Job Methods

func (dal *PayloadDAL) CreateBuildJob(ctx context.Context, job *models.PayloadJob) error {
	query := fmt.Sprintf(`
        INSERT INTO %s.payload_build_jobs (
            job_id, payload_id, status, architecture, parameters, 
            start_time, end_time, output_path, error_message, build_log
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		// Convert parameters to JSON
		paramsJSON, err := json.Marshal(job.Parameters)
		if err != nil {
			return fmt.Errorf("failed to marshal parameters: %w", err)
		}

		_, err = stmt.ExecContext(ctx,
			job.ID,
			job.PayloadID,
			job.Status,
			job.Architecture,
			paramsJSON,
			job.StartTime,
			sql.NullTime{Time: job.EndTime, Valid: !job.EndTime.IsZero()},
			job.OutputPath,
			job.ErrorMessage,
			job.BuildLog)

		if err != nil {
			return fmt.Errorf("failed to create build job: %w", err)
		}

		return nil
	})
}

func (dal *PayloadDAL) UpdateBuildJob(ctx context.Context, job *models.PayloadJob) error {
	query := fmt.Sprintf(`
        UPDATE %s.payload_build_jobs SET
            status = $1,
            end_time = $2,
            output_path = $3,
            error_message = $4,
            build_log = $5,
            updated_at = CURRENT_TIMESTAMP
        WHERE job_id = $6`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		result, err := stmt.ExecContext(ctx,
			job.Status,
			sql.NullTime{Time: job.EndTime, Valid: !job.EndTime.IsZero()},
			job.OutputPath,
			job.ErrorMessage,
			job.BuildLog,
			job.ID)

		if err != nil {
			return fmt.Errorf("failed to update build job: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return fmt.Errorf("build job not found: %s", job.ID)
		}

		return nil
	})
}

func (dal *PayloadDAL) GetBuildJob(ctx context.Context, jobID string) (*models.PayloadJob, error) {
	query := fmt.Sprintf(`
        SELECT job_id, payload_id, status, architecture, parameters, 
               start_time, end_time, output_path, error_message, build_log
        FROM %s.payload_build_jobs
        WHERE job_id = $1`, dal.schema)

	var job models.PayloadJob
	var paramsJSON []byte
	var endTime sql.NullTime

	err := utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		row := stmt.QueryRowContext(ctx, jobID)
		return row.Scan(
			&job.ID,
			&job.PayloadID,
			&job.Status,
			&job.Architecture,
			&paramsJSON,
			&job.StartTime,
			&endTime,
			&job.OutputPath,
			&job.ErrorMessage,
			&job.BuildLog)
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("build job not found: %s", jobID)
		}
		return nil, fmt.Errorf("failed to get build job: %w", err)
	}

	// Parse parameters JSON
	if err := json.Unmarshal(paramsJSON, &job.Parameters); err != nil {
		return nil, fmt.Errorf("failed to unmarshal parameters: %w", err)
	}

	// Handle nullable end_time
	if endTime.Valid {
		job.EndTime = endTime.Time
	}

	return &job, nil
}

func (dal *PayloadDAL) GetPayloadBuildJobs(ctx context.Context, payloadID string) ([]*models.PayloadJob, error) {
	query := fmt.Sprintf(`

        SELECT job_id, payload_id, status, architecture, parameters, 
               start_time, end_time, output_path, error_message, build_log
        FROM %s.payload_build_jobs
        WHERE payload_id = $1
        ORDER BY start_time DESC`, dal.schema)

	var jobs []*models.PayloadJob

	err := utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		rows, err := stmt.QueryContext(ctx, payloadID)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var job models.PayloadJob
			var paramsJSON []byte
			var endTime sql.NullTime

			if err := rows.Scan(
				&job.ID,
				&job.PayloadID,
				&job.Status,
				&job.Architecture,
				&paramsJSON,
				&job.StartTime,
				&endTime,
				&job.OutputPath,
				&job.ErrorMessage,
				&job.BuildLog); err != nil {
				return err
			}

			// Parse parameters JSON
			if err := json.Unmarshal(paramsJSON, &job.Parameters); err != nil {
				return fmt.Errorf("failed to unmarshal parameters: %w", err)
			}

			// Handle nullable end_time
			if endTime.Valid {
				job.EndTime = endTime.Time
			}

			jobs = append(jobs, &job)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get build jobs: %w", err)
	}

	return jobs, nil
}

func (dal *PayloadDAL) DeleteBuildJob(ctx context.Context, jobID string) error {
	query := fmt.Sprintf(`
        DELETE FROM %s.payload_build_jobs WHERE job_id = $1`, dal.schema)

	err := utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		result, err := stmt.ExecContext(ctx, jobID)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}

		if rowsAffected == 0 {
			return fmt.Errorf("build job not found: %s", jobID)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to delete build job: %w", err)
	}

	return nil
}

// Payload methods

func (dal *PayloadDAL) CreatePayload(ctx context.Context, payload *models.Payload) error {
	query := fmt.Sprintf(`
        INSERT INTO %s.payloads (
            payload_id, listener_id, config_id, manifest_id, name, 
            architecture, public_key, private_key, token, created_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx,
			payload.ID,
			payload.ListenerID,
			payload.ConfigID,
			payload.ManifestID,
			payload.Name,
			payload.Arch,
			payload.PublicKey,
			payload.PrivateKey,
			payload.Token,
			payload.CreatedAt)

		if err != nil {
			return fmt.Errorf("failed to create payload: %w", err)
		}

		return nil
	})
}

func (dal *PayloadDAL) GetPayload(ctx context.Context, payloadID string) (*models.Payload, error) {
	query := fmt.Sprintf(`
        SELECT payload_id, listener_id, config_id, manifest_id, name, 
               architecture, public_key, private_key, token, created_at
        FROM %s.payloads
        WHERE payload_id = $1`, dal.schema)

	var payload models.Payload

	err := utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		row := stmt.QueryRowContext(ctx, payloadID)
		return row.Scan(
			&payload.ID,
			&payload.ListenerID,
			&payload.ConfigID,
			&payload.ManifestID,
			&payload.Name,
			&payload.Arch,
			&payload.PublicKey,
			&payload.PrivateKey,
			&payload.Token,
			&payload.CreatedAt)
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("payload not found: %s", payloadID)
		}
		return nil, fmt.Errorf("failed to get payload: %w", err)
	}

	return &payload, nil
}

func (dal *PayloadDAL) GetPayloads(ctx context.Context) ([]*models.Payload, error) {
	query := fmt.Sprintf(`
        SELECT payload_id, listener_id, config_id, manifest_id, name, 
               architecture, public_key, private_key, token, created_at
        FROM %s.payloads
        ORDER BY created_at DESC`, dal.schema)

	var payloads []*models.Payload

	err := utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		rows, err := stmt.QueryContext(ctx)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var payload models.Payload
			if err := rows.Scan(
				&payload.ID,
				&payload.ListenerID,
				&payload.ConfigID,
				&payload.ManifestID,
				&payload.Name,
				&payload.Arch,
				&payload.PublicKey,
				&payload.PrivateKey,
				&payload.Token,
				&payload.CreatedAt); err != nil {
				return err
			}

			payloads = append(payloads, &payload)
		}

		return rows.Err()
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get payloads: %w", err)
	}

	return payloads, nil
}

func (dal *PayloadDAL) DeletePayload(ctx context.Context, payloadID string) error {
	query := fmt.Sprintf(`
        DELETE FROM %s.payloads WHERE payload_id = $1`, dal.schema)

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
			return fmt.Errorf("payload not found: %s", payloadID)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to delete payload: %w", err)
	}

	return nil
}

func (dal *PayloadDAL) DownloadPayloadBuild(ctx context.Context, jobID string) ([]byte, string, error) {
	// First, get the build job to find the output path
	job, err := dal.GetBuildJob(ctx, jobID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get build job: %w", err)
	}

	if job.Status != "completed" {
		return nil, "", fmt.Errorf("build job is not completed: %s", job.Status)
	}

	if job.OutputPath == "" {
		return nil, "", fmt.Errorf("build job has no output path")
	}

	// Read the file
	data, err := os.ReadFile(job.OutputPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read build output file: %w", err)
	}

	// Get the filename from the path
	filename := filepath.Base(job.OutputPath)

	return data, filename, nil
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