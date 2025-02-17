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
	DeleteModule(ctx context.Context, moduleId string) error
	GetAllModules(ctx context.Context) ([]models.Module, error)
	GetModuleById(ctx context.Context, moduleId string) (*models.Module, error)
	DeleteAllModules(ctx context.Context) error
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
        INSERT INTO %s.modules (id, name, author, description, file_name, commands)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, dal.Schema)
	_, err = dal.DB.ExecContext(ctx, query, module.Id, module.Name, module.Author, module.Description, module.FileName, commandsJSON)
	if err != nil {
		return fmt.Errorf("failed to insert module: %w", err)
	}

	return nil
}

func (dal *ModuleDAL) DeleteModule(ctx context.Context, moduleId string) error {
	query := fmt.Sprintf(`DELETE FROM %s.modules WHERE id = $1`, dal.Schema)
	_, err := dal.DB.ExecContext(ctx, query, moduleId)
	if err != nil {
		return fmt.Errorf("failed to delete module: %w", err)
	}
	return nil
}

func (dal *ModuleDAL) GetAllModules(ctx context.Context) ([]models.Module, error) {
	query := fmt.Sprintf(`SELECT id, name, author, description, file_name, commands FROM %s.modules`, dal.Schema)
	rows, err := dal.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all modules: %w", err)
	}
	defer rows.Close()

	var modules []models.Module
	for rows.Next() {
		var module models.Module
		var commandsJSON []byte
		if err := rows.Scan(&module.Id, &module.Name, &module.Author, &module.Description, &module.FileName, &commandsJSON); err != nil {
			return nil, fmt.Errorf("failed to scan module: %w", err)
		}
		if err := json.Unmarshal(commandsJSON, &module.Commands); err != nil {
			return nil, fmt.Errorf("failed to unmarshal commands: %w", err)
		}
		modules = append(modules, module)
	}
	return modules, nil
}

func (dal *ModuleDAL) GetModuleById(ctx context.Context, moduleId string) (*models.Module, error) {
	query := fmt.Sprintf(`SELECT id, name, author, description, file_name, commands FROM %s.modules WHERE id = $1`, dal.Schema)
	row := dal.DB.QueryRowContext(ctx, query, moduleId)

	var module models.Module
	var commandsJSON []byte
	if err := row.Scan(&module.Id, &module.Name, &module.Author, &module.Description, &module.FileName, &commandsJSON); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get module by id: %w", err)
	}
	if err := json.Unmarshal(commandsJSON, &module.Commands); err != nil {
		return nil, fmt.Errorf("failed to unmarshal commands: %w", err)
	}
	return &module, nil
}

func (dal *ModuleDAL) DeleteAllModules(ctx context.Context) error {
	query := fmt.Sprintf(`DELETE FROM %s.modules`, dal.Schema)
	_, err := dal.DB.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete all modules: %w", err)
	}
	return nil
}
