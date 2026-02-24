package factory

import (
	"database/sql"
	"fmt"
	"goauth/config"
	"goauth/database"
)

type Factory struct {
	DB *sql.DB
}

func New(cfg *config.Config) (*Factory, error) {
	db, err := database.OpenSQLite(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open sqlite factory: %w", err)
	}

	return &Factory{
		DB: db,
	}, nil
}
