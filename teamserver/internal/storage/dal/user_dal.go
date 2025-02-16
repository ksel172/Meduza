package dal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
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
	logger.Debug(logLevel, "Adding user: "+user.ID)
	query := fmt.Sprintf(`INSERT INTO %s.users(username,pw_hash,role) VALUES($1,$2,$3)`, dal.schema)
	_, err := dal.db.ExecContext(ctx, query, user.Username, user.PasswordHash, user.Role)
	if err != nil {
		logger.Error(logLevel, "Unable to add user: ", err)
		return fmt.Errorf("failed to add user: %v", err)
	}
	return err
}

func (dal *UserDAL) GetUsers(ctx context.Context) ([]models.User, error) {
	logger.Debug(logLevel, "Fetching all users")
	rows, err := dal.db.QueryContext(ctx, fmt.Sprintf("SELECT id, username, pw_hash, role, created_at, updated_at FROM %s.users WHERE deleted_at IS NULL ORDER BY created_at DESC", dal.schema))
	if err != nil {
		logger.Error(logLevel, "Failed to fetch users: ", err)
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			logger.Error(logLevel, "Failed to scan user row: ", err)
			return nil, fmt.Errorf("failed to fetch users: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (dal *UserDAL) GetUserByUsername(ctx context.Context, username string) (*models.ResUser, error) {
	logger.Debug(logLevel, "Fetching user by username: "+username)
	query := fmt.Sprintf(`SELECT id , username , pw_hash, role FROM %s.users WHERE username = $1 AND deleted_at IS NULL`, dal.schema)
	var user models.ResUser
	err := dal.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		logger.Error(logLevel, "Failed to fetch user: ", err)
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	return &user, nil
}

func (dal *UserDAL) GetUserById(ctx context.Context, id string) (*models.ResUser, error) {
	logger.Debug(logLevel, "Fetching user by id: "+id)
	query := fmt.Sprintf(`SELECT id , username , pw_hash, role FROM %s.users WHERE id = $1 AND deleted_at IS NULL`, dal.schema)
	var user models.ResUser
	err := dal.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		logger.Error(logLevel, "Failed to fetch user: ", err)
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	return &user, nil
}
