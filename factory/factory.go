package factory

import (
	"database/sql"
	"fmt"
	"goauth/config"
	"goauth/database"
	"goauth/handler"
	"goauth/repository"
	"goauth/service"
	"goauth/usecase"
)

type Factory struct {
	DB           *sql.DB
	AuthHandler  *handler.AuthHandler
	TokenService service.TokenService
	UserHandler  *handler.UserHandler
}

func New(cfg *config.Config) (*Factory, error) {
	db, err := database.OpenSQLite(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open sqlite factory: %w", err)
	}

	userRepo := repository.NewSQLiteUserRepo(db)
	tokenService := service.NewJWTTokenService(cfg.JWTSecret, cfg.AccessTTL, cfg.RefreshTTL)
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg.BcryptCost, tokenService)
	authHandler := handler.NewAuthHandler(authUsecase)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	return &Factory{
		DB:           db,
		AuthHandler:  authHandler,
		TokenService: tokenService,
		UserHandler:  userHandler,
	}, nil
}
