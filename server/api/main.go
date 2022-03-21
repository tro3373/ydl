package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"log"

	"github.com/gin-gonic/gin"

	"github.com/tro3373/ydl/middleware"
	"github.com/tro3373/ydl/request"
	"go.uber.org/zap"
)

func main() {
	engine := gin.Default()
	engine.Use(middleware.RecordUaAndTime)
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err.Error())
	}

	workd := filepath.Join(".", "work")
	libd := filepath.Join(workd, "lib")

	queued := filepath.Join(workd, "queue")
	doned := filepath.Join(workd, "done")

	os.MkdirAll(libd, os.ModePerm)
	os.MkdirAll(doned, os.ModePerm)

	engine.GET("/api", func(c *gin.Context) {
		var key = c.Query("key")
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

		// if len(key) > 0 {
		// 	// logger.Warn("[WARN] ==> Params key is empty. Set default test key.")
		// 	// key = "zkZARKFuzNQ"
		// 	dir := filepath.Join(doned, key)
		// 	if exists(dir) {
		// 		logger.Info("Getting ", zap.String("Key", key))
		// 		// TODO
		// 	}
		// }
		// err := downloadAudio(libd, workd, key)
		// if err != nil {
		// 	log.Fatalf("Failed to download. key:%s err:%v", key, err)
		// 	return
		// }

		// c.JSON(http.StatusOK, gin.H{
		// 	"message": fmt.Sprintf("Success to download key:%s", key),
		// })
	})

	engine.POST("/api", func(c *gin.Context) {
		var req request.Exec
		if err := c.Bind(&req); err != nil {
			logger.Error("[ERROR] ==> Failed to bind request.", zap.String("Error:", err.Error()))
			c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
			return
		}
		err := saveRequest(logger, queued, req)
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

func readJsons(files []string) ([]request.Exec, error) {
	var res []request.Exec
	for _, file := range files {
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

func saveRequest(logger *zap.Logger, dstRootDir string, req request.Exec) error {
	key := req.Key()

	timestamp := time.Now().Format("20060102_150405")
	req.CreatedAt = timestamp

	dstDir := filepath.Join(dstRootDir, key)
	dstFile := filepath.Join(dstDir, "req.json")

	if !exists(dstDir) {
		os.MkdirAll(dstDir, os.ModePerm)
	}

	logger.Info("==> Saving request..", zap.String("dstFile", dstFile))
	data, _ := json.MarshalIndent(req, "", " ")
	return ioutil.WriteFile(dstFile, data, 0644)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// func findResults() {
// 	// out, err := exec.Command("ls").Output()
// 	var str = filepath.Join(".", workd, key+"*")
// 	logger.Info("Key is", zap.String("key", key), zap.String("str", str))
// 	out, err := exec.Command("ls", "-la", str).Output()
// 	if err != nil {
// 		// log.Fatalf("Failed list workd: %v", err)
// 		logger.Error("Failed to list %v", zap.Error(err))
// 		c.JSON(http.StatusOK, gin.H{
// 			"message": fmt.Sprintf("Failed to list %s %v", str, err),
// 		})
// 		return
// 	}
// 	logger.Info("### ls -la result", zap.String("out", string(out)))
// 	fmt.Printf("@@@ %s", string(out))
// 	logger.Info("test")
// 	fmt.Printf("testt")
// }

func downloadAudio(libd, workd, key string) error {
	return startDownload(libd, workd, key, "", true)
}

func downloadMovie(libd, workd, key string) error {
	return startDownload(libd, workd, key, "", false)
}

func startDownload(libd, workd, key, format string, audio bool) error {
	res, err := existsPrefix(filepath.Join(workd, key))
	if err != nil {
		return err
	}
	if res {
		return nil
	}
	return executeYoutubeDl(libd, workd, key, format, audio)
}

func existsPrefix(name string) (bool, error) {
	matches, err := filepath.Glob(name + ".*")
	if err != nil {
		return false, err
	}
	return len(matches) > 0, nil
}

func executeYoutubeDl(libd, workd, key, format string, audio bool) error {
	var args []string
	args = append(args, key)
	args = append(args, "-o")
	if len(format) == 0 {
		format = "%(id)s_%(title)s.%(ext)s"
	}
	args = append(args, filepath.Join(workd, format))
	if audio {
		args = append(args, "-x")
		args = append(args, "--audio-format")
		args = append(args, "mp3")
	}
	err := exec.Command(filepath.Join(libd, "youtube-dl"), args...).Run()
	if err != nil {
		log.Fatalf("Failed to executeYoutubeDl %s: %v", key, err)
		return err
	}
	return nil
}
