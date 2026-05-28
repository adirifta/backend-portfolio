package handler

import (
	"backend-portfolio/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetDashboardStats returns total counts for dashboard summary cards.
func (h *Handler) GetDashboardStats(c *gin.Context) {
	var totalSkills int64
	if err := h.db.Model(&models.Skill{}).Count(&totalSkills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count skills"})
		return
	}

	var totalProjects int64
	if err := h.db.Model(&models.Portfolio{}).Count(&totalProjects).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count projects"})
		return
	}

	var totalQualifications int64
	if err := h.db.Model(&models.Qualification{}).Count(&totalQualifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count qualifications"})
		return
	}

	var totalVisitors int64
	if err := h.db.Model(&models.Visitor{}).Count(&totalVisitors).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count visitors"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_skill":          totalSkills,
		"total_project":        totalProjects,
		"total_qualifications": totalQualifications,
		"total_visitors":       totalVisitors,
	})
}
