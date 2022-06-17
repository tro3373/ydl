package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tro3373/ydl/cmd/api/request"
	"github.com/tro3373/ydl/cmd/api/response"
	"github.com/tro3373/ydl/cmd/util"
	"github.com/tro3373/ydl/cmd/worker/ctx"
	"github.com/tro3373/ydl/cmd/worker/task"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
	ctx    ctx.Ctx
}

func NewHandler(logger *zap.Logger, ctx ctx.Ctx) *Handler {
	return &Handler{
		logger,
		ctx,
	}
}

func (h *Handler) GetDones(c *gin.Context) {
	uuid := c.GetHeader("x-uuid")
	key := c.Query("key")
	h.logger.Info("[GET] ==> ", zap.String("uuid", uuid), zap.String("key", key))

	doingFiles, err := h.findJsons(h.ctx.WorkDirs.Doing, key)
	if err != nil {
		h.handleAsServerError(c, "Failed to find doing jsons.", err)
		return
	}
	doneFiles, err := h.findJsons(h.ctx.WorkDirs.Done, key)
	if err != nil {
		h.handleAsServerError(c, "Failed to find done jsons.", err)
		return
	}

	jsons, err := h.readJsons(doingFiles, uuid, true)
	if err != nil {
		h.handleAsServerError(c, "Failed to read doing jsons.", err)
		return
	}
	doneJsons, err := h.readJsons(doneFiles, uuid, false)
	if err != nil {
		h.handleAsServerError(c, "Failed to read done jsons.", err)
		return
	}
	for _, doneJson := range doneJsons {
		jsons = append(jsons, doneJson)
	}
	h.list(c, jsons)
}

func (h *Handler) CreateQueue(c *gin.Context) {
	var req request.Req
	if err := c.Bind(&req); err != nil {
		h.handleAsBadRequest(c, "Failed to bind request.", err)
		return
	}
	req.Uuid = c.GetHeader("x-uuid")
	h.logger.Info("[POST] ==> ", zap.Object("req", req))
	err := h.saveRequest(h.ctx.WorkDirs.Queue, req)
	if err != nil {
		h.logger.Error("Failed to save request.", zap.String("Error:", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"status": "StatusInternalServerError"})
		return
	}
	h.ok(c)
}

func (h *Handler) DeleteDone(c *gin.Context) {
	uuid := c.GetHeader("x-uuid")
	key := c.Param("key")
	if len(key) == 0 {
		h.handleAsBadRequest(c, "Empty key request.", nil)
		return
	}
	h.logger.Info("[DELETE] ==> ", zap.String("uuid", uuid), zap.String("key", key))
	if err := h.removeRequest(h.ctx.WorkDirs.Done, key); err != nil {
		h.handleAsServerError(c, fmt.Sprintf("Failed to remove %s.", key), err)
		return
	}
	h.logger.Info("[DELETE] Deleted", zap.String("key", key))
	h.ok(c)
}

func (h *Handler) findJsons(rootDir, key string) ([]string, error) {
	if len(key) == 0 {
		key = "*"
	}
	return filepath.Glob(filepath.Join(rootDir, key, "task.json"))
}

func (h *Handler) getJsonPath(dir, key string) string {
	return filepath.Join(dir, fmt.Sprintf("%s.json", key))
}

func (h *Handler) readJsons(doneFiles []string, uuid string, doing bool) ([]response.Res, error) {
	var reses []response.Res
	for _, file := range doneFiles {
		h.logger.Debug("[INFO] ==> ", zap.String("uuid", uuid), zap.String("file", file))
		raw, err := ioutil.ReadFile(filepath.Clean(file))
		if err != nil {
			return nil, err
		}
		var task task.Task

		if err := json.Unmarshal(raw, &task); err != nil {
			return nil, err
		}
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

func (h *Handler) saveRequest(dstRootDir string, req request.Req) error {
	key := request.Key(req.Url)

	timestamp := time.Now().Format("20060102_150405")
	req.CreatedAt = timestamp

	if !util.Exists(dstRootDir) {
		//#nosec G301
		if err := os.MkdirAll(dstRootDir, 0775); err != nil {
			return fmt.Errorf("Failed to create directory %s %d: %w", dstRootDir, 0775, err)
		}
	}
	dstFile := h.getJsonPath(dstRootDir, key)

	h.logger.Info("==> Saving request..", zap.String("dstFile", dstFile))
	data, _ := json.MarshalIndent(req, "", " ")
	//#nosec G306
	return ioutil.WriteFile(dstFile, data, 0664)
}

func (h *Handler) removeRequest(rootDir, key string) error {
	if len(key) == 0 {
		return errors.New("Empty key specified")
	}
	return os.RemoveAll(filepath.Join(rootDir, key))
}
