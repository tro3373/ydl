package api

import (
	"log"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/tro3373/ydl/cmd/api/handler"
	"github.com/tro3373/ydl/cmd/api/middleware"
	"github.com/tro3373/ydl/cmd/worker/ctx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func Start(ctx ctx.Ctx) {
	engine := gin.Default()

	engine.Use(static.Serve("/", static.LocalFile("./assets", true)))
	engine.Use(middleware.RecordUaAndTime)
	engine.NoRoute(func(c *gin.Context) {
		c.File("./assets/index.html")
	})

	// engine.Static("/", "./assets")
	// zapLogger, err := zap.NewProduction()
	zapLogger, err := newLogger()
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}
	logger = zapLogger
	defer logger.Sync()

	handler := handler.NewHandler(logger, ctx)

	v1 := engine.Group("/api")
	v1.GET("", handler.GetDones)
	v1.POST("", handler.CreateQueue)
	v1.DELETE("/:key", handler.DeleteDone)

	engine.Run(":3000")
}

func newLogger() (*zap.Logger, error) {
	// logger, err := zap.NewProduction()
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "console"
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	cfg.EncoderConfig.ConsoleSeparator = " "
	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return cfg.Build()
}
