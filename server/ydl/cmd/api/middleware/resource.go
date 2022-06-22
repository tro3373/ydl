package middleware

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func NewResourceHandler(prefix, localFilePath string, logger *zap.Logger) func(c *gin.Context) {

	fs := static.LocalFile(localFilePath, true)
	fn := static.Serve(prefix, fs)
	return func(c *gin.Context) {
		if !fs.Exists(prefix, c.Request.URL.Path) {
			c.Next()
			return
		}
		fnm := c.Query("f")
		if len(fnm) > 0 {
			fnm = " filename=" + fnm
		}
		c.Header("Content-Disposition", "attachment;"+fnm)
		c.Header("Content-Type", "application/octet-stream")
		fn(c)
		c.Abort()
	}
}
