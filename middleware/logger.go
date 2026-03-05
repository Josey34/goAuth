package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func Logger(log zerolog.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := uuid.New().String()
		ctx.Set("requestID", requestID)

		start := time.Now()
		ctx.Next()
		duration := time.Since(start)

		log.Info().
			Str("request_id", requestID).
			Str("method", ctx.Request.Method).
			Str("path", ctx.Request.RequestURI).
			Int("status", ctx.Writer.Status()).
			Dur("duration_ms", duration).
			Msg("request completed")
	}
}
