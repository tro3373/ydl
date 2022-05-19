package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tro3373/ydl/cmd/middleware"
	"github.com/tro3373/ydl/cmd/request"
	"github.com/tro3373/ydl/cmd/util"
	"github.com/tro3373/ydl/cmd/worker"
	"go.uber.org/zap"
)

// // apiCmd represents the api command
// var apiCmd = &cobra.Command{
// 	Use:   "api",
// 	Short: "A brief description of your command",
// 	Long: `A longer description that spans multiple lines and likely contains examples
// and usage of using your command. For example:
//
// Cobra is a CLI library for Go that empowers applications.
// This application is a tool to generate the needed files
// to quickly create a Cobra application.`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		StartApi(args)
// 	},
// }

var logger *zap.Logger

// func init() {
// 	rootCmd.AddCommand(apiCmd)
//
// 	// Here you will define your flags and configuration settings.
//
// 	// Cobra supports Persistent Flags which will work for this command
// 	// and all subcommands, e.g.:
// 	// apiCmd.PersistentFlags().String("foo", "", "A help for foo")
//
// 	// Cobra supports local flags which will only run when this command
// 	// is called directly, e.g.:
// 	// apiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
// }
//
func StartApi(ctx worker.Ctx) {
	engine := gin.Default()
	engine.Use(middleware.RecordUaAndTime)
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err.Error())
	}
	logger = zapLogger

	engine.GET("/api", func(c *gin.Context) {
		var key = c.Query("key")
		logger.Info("[INFO] ==> ", zap.String("key", key))
		files, err := findJsons(ctx.DoneDir, key)
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
		err := saveRequest(ctx.QueueDir, req)
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

	if !util.Exists(dstRootDir) {
		os.MkdirAll(dstRootDir, os.ModePerm)
	}
	dstFile := getJsonPath(dstRootDir, key)

	logger.Info("==> Saving request..", zap.String("dstFile", dstFile))
	data, _ := json.MarshalIndent(req, "", " ")
	return ioutil.WriteFile(dstFile, data, 0644)
}
