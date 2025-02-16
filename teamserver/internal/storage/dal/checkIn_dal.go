package dal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type ICheckInDAL interface {
	CreateAgent(context.Context, models.Agent) error
}

type CheckInDAL struct {
	db     *sql.DB
	schema string
}

func NewCheckInDAL(db *sql.DB, schema string) *CheckInDAL {
	return &CheckInDAL{db: db, schema: schema}
}

func (dal *CheckInDAL) CreateAgent(ctx context.Context, agent models.Agent) error {

	query := fmt.Sprintf(`
		INSERT INTO %s.agents (id, config_id, name, note, status, first_callback, last_callback, modified_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx, query, agent.AgentID, agent.ConfigID, agent.Name, agent.Note, agent.Status,
			agent.FirstCallback, agent.LastCallback, agent.ModifiedAt)
		if err != nil {
			logger.Error(logLevel, logDetailCheckIn, fmt.Sprintf("failed to create agent: %v", err))
			return fmt.Errorf("failed to create agent: %w", err)
		}
		return nil
	})
}
