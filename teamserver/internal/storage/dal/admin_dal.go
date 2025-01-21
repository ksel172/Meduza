package dal

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
)

type IAdminDAL interface {
	CreateDefaultAdmins(ctx context.Context, admin *models.ResAdmin) error
}

type AdminDAL struct {
	db     *sql.DB
	schema string
}

func NewAdminsDAL(db *sql.DB, schema string) *AdminDAL {
	return &AdminDAL{db: db, schema: schema}
}

func (adDal *AdminDAL) CreateDefaultAdmins(ctx context.Context, admin *models.ResAdmin) error {
	query := fmt.Sprintf(`INSERT INTO %s.users(username,pw_hash,role) VALUES($1,$2,$3)`, adDal.schema)

	_, err := adDal.db.ExecContext(ctx, query, admin.Adminname, admin.PasswordHash, "admin")
	if err != nil {
		logger.Error(fmt.Sprintf("Unable to create default admin: %v", err))
	}

	return err
}
