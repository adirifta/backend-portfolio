package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck verifies the database connection is alive.
func (h *Handler) HealthCheck(c *gin.Context) {
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERROR", "message": "Database connection failed"})
		return
	}
	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "ERROR", "message": "Database ping failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Server is running"})
}

// Info returns basic API info.
func (h *Handler) Info(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Backend Portfolio API is running!",
		"version": "2.0.0",
	})
}
