package main

import (
	"fmt"
	"goauth/config"
	"goauth/factory"
	"goauth/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	f, err := factory.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create factory: %v", err)
	}

	r := gin.Default()

	api := r.Group("/api")
	auth := api.Group("/auth")
	auth.POST("/register", f.AuthHandler.Register)
	auth.POST("/login", f.AuthHandler.Login)

	protected := api.Group("/protected")
	protected.Use(middleware.Auth(f.TokenService))
	protected.GET("/profile", func(ctx *gin.Context) {
		userID := ctx.GetString("userID")
		userRole := ctx.GetString("userRole")
		ctx.JSON(200, gin.H{"user_id": userID, "role": userRole})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
