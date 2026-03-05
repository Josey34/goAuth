package middleware

import "github.com/gin-gonic/gin"

func SecurityHeaders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.Header("X-Frame-Options", "DENY")
		ctx.Header("X-XSS-Protection", "1; mode=block")
		ctx.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		ctx.Next()
	}
}
