package api

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"
	"github.com/tro3373/ydl/cmd/api/handler"
	"github.com/tro3373/ydl/cmd/api/middleware"
	"github.com/tro3373/ydl/cmd/worker/ctx"
	_ "github.com/tro3373/ydl/statik"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func Start(ctx ctx.Ctx) {
	err := startInner(ctx)
	if err != nil {
		log.Fatalf("Failed to start api: %v", err)
	}
}

func startInner(ctx ctx.Ctx) error {
	r := gin.Default()

	zapLogger, err := newLogger()
	if err != nil {
		return fmt.Errorf("Failed to initialize zap logger : %w", err)
	}
	logger = zapLogger
	defer logger.Sync()

	r.Use(middleware.NewRecordUaAndTimeHandler(logger))
	r.Use(middleware.NewResourceHandler("/resource", "./work", logger))
	statikFS, err := fs.New()
	if err != nil {
		return fmt.Errorf("Failed to initialize statik fs: %w", err)
	}
	r.Use(static.Serve("/", middleware.StatikFileSystem(statikFS)))

	handler := handler.NewHandler(logger, ctx)

	v1 := r.Group("/api")
	v1.GET("", handler.GetDones)
	v1.POST("", handler.CreateQueue)
	v1.DELETE("/:key", handler.DeleteDone)

	if err := r.Run(":3000"); err != nil {
		return fmt.Errorf("Failed to run server: %w", err)
	}
	return nil
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
