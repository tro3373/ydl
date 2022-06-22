package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewRecordUaAndTimeHandler(logger *zap.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		oldTime := time.Now()
		ua := c.GetHeader("User-Agent")
		uuid := c.GetHeader("x-uuid")
		c.Next()
		logger.Info("Incoming Request => ",
			zap.String("path", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("uuid", uuid),
			zap.String("ua", ua),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("elapsed", time.Now().Sub(oldTime)),
		)
	}
}
