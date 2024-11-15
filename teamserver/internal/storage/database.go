package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type Service interface {
	// Ping the storage
	// It returns an error if the connection is not made
	Ping() error

	// Close terminates the storage connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

type service struct {
	db *sql.DB
}

var (
	database   = utils.GetEnvString("DB_DATABASE", "")
	password   = utils.GetEnvString("DB_PASSWORD", "")
	username   = utils.GetEnvString("DB_USERNAME", "")
	port       = utils.GetEnvString("DB_PORT", "")
	host       = utils.GetEnvString("DB_HOST", "")
	schema     = utils.GetEnvString("DB_SCHEMA", "")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func (s *service) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return s.db.PingContext(ctx)
}

// Close closes the storage connection.
// It logs a message indicating the disconnection from the specific storage.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from storage: %s", database)
	return s.db.Close()
}
