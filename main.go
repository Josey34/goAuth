package main

import (
	"fmt"
	"goauth/config"
	"goauth/factory"
	"goauth/middleware"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if cfg.LogLevel == "debug" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	f, err := factory.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create factory: %v", err)
	}

	rateLimiter := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)

	r := gin.Default()

	r.Use(gin.Recovery())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORS(cfg.AllowedOrigins))
	r.Use(middleware.Logger(zlog.Logger))
	r.Use(rateLimiter.Limit())

	api := r.Group("/api")
	auth := api.Group("/auth")
	auth.POST("/register", f.AuthHandler.Register)
	auth.POST("/login", f.AuthHandler.Login)
	auth.POST("/refresh", f.AuthHandler.Refresh)

	protected := api.Group("/auth")
	protected.Use(middleware.Auth(f.TokenService))
	protected.GET("/profile", f.UserHandler.GetProfile)
	protected.PUT("/profile", f.UserHandler.UpdateProfile)

	admin := api.Group("/admin")
	admin.Use(middleware.Auth(f.TokenService))
	admin.Use(middleware.RequireRole("admin"))
	admin.GET("/dashboard", f.UserHandler.AdminDashboard)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
