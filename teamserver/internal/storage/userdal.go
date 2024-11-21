package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ksel172/Meduza/teamserver/models"
)

var db *sql.DB

type UserDAL struct {
	db     *sql.DB
	schema string
}

func NewUsersDAL(db *sql.DB, schema string) *UserDAL {
	return &UserDAL{db: db, schema: schema}
}

func (dal *UserDAL) AddUsers(ctx context.Context, user *models.ResUser) error {
	query := fmt.Sprintf(`INSERT INTO %s.users(username,pw_hash,role) VALUES($1,$2,$3)`, dal.schema)
	_, err := dal.db.ExecContext(ctx, query, user.Username, user.PasswordHash, user.Role)
	return err
}

func (dal *UserDAL) GetUsers(ctx context.Context) ([]models.User, error) {
	rows, err := dal.db.QueryContext(ctx, fmt.Sprintf("SELECT id, username, pw_hash, role, created_at, updated_at FROM %s.users WHERE deleted_at IS NULL ORDER BY created_at DESC", dal.schema))
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("Failed to fetch users: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (dal *UserDAL) GetUserByUsername(ctx context.Context, username string) (*models.ResUser, error) {
	query := fmt.Sprintf(`SELECT id , username , pw_hash, role FROM %s.users WHERE username = $1 AND deleted_at IS NULL`, dal.schema)
	var user models.ResUser
	err := dal.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("Failed to fetch user: %w", err)
	}
	return &user, nil
}

func (dal *UserDAL) GetUserById(ctx context.Context, id string) (*models.ResUser, error) {
	query := fmt.Sprintf(`SELECT id , username , pw_hash, role FROM %s.users WHERE id = $1 AND deleted_at IS NULL`, dal.schema)
	var user models.ResUser
	err := dal.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("Failed to fetch user: %w", err)
	}
	return &user, nil
}
