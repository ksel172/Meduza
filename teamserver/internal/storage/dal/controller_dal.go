package dal

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type IControllerDAL interface {
	RegisterController(ctx context.Context, registration models.ControllerRegistration) error
	UpdateHeartbeat(ctx context.Context, controllerID string, heartbeat models.HeartbeatRequest) error
}

type ControllerDAL struct {
	db     *sql.DB
	schema string
}

func NewControllerDAL(db *sql.DB, schema string) IControllerDAL {
	return &ControllerDAL{db: db, schema: schema}
}

func (dal *ControllerDAL) RegisterController(ctx context.Context, controller models.ControllerRegistration) error {
	query := fmt.Sprintf(`
        INSERT INTO %s.controllers (
            id, endpoint
        ) VALUES($1, $2)
    `, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		logger.Debug(logLevel, logDetailController, fmt.Sprintf("Registering controller: %s", controller.ID))

		_, err := stmt.ExecContext(
			ctx,
			controller.ID,
			controller.Endpoint,
		)

		if err != nil {
			logger.Error(logLevel, logDetailController, "Failed to register controller: ", err)
			return fmt.Errorf("failed to register controller: %w", err)
		}
		return nil
	})
}

func (dal *ControllerDAL) UpdateHeartbeat(ctx context.Context, controllerID string, heartbeat models.HeartbeatRequest) error {

	// Map of status to UUIDs with the same status
	statuses := map[string][]string{}
	for id, status := range heartbeat.Listeners {
		statuses[status] = append(statuses[status], id)
	}

	selectQuery := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s.controllers WHERE id = $1)`, dal.schema)
	updateHeartbeatQuery := fmt.Sprintf(`UPDATE %s.controllers SET heartbeat = CURRENT_TIMESTAMP WHERE id = $1`, dal.schema)

	// updateListenersQuery := fmt.Sprintf(`UPDATE %s.listeners SET status = $1 WHERE id IN $2`, dal.schema)

	return utils.WithTransactionTimeout(ctx, dal.db, 5, sql.TxOptions{}, func(ctx context.Context, tx *sql.Tx) error {
		logger.Debug(logLevel, logDetailController, fmt.Sprintf("Updating heartbeat for controller: %s", controllerID))

		// Check if the controller exists by its ID
		var exists bool
		if err := tx.QueryRowContext(ctx, selectQuery, controllerID).Scan(&exists); err != nil {
			logger.Error(logLevel, logDetailController, fmt.Sprintf("Failed to execute select controller query: %v", err))
			return fmt.Errorf("failed to execute select controller query: %w", err)
		}
		if !exists {
			logger.Info(logLevel, logDetailController, fmt.Sprintf("Controller with ID %s does not exist", controllerID))
			return fmt.Errorf("controller with ID %s does not exist", controllerID)
		}

		// Update the listener status rows on the listeners table
		for status, ids := range statuses {
			if len(ids) == 0 {
				continue
			}

			// Build query with right number of placeholders
			placeholders := make([]string, len(ids))
			args := make([]any, len(ids)+1)
			args[0] = status

			for i, id := range ids {
				placeholders[i] = fmt.Sprintf("$%d", i+2) // +2 because status is $1
				args[i+1] = id
			}

			// Prepare statement
			updateListenersQuery := fmt.Sprintf(
				`UPDATE %s.listeners SET STATUS = $1 WHERE id IN (%s)`,
				dal.schema,
				strings.Join(placeholders, ","),
			)
			updateListenersStmt, err := dal.db.PrepareContext(ctx, updateListenersQuery)
			if err != nil {
				logger.Error(logLevel, logDetailController, fmt.Sprintf("Failed to prepare update listener query: %v", err))
				return fmt.Errorf("failed to prepare update listener query: %w", err)
			}

			// Execute with statement in transaction
			result, err := tx.StmtContext(ctx, updateListenersStmt).ExecContext(ctx, args...)
			if err != nil {
				logger.Error(logLevel, logDetailController, fmt.Sprintf("Failed to update listeners statuses: %v", err))
				return fmt.Errorf("failed to update listeners statuses: %w", err)
			}

			rowsAffected, err := result.RowsAffected()
			if err == nil && rowsAffected == 0 {
				logger.Warn(logLevel, logDetailController, fmt.Sprintf("no listeners with status %s were updated", status))
			}
		}

		// Finally, update the heartbeat on the controllers table
		_, err := tx.ExecContext(ctx, updateHeartbeatQuery, controllerID)
		if err != nil {
			logger.Error(logLevel, logDetailController, "Failed to update heartbeat: ", err)
			return fmt.Errorf("failed to update heartbeat: %w", err)
		}

		return nil
	})
}
