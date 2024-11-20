package dal

import (
	"context"
	"fmt"
	"github.com/ksel172/Meduza/teamserver/internal/models"
	"github.com/ksel172/Meduza/teamserver/internal/storage"
)

type UserDAL struct {
	db     storage.Database
	schema string
}

func NewUsersDAL(db storage.Database, schema string) *UserDAL {
	return &UserDAL{db: db, schema: schema}
}

func (dal *UserDAL) GetUsers(ctx context.Context) ([]models.User, error) {
	rows, err := dal.db.QueryContext(ctx, fmt.Sprintf("SELECT id, username, pw_hash, role_id, created_ts, updated_ts FROM %s", dal.schema))
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

// TODO add create user
func (dal *UserDAL) CreateUser(ctx context.Context, user models.User) error {
	return nil
}
