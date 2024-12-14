package dal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/lib/pq"
)

type ListenersDAL struct {
	db     *sql.DB
	schema string
}

func NewListenersDAL(db *sql.DB, schema string) *ListenersDAL {
	return &ListenersDAL{db: db, schema: schema}
}

func (dal *ListenersDAL) CreateListeners(ctx context.Context, listener *models.Listener) error {
	query := fmt.Sprintf(`
		INSERT INTO %s.listeners (
			type, host, port, status, description, 
			created_at, udpated_at, started_at, stopped_at, 
			certificate_path, key_path, whitelist_enabled, whitelist, 
			blacklist_enabled, blacklist, logging_enabled, log_path, log_level
		) VALUES (
			$1, $2, $3, $4, $5, 
			$6, $7, $8, $9, 
			$10, $11, $12, $13, 
			$14, $15, $16, $17, $18
		) RETURNING id`, dal.schema)

	var id string

	err := dal.db.QueryRowContext(ctx, query,
		listener.Type, listener.Host,
		listener.Port, listener.Status,
		listener.Description, listener.CreatedAt,
		listener.UpdatedAt, listener.StartedAt,
		listener.StoppedAt, listener.CertPath, listener.KeyPath,
		listener.WhitelistEnabled, pq.Array(listener.Whitelist),
		listener.BlacklistEnabled, pq.Array(listener.Blacklist),
		listener.LoggingEnabled, listener.LogPath,
		listener.LogLevel).Scan(&id)

	if err != nil {
		logger.Error("Failed to insert Listener:", err)
		return fmt.Errorf("Failed to insert listener")
	}

	logger.Good("Inserted Listener with ID - ", id)

	return nil
}
