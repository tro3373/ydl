package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tro3373/ydl/cmd/middleware"
	"github.com/tro3373/ydl/cmd/request"
	"github.com/tro3373/ydl/cmd/response"
	"github.com/tro3373/ydl/cmd/util"
	"github.com/tro3373/ydl/cmd/worker"
	"go.uber.org/zap"
)

var logger *zap.Logger

func StartApi(ctx worker.Ctx) {
	engine := gin.Default()
	engine.Use(middleware.RecordUaAndTime)
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize zap logger: %v", err)
	}
	logger = zapLogger
	defer logger.Sync()

	handleAsServerError := func(c *gin.Context, message string, err error) {
		logger.Error(message, zap.String("Error:", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"status": "StatusInternalServerError"})
	}
	handleAsBadRequest := func(c *gin.Context, message string, err error) {
		logger.Error(message, zap.String("Error:", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
	}

	engine.GET("/api", func(c *gin.Context) {
		uuid := c.GetHeader("x-uuid")
		var key = c.Query("key")
		logger.Info("[GET] ==> ", zap.String("uuid", uuid), zap.String("key", key))

		doingFiles, err := findJsons(ctx.WorkDirs.Doing, key)
		if err != nil {
			handleAsServerError(c, "Failed to find doing jsons.", err)
			return
		}
		doneFiles, err := findJsons(ctx.WorkDirs.Done, key)
		if err != nil {
			handleAsServerError(c, "Failed to find done jsons.", err)
			return
		}

		jsons, err := readJsons(doingFiles, uuid, true)
		if err != nil {
			handleAsServerError(c, "Failed to read doing jsons.", err)
			return
		}
		doneJsons, err := readJsons(doneFiles, uuid, false)
		if err != nil {
			handleAsServerError(c, "Failed to read done jsons.", err)
			return
		}
		for _, doneJson := range doneJsons {
			jsons = append(jsons, doneJson)
		}
		c.JSON(http.StatusOK, gin.H{
			"list": jsons,
		})
	})

	engine.POST("/api", func(c *gin.Context) {
		var req request.Req
		if err := c.Bind(&req); err != nil {
			handleAsBadRequest(c, "Failed to bind request.", err)
			return
		}
		req.Uuid = c.GetHeader("x-uuid")
		logger.Info("[POST] ==> ", zap.Object("req", req))
		err := saveRequest(ctx.WorkDirs.Queue, req)
		if err != nil {
			logger.Error("Failed to save request.", zap.String("Error:", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"status": "StatusInternalServerError"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("ok"),
		})
	})

	engine.Run(":3000")
}

func findJsons(rootDir, key string) ([]string, error) {
	if len(key) == 0 {
		key = "*"
	}
	return filepath.Glob(filepath.Join(rootDir, key, "task.json"))
}

func getJsonPath(dir, key string) string {
	return filepath.Join(dir, fmt.Sprintf("%s.json", key))
}

func readJsons(doneFiles []string, uuid string, doing bool) ([]response.Res, error) {
	var reses []response.Res
	for _, file := range doneFiles {
		logger.Debug("[INFO] ==> ", zap.String("uuid", uuid), zap.String("file", file))
		raw, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		var task worker.Task
		json.Unmarshal(raw, &task)
		if !util.Contains(task.Uuids, uuid) {
			continue
		}
		reses = append(reses, response.NewRes(task, doing))
	}
	sort.Slice(reses, func(i, j int) bool {
		return reses[i].CreatedAt > reses[j].CreatedAt
	})
	return reses, nil
}

func saveRequest(dstRootDir string, req request.Req) error {
	key := request.Key(req.Url)

	timestamp := time.Now().Format("20060102_150405")
	req.CreatedAt = timestamp

	if !util.Exists(dstRootDir) {
		os.MkdirAll(dstRootDir, os.ModePerm)
	}
	dstFile := getJsonPath(dstRootDir, key)

	logger.Info("==> Saving request..", zap.String("dstFile", dstFile))
	data, _ := json.MarshalIndent(req, "", " ")
	return ioutil.WriteFile(dstFile, data, 0644)
}
