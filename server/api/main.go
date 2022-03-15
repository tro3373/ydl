package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"log"

	"github.com/gin-gonic/gin"

	"github.com/tro3373/ydl/middleware"
	"go.uber.org/zap"
)

type Req struct {
	Url    string `json:"url" binding:"required"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
	Jenre  string `json:"jenre"`
}

func (r Req) Key() string {
	if len(r.Url) == 0 {
		return ""
	}
	exp := regexp.MustCompile(`^http.*watch\?v\=`)
	if exp.MatchString(r.Url) {
		return exp.ReplaceAllString(r.Url, "")
	}
	return r.Url
}

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

	os.MkdirAll(libd, os.ModePerm)

	engine.GET("/api", func(c *gin.Context) {
		var key = c.Query("key")
		if len(key) == 0 {
			logger.Warn("[WARN] ==> Params key is empty. Set default test key.")
			key = "zkZARKFuzNQ"
		}
		err := downloadAudio(libd, workd, key)
		if err != nil {
			log.Fatalf("Failed to download. key:%s err:%v", key, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Success to download key:%s", key),
		})
	})

	engine.POST("/api", func(c *gin.Context) {
		var req Req
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

func saveRequest(logger *zap.Logger, queued string, req Req) error {
	key := req.Key()

	timestamp := time.Now().Format("20060102_150405")

	dstDir := filepath.Join(queued, fmt.Sprintf("%s.%s", timestamp, key))
	os.MkdirAll(dstDir, os.ModePerm)
	jsonPath := filepath.Join(dstDir, fmt.Sprintf("%s.json", key))
	// logger.Info("==> Saving request..", zap.String("jsonPath", jsonPath), zap.String("req", fmt.Sprintf("%#+v", req)))
	logger.Info("==> Saving request..", zap.String("jsonPath", jsonPath))

	data, _ := json.MarshalIndent(req, "", " ")
	err := ioutil.WriteFile(jsonPath, data, 0644)
	if err == nil {
	}
	return err
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
