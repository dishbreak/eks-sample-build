package models

import (
	"database/sql"

	"github.com/dishbreak/sample-store-backend/config"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func Open(d config.DatabaseConfig) (*sql.DB, error) {
	return sql.Open("pgx", d.ToDSN())
}
