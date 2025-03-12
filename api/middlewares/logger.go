package middlewares

import (
	"context"
	"strings"

	"github.com/1755/bookstore-api/internal/lgr"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/metrics" || strings.HasPrefix(c.Request.URL.Path, "/swagger") {
			c.Next()
			return
		}

		ctx := c.Request.Context()
		logger := lgr.GetLogger(ctx).With(
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
		)

		ctx = context.WithValue(ctx, lgr.CtxKey, logger)
		c.Request = c.Request.WithContext(ctx)
		c.Next()

		logger.Info("request completed", zap.Int("status", c.Writer.Status()))
	}
}
