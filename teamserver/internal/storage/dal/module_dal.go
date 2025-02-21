package dal

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type IModuleDAL interface {
	CreateModule(ctx context.Context, module *models.Module) error
	DeleteModule(ctx context.Context, moduleId string) error
	GetAllModules(ctx context.Context) ([]models.Module, error)
	GetModuleById(ctx context.Context, moduleId string) (*models.Module, error)
	DeleteAllModules(ctx context.Context) error
}
type ModuleDAL struct {
	db     *sql.DB
	schema string
}

func NewModuleDAL(db *sql.DB, schema string) *ModuleDAL {
	return &ModuleDAL{
		db:     db,
		schema: schema,
	}
}

func (dal *ModuleDAL) CreateModule(ctx context.Context, module *models.Module) error {
	query := fmt.Sprintf(`
        INSERT INTO %s.modules (id, name, author, description, file_name, commands)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		commandsJSON, err := json.Marshal(module.Commands)
		if err != nil {
			logger.Error(logLevel, logDetailModule, fmt.Sprintf("failed to marshal commands: %v", err))
			return fmt.Errorf("failed to marshal commands: %w", err)
		}

		_, err = stmt.ExecContext(ctx, module.Id, module.Name, module.Author, module.Description, module.FileName, commandsJSON)
		if err != nil {
			logger.Error(logLevel, logDetailModule, fmt.Sprintf("failed to create module: %v", err))
			return fmt.Errorf("failed to create module: %w", err)
		}

		return nil
	})
}

func (dal *ModuleDAL) DeleteModule(ctx context.Context, moduleId string) error {
	query := fmt.Sprintf(`DELETE FROM %s.modules WHERE id = $1`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx, moduleId)
		if err != nil {
			logger.Error(logLevel, logDetailModule, fmt.Sprintf("failed to delete module: %v", err))
			return fmt.Errorf("failed to delete module: %w", err)
		}
		return nil
	})
}

func (dal *ModuleDAL) GetAllModules(ctx context.Context) ([]models.Module, error) {
	query := fmt.Sprintf(`SELECT id, name, author, description, file_name, commands FROM %s.modules`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) ([]models.Module, error) {
		rows, err := stmt.QueryContext(ctx)
		if err != nil {
			logger.Error(logLevel, logDetailModule, fmt.Sprintf("failed to get modules: %v", err))
			return nil, fmt.Errorf("failed to get modules: %w", err)
		}
		defer rows.Close()

		var modules []models.Module
		for rows.Next() {
			var module models.Module
			var commandsJSON []byte
			if err := rows.Scan(&module.Id, &module.Name, &module.Author, &module.Description, &module.FileName, &commandsJSON); err != nil {
				logger.Error(logLevel, logDetailModule, fmt.Sprintf("failed to scan module: %v", err))
				return nil, fmt.Errorf("failed to scan module: %w", err)
			}
			if err := json.Unmarshal(commandsJSON, &module.Commands); err != nil {
				logger.Error(logLevel, logDetailModule, fmt.Sprintf("failed to unmarshal commands: %v", err))
				return nil, fmt.Errorf("failed to unmarshal commands: %w", err)
			}
			modules = append(modules, module)
		}
		return modules, nil
	})

}

func (dal *ModuleDAL) GetModuleById(ctx context.Context, moduleId string) (*models.Module, error) {
	query := fmt.Sprintf(`SELECT id, name, author, description, file_name, commands FROM %s.modules WHERE id = $1`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) (*models.Module, error) {
		row := stmt.QueryRowContext(ctx, moduleId)

		var module models.Module
		var commandsJSON []byte
		if err := row.Scan(&module.Id, &module.Name, &module.Author, &module.Description, &module.FileName, &commandsJSON); err != nil {
			if err == sql.ErrNoRows {
				logger.Error(logLevel, logDetailModule, fmt.Sprintf("no module found with ID '%s': %v", moduleId, err))
				return nil, fmt.Errorf("no module found: %w", err)
			}
			return nil, fmt.Errorf("failed to get module by id: %w", err)
		}
		if err := json.Unmarshal(commandsJSON, &module.Commands); err != nil {
			logger.Error(logLevel, logDetailModule, fmt.Sprintf("failed to unmarshal commands: %v", err))
			return nil, fmt.Errorf("failed to unmarshal commands: %w", err)
		}
		return &module, nil
	})
}

func (dal *ModuleDAL) DeleteAllModules(ctx context.Context) error {
	query := fmt.Sprintf(`DELETE FROM %s.modules`, dal.schema)
	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx)
		if err != nil {
			logger.Error(logLevel, logDetailModule, fmt.Sprintf("failed to delete all modules: %v", err))
			return fmt.Errorf("failed to delete all modules: %w", err)
		}
		return nil
	})
}
