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

	// Build job methods
	CreateBuildJob(ctx context.Context, job *models.PayloadJob) error
	UpdateBuildJob(ctx context.Context, job *models.PayloadJob) error
	GetBuildJob(ctx context.Context, jobID string) (*models.PayloadJob, error)
	GetPayloadBuildJobs(ctx context.Context, payloadID string) ([]*models.PayloadJob, error)
	DeleteBuildJob(ctx context.Context, jobID string) error
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
