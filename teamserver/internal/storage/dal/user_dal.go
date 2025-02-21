package dal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type IUserDAL interface {
	AddUsers(context.Context, *models.ResUser) error
	GetUsers(context.Context) ([]models.User, error)
	GetUserByUsername(context.Context, string) (*models.ResUser, error)
	GetUserById(context.Context, string) (*models.ResUser, error)
}

type UserDAL struct {
	db     *sql.DB
	schema string
}

func NewUsersDAL(db *sql.DB, schema string) *UserDAL {
	return &UserDAL{db: db, schema: schema}
}

func (dal *UserDAL) AddUsers(ctx context.Context, user *models.ResUser) error {
	query := fmt.Sprintf(`INSERT INTO %s.users(username,pw_hash,role) VALUES($1,$2,$3)`, dal.schema)

	return utils.WithTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) error {
		_, err := stmt.ExecContext(ctx, user.Username, user.PasswordHash, user.Role)
		if err != nil {
			logger.Error(logLevel, logDetailUser, fmt.Sprintf("Unable to add user: %v", err))
			return fmt.Errorf("failed to add user: %w", err)
		}
		return err
	})
}

func (dal *UserDAL) GetUsers(ctx context.Context) ([]models.User, error) {
	query := fmt.Sprintf(`SELECT id, username, pw_hash, role, created_at, updated_at
		FROM %s.users WHERE deleted_at IS NULL ORDER BY created_at DESC`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) ([]models.User, error) {
		rows, err := stmt.QueryContext(ctx)
		if err != nil {
			logger.Error(logLevel, logDetailUser, fmt.Sprintf("Failed to get users: %v", err))
			return nil, fmt.Errorf("failed to get users: %w", err)
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
				logger.Error(logLevel, logDetailUser, fmt.Sprintf("Failed to scan user row: %v", err))
				return nil, fmt.Errorf("failed to get users: %w", err)
			}
			users = append(users, user)
		}

		return users, nil
	})

}

func (dal *UserDAL) GetUserByUsername(ctx context.Context, username string) (*models.ResUser, error) {
	query := fmt.Sprintf(`SELECT id, username, pw_hash, role FROM %s.users WHERE username = $1 AND deleted_at IS NULL`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) (*models.ResUser, error) {
		var user models.ResUser
		err := stmt.QueryRowContext(ctx, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("user not found")
			}
			logger.Error(logLevel, logDetailUser, fmt.Sprintf("Failed to fetch user: %v", err))
			return nil, fmt.Errorf("failed to fetch user: %w", err)
		}
		return &user, nil
	})
}

func (dal *UserDAL) GetUserById(ctx context.Context, id string) (*models.ResUser, error) {
	query := fmt.Sprintf(`SELECT id , username , pw_hash, role FROM %s.users WHERE id = $1 AND deleted_at IS NULL`, dal.schema)

	return utils.WithResultTimeout(ctx, dal.db, query, 5, func(ctx context.Context, stmt *sql.Stmt) (*models.ResUser, error) {
		var user models.ResUser
		err := stmt.QueryRowContext(ctx, id).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("user not found")
			}
			logger.Error(logLevel, logDetailUser, fmt.Sprintf("Failed to get user: %v", err))
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
		return &user, nil
	})
}
