package middleware

import (
	"goauth/dto"
	"goauth/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(tokenService service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := tokenService.Validate(token)
		if err != nil {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("userID", claims["sub"])
		c.Set("userRole", claims["role"])
		c.Next()
	}
}
