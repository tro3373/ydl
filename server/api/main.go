package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"log"

	"github.com/gin-gonic/gin"

	"github.com/tro3373/ydl/api/middleware"
	"github.com/tro3373/ydl/api/request"
	"go.uber.org/zap"
)

var logger *zap.Logger

func main() {
	engine := gin.Default()
	engine.Use(middleware.RecordUaAndTime)
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err.Error())
	}
	logger = zapLogger

	workd := filepath.Join(".", "work")
	libd := filepath.Join(workd, "lib")

	queued := filepath.Join(workd, "queue")
	doned := filepath.Join(workd, "done")

	os.MkdirAll(libd, os.ModePerm)
	os.MkdirAll(doned, os.ModePerm)

	engine.GET("/api", func(c *gin.Context) {
		var key = c.Query("key")
		logger.Info("[INFO] ==> ", zap.String("key", key))
		files, err := findJsons(doned, key)
		if err != nil {
			logger.Error("[ERROR] ==> Failed to find jsons.", zap.String("Error:", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"status": "StatusInternalServerError"})
		}
		jsons, err := readJsons(files)
		if err != nil {
			logger.Error("[ERROR] ==> Failed to read jsons.", zap.String("Error:", err.Error()))
			c.JSON(http.StatusInternalServerError, gin.H{"status": "StatusInternalServerError"})
		}

		c.JSON(http.StatusOK, gin.H{
			"list": jsons,
		})
	})

	engine.POST("/api", func(c *gin.Context) {
		var req request.Exec
		if err := c.Bind(&req); err != nil {
			logger.Error("[ERROR] ==> Failed to bind request.", zap.String("Error:", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
			return
		}
		err := saveRequest(queued, req)
		if err != nil {
			logger.Error("[ERROR] ==> Failed to save request.", zap.String("Error:", err.Error()))
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
	return filepath.Glob(filepath.Join(rootDir, key, "req.json"))
}

func getJsonPath(dir, key string) string {
	return filepath.Join(dir, fmt.Sprintf("%s.json", key))
}

func readJsons(files []string) ([]request.Exec, error) {
	var res []request.Exec
	for _, file := range files {
		logger.Info("[INFO] ==> ", zap.String("file", file))
		raw, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		var req request.Exec
		json.Unmarshal(raw, &req)
		res = append(res, req)
	}
	return res, nil
}

func saveRequest(dstRootDir string, req request.Exec) error {
	key := req.Key()

	timestamp := time.Now().Format("20060102_150405")
	req.CreatedAt = timestamp

	if !exists(dstRootDir) {
		os.MkdirAll(dstRootDir, os.ModePerm)
	}
	dstFile := getJsonPath(dstRootDir, key)

	logger.Info("==> Saving request..", zap.String("dstFile", dstFile))
	data, _ := json.MarshalIndent(req, "", " ")
	return ioutil.WriteFile(dstFile, data, 0644)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
