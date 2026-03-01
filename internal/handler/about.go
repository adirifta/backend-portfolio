package handler

import (
	"backend-portfolio/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAbout returns the first about entry (public).
func (h *Handler) GetAbout(c *gin.Context) {
	about, err := h.abouts.FindFirst()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "About information not found"})
		return
	}
	c.JSON(http.StatusOK, about)
}

// UpdateAbout updates an existing about entry by ID.
func (h *Handler) UpdateAbout(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		return
	}

	about, err := h.abouts.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "About information not found"})
		return
	}

	if err := c.ShouldBindJSON(about); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.abouts.Save(about); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update about"})
		return
	}
	c.JSON(http.StatusOK, about)
}

// CreateOrUpdateAbout creates a new about entry or updates the existing one.
func (h *Handler) CreateOrUpdateAbout(c *gin.Context) {
	var about models.About
	if err := c.ShouldBindJSON(&about); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing, err := h.abouts.FindFirst()
	if err != nil {
		// No existing entry — create new
		if err := h.abouts.Create(&about); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create about"})
			return
		}
		c.JSON(http.StatusCreated, about)
		return
	}

	// Update existing entry
	about.ID = existing.ID
	if err := h.abouts.Save(&about); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update about"})
		return
	}
	c.JSON(http.StatusOK, about)
}
