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

type IControllerDAL interface {
	RegisterController(ctx context.Context, registration models.Controller) error
	ControllerExists(ctx context.Context, controllerID string) (bool, error)
	UpdateHeartbeat(ctx context.Context, controllerID string, heartbeat models.HeartbeatRequest) error
}

type ControllerDAL struct {
	db     *sql.DB
	schema string
}

func NewControllerDAL(db *sql.DB, schema string) IControllerDAL {
	return &ControllerDAL{db: db, schema: schema}
}

func (dal *ControllerDAL) RegisterController(ctx context.Context, controller models.Controller) error {
	query := fmt.Sprintf(`
        INSERT INTO %s.controllers (
            id, endpoint, public_key, private_key, created_at, updated_at
        ) VALUES($1, $2, $3, $4, $5, $5)
    `, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		logger.Debug(logLevel, logDetailController, fmt.Sprintf("Registering controller: %s", controller.ID))

		_, err := stmt.ExecContext(
			ctx,
			controller.ID,
			controller.Endpoint,
			controller.PublicKey,
			controller.PrivateKey,
			time.Now().UTC(),
		)

		if err != nil {
			logger.Error(logLevel, logDetailController, "Failed to register controller: ", err)
			return fmt.Errorf("failed to register controller: %w", err)
		}
		return nil
	})
}

func (dal *ControllerDAL) ControllerExists(ctx context.Context, controllerID string) (bool, error) {
	query := fmt.Sprintf(`
        SELECT EXISTS(
            SELECT 1 FROM %s.controllers WHERE id = $1
        )
    `, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) (bool, error) {
		logger.Debug(logLevel, logDetailController, fmt.Sprintf("Checking if controller exists: %s", controllerID))

		var exists bool
		err := stmt.QueryRowContext(ctx, controllerID).Scan(&exists)
		if err != nil {
			logger.Error(logLevel, logDetailController, "Error checking controller existence: ", err)
			return false, fmt.Errorf("failed to check controller existence: %w", err)
		}

		return exists, nil
	})
}

func (dal *ControllerDAL) UpdateHeartbeat(ctx context.Context, controllerID string, heartbeat models.HeartbeatRequest) error {
	// First convert the listeners map to JSON for storage
	listenersJSON, err := json.Marshal(heartbeat.Listeners)
	if err != nil {
		logger.Error(logLevel, logDetailController, "Failed to marshal listeners to JSON: ", err)
		return fmt.Errorf("failed to marshal listeners to JSON: %v", err)
	}

	query := fmt.Sprintf(`
        UPDATE %s.controllers
        SET updated_at = $1,
            heartbeat_timestamp = $2,
            listeners_status = $3
        WHERE id = $4
    `, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		logger.Debug(logLevel, logDetailController, fmt.Sprintf("Updating heartbeat for controller: %s", controllerID))

		_, err := stmt.ExecContext(
			ctx,
			time.Now().UTC(),
			time.Unix(heartbeat.Timestamp, 0).UTC(),
			listenersJSON,
			controllerID,
		)

		if err != nil {
			logger.Error(logLevel, logDetailController, "Failed to update heartbeat: ", err)
			return fmt.Errorf("failed to update heartbeat: %w", err)
		}
		return nil
	})
}
