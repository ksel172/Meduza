package repos

import (
	/* 	"context" */
	"database/sql"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/pkg/conf"
	_ "github.com/lib/pq"
)

func Setup() (*sql.DB, error) {

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable search_path=%s", conf.GetMeduzaDbHostname(), conf.GetMeduzaDbPort(), conf.GetMeduzaDbUsername(), conf.GetMeduzaDbPassword(), conf.GetMeduzaDbName(), conf.GetMeduzaDbSchema())

	return sql.Open("postgres", connectionString)
}
