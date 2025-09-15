package handlers

import (
	"backend-portfolio/database"
	"backend-portfolio/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAbout(c *gin.Context) {
	var about models.About
	if err := database.GetDB().First(&about).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "About information not found"})
		return
	}

	c.JSON(http.StatusOK, about)
}

func CreateOrUpdateAbout(c *gin.Context) {
	var about models.About
	if err := c.ShouldBindJSON(&about); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if about entry exists
	var existingAbout models.About
	if err := database.GetDB().First(&existingAbout).Error; err != nil {
		// Create new about entry
		if err := database.GetDB().Create(&about).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create about"})
			return
		}
		c.JSON(http.StatusCreated, about)
	} else {
		// Update existing about entry
		about.ID = existingAbout.ID
		if err := database.GetDB().Save(&about).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update about"})
			return
		}
		c.JSON(http.StatusOK, about)
	}
}