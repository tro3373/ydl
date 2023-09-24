package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tro3373/ydl/cmd/api/response"
	"go.uber.org/zap"
)

func (h *Handler) ok(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}

func (h *Handler) list(c *gin.Context, list []response.Res) {
	c.JSON(http.StatusOK, gin.H{
		"list": list,
	})
}

func (h *Handler) handleAsServerError(c *gin.Context, message string, err error) {
	h.logger.Error(message, zap.String("Error:", err.Error()))
	c.JSON(http.StatusInternalServerError, gin.H{"status": "StatusInternalServerError"})
}

func (h *Handler) handleAsBadRequest(c *gin.Context, message string, err error) {
	if err != nil {
		h.logger.Error(message, zap.String("Error:", err.Error()))
	}
	c.JSON(http.StatusBadRequest, gin.H{"status": "BadRequest"})
}
