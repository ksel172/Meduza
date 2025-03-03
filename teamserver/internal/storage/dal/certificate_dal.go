package dal

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

// CertificateDAL interface for certificate operations
type CertificateDAL interface {
	SaveCertificate(ctx context.Context, certType, path, filename string) error
	GetAllCertificates(ctx context.Context) ([]models.Certificate, error)
	GetCertificateByType(ctx context.Context, certType string) (*models.Certificate, error)
	DeleteCertificate(ctx context.Context, id string) error
}

// CertificateDALImpl implements CertificateDAL
type CertificateDALImpl struct {
	db     *sql.DB
	schema string
}

// NewCertificateDAL creates a new certificate DAL
func NewCertificateDAL(db *sql.DB, schema string) CertificateDAL {
	return &CertificateDALImpl{
		db:     db,
		schema: schema,
	}
}

// SaveCertificate saves a certificate to the database
func (d *CertificateDALImpl) SaveCertificate(ctx context.Context, certType, path, filename string) error {
	now := time.Now()
	id := uuid.New().String()

	// Used upsert: https://neon.tech/postgresql/postgresql-tutorial/postgresql-upsert
	query := fmt.Sprintf(`
        INSERT INTO "%s"."certificates" (id, type, path, filename, created_at, updated_at) 
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (type) 
        DO UPDATE SET 
            path = $3,
            filename = $4,
            updated_at = $6
    `, d.schema)

	return utils.WithTimeout(ctx, d.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		logger.Debug(logLevel, logDetailCert, fmt.Sprintf("Upserting certificate of type: %s", certType))
		_, err := stmt.ExecContext(ctx, id, certType, path, filename, now, now)
		if err != nil {
			logger.Error(logLevel, logDetailCert, fmt.Sprintf("Failed to upsert certificate: %v", err))
			return fmt.Errorf("failed to upsert certificate: %w", err)
		}
		return nil
	})
}

// GetAllCertificates retrieves all certificates from the database
func (d *CertificateDALImpl) GetAllCertificates(ctx context.Context) ([]models.Certificate, error) {
	query := fmt.Sprintf(`SELECT id, type, path, filename, created_at, updated_at FROM "%s"."certificates"`, d.schema)

	return utils.WithResultTimeout(ctx, d.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) ([]models.Certificate, error) {
		logger.Debug(logLevel, logDetailCert, "Getting all certificates")
		rows, err := stmt.QueryContext(ctx)
		if err != nil {
			logger.Error(logLevel, logDetailCert, fmt.Sprintf("Failed to get certificates: %v", err))
			return nil, fmt.Errorf("failed to get certificates: %w", err)
		}
		defer rows.Close()

		var certificates []models.Certificate
		for rows.Next() {
			var cert models.Certificate
			err := rows.Scan(&cert.ID, &cert.Type, &cert.Path, &cert.Filename, &cert.CreatedAt, &cert.UpdatedAt)
			if err != nil {
				logger.Error(logLevel, logDetailCert, fmt.Sprintf("Failed to scan certificate row: %v", err))
				return nil, fmt.Errorf("failed to scan certificate row: %w", err)
			}
			certificates = append(certificates, cert)
		}

		if err = rows.Err(); err != nil {
			logger.Error(logLevel, logDetailCert, fmt.Sprintf("Error after scanning certificates: %v", err))
			return nil, fmt.Errorf("error after scanning certificates: %w", err)
		}

		return certificates, nil
	})
}

// GetCertificateByType retrieves a certificate by its type
func (d *CertificateDALImpl) GetCertificateByType(ctx context.Context, certType string) (*models.Certificate, error) {
	query := fmt.Sprintf(`SELECT id, type, path, filename, created_at, updated_at FROM "%s"."certificates" WHERE type = $1`, d.schema)

	return utils.WithResultTimeout(ctx, d.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) (*models.Certificate, error) {
		logger.Debug(logLevel, logDetailCert, fmt.Sprintf("Getting certificate by type: %s", certType))

		var cert models.Certificate
		err := stmt.QueryRowContext(ctx, certType).Scan(&cert.ID, &cert.Type, &cert.Path, &cert.Filename, &cert.CreatedAt, &cert.UpdatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				logger.Debug(logLevel, logDetailCert, fmt.Sprintf("Certificate of type %s not found", certType))
				return nil, nil
			}
			logger.Error(logLevel, logDetailCert, fmt.Sprintf("Failed to get certificate by type: %v", err))
			return nil, fmt.Errorf("failed to get certificate by type: %w", err)
		}

		return &cert, nil
	})
}

// DeleteCertificate deletes a certificate from the database
func (d *CertificateDALImpl) DeleteCertificate(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM "%s"."certificates" WHERE id = $1`, d.schema)

	return utils.WithTimeout(ctx, d.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		logger.Debug(logLevel, logDetailCert, fmt.Sprintf("Deleting certificate with ID: %s", id))

		result, err := stmt.ExecContext(ctx, id)
		if err != nil {
			logger.Error(logLevel, logDetailCert, fmt.Sprintf("Failed to delete certificate: %v", err))
			return fmt.Errorf("failed to delete certificate: %w", err)
		}

		rows, err := result.RowsAffected()
		if err != nil {
			logger.Error(logLevel, logDetailCert, fmt.Sprintf("Failed to get rows affected: %v", err))
			return fmt.Errorf("failed to get rows affected: %w", err)
		}

		if rows == 0 {
			logger.Warn(logLevel, logDetailCert, "Certificate not found")
			return fmt.Errorf("certificate not found")
		}

		return nil
	})
}
