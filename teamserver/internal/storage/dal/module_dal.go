package dal

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/models"
)

type IModuleDAL interface {
	CreateModule(ctx context.Context, module *models.Module) error
}

type ModuleDAL struct {
	DB     *sql.DB
	Schema string
}

func NewModuleDAL(db *sql.DB, schema string) *ModuleDAL {
	return &ModuleDAL{
		DB:     db,
		Schema: schema,
	}
}

func (dal *ModuleDAL) CreateModule(ctx context.Context, module *models.Module) error {
	commandsJSON, err := json.Marshal(module.Commands)
	if err != nil {
		return fmt.Errorf("failed to marshal commands: %w", err)
	}

	query := fmt.Sprintf(`
        INSERT INTO %s.modules (id, name, author, description, file_name, usage, commands)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, dal.Schema)
	_, err = dal.DB.ExecContext(ctx, query, module.Id, module.Name, module.Author, module.Description, module.FileName, module.Usage, commandsJSON)
	if err != nil {
		return fmt.Errorf("failed to insert module: %w", err)
	}

	return nil
}
