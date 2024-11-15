package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ksel172/Meduza/teamserver/models"
)

var db *sql.DB

type UserDAL struct {
	db     Database
	schema string
}

func NewUsersDAL(db Database, schema string) *UserDAL {
	return &UserDAL{db: db, schema: schema}
}

func (dal *UserDAL) GetUsers(ctx context.Context) ([]models.User, error) {
	rows, err := dal.db.QueryContext(ctx, fmt.Sprintf("SELECT * FROM %s", dal.schema))
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
