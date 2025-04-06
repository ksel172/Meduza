package dal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type IControllerDal interface {
	RegisterExternalController()
}

type ControllerDal struct {
	db     *sql.DB
	schema string
}

func NewControllerDal(db *sql.DB, schema string) *ControllerDal {
	return &ControllerDal{
		db:     db,
		schema: schema,
	}
}

func (dal *ControllerDal) RegisterExternalController(ctx context.Context, controller models.ListenerController) error {

	query := fmt.Sprintf(`
	INSERT INTO %s.external_controllers (callback_ip, callback_port, listener_config, created_at, last_seen) 
	VALUES ($1, $2, $3, $4, $5)`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx, controller.CallbackIP, controller.CallbackPort, controller.ListenerConfig, controller.CreatedAt, controller.LastSeen)
		if err != nil {
			logger.Error(logLevel, logDetailController, "Failed to create listener: ", err)
			return fmt.Errorf("failed to register external controller: %w", err)
		}
		return nil
	})
}

func (dal *ControllerDal) CheckInController(ctx context.Context, controllerID string) error {
	query := fmt.Sprintf(`
	UPDATE %s.external_controllers 
	SET last_seen = NOW() 
	WHERE id = $1`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx, controllerID)
		if err != nil {
			logger.Error(logLevel, logDetailController, "Failed to update controller check-in: ", err)
			return fmt.Errorf("failed to check in controller: %w", err)
		}
		return nil
	})
}

func (dal *ControllerDal) QueryByKind(ctx context.Context) error {

	query := fmt.Sprintf(`
	SELECT id, callback_ip, callback_port, listener_config, created_at, last_seen
	FROM %s.external_controllers
	WHERE listener_kind = $1`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.QueryContext(ctx, "external")
		if err != nil {
			logger.Error(logLevel, logDetailController, "Failed to query external controllers: ", err)
			return fmt.Errorf("failed to query external controllers: %w", err)
		}
		return nil
	})
}
