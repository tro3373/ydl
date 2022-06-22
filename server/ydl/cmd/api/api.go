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
	// @see [GolangのGin/bindataでシングルバイナリを試してみた(+React) - Qiita](https://qiita.com/wadahiro/items/4173788d54f028936723)
	// @see [【GO】gin + statikのシングルバイナリファイルサーバ | Narumium Blog](https://blog.narumium.net/2019/06/07/%E3%80%90go%E3%80%91gin-statik%E3%81%AE%E3%82%B7%E3%83%B3%E3%82%B0%E3%83%AB%E3%83%90%E3%82%A4%E3%83%8A%E3%83%AA%E3%83%95%E3%82%A1%E3%82%A4%E3%83%AB%E3%82%B5%E3%83%BC%E3%83%90/)
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
