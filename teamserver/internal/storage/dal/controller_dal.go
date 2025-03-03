package dal

import (
	"context"
	"database/sql"

	"github.com/ksel172/Meduza/teamserver/models"
)

type IControllerDAL interface {
	RegisterController(ctx context.Context, registration models.ControllerRegistration) error

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

func (dal *ControllerDAL) RegisterController(ctx context.Context, registration models.ControllerRegistration) error {
	return nil
}

func (dal *ControllerDAL) ControllerExists(ctx context.Context, controllerID string) (bool, error) {
	return false, nil
}

func (dal *ControllerDAL) UpdateHeartbeat(ctx context.Context, controllerID string, heartbeat models.HeartbeatRequest) error {
	return nil
}
