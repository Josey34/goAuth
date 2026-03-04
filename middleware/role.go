package middleware

import (
	"goauth/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole := ctx.GetString("userRole")

		if userRole != requiredRole {
			ctx.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "Insufficient permissions"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
