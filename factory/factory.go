package factory

import (
	"database/sql"
	"fmt"
	"goauth/config"
	"goauth/database"
	"goauth/handler"
	"goauth/repository"
	"goauth/usecase"
)

type Factory struct {
	DB          *sql.DB
	AuthHandler *handler.AuthHandler
}

func New(cfg *config.Config) (*Factory, error) {
	db, err := database.OpenSQLite(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open sqlite factory: %w", err)
	}

	userRepo := repository.NewSQLiteUserRepo(db)
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg.BcryptCost)
	authHandler := handler.NewAuthHandler(authUsecase)

	return &Factory{
		DB:          db,
		AuthHandler: authHandler,
	}, nil
}
