package dal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
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
	tx, err := dal.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert agent
	logger.Debug(layer, "Creating agent: "+agent.ID)
	agentQuery := fmt.Sprintf(`
        INSERT INTO %s.agents (id, name, note, status, first_callback, last_callback, modified_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`, dal.schema)

	_, err = tx.Exec(agentQuery, agent.AgentID, agent.Name, agent.Note, agent.Status,
		agent.FirstCallback, agent.LastCallback, agent.ModifiedAt)
	if err != nil {
		logger.Error(layer, fmt.Sprintf("failed to insert agent in database: %v", err))
		return fmt.Errorf("failed to insert agent: %w", err)
	}

	return tx.Commit()
}
