package handlers

import (
	"net/http"
	"backend-portfolio/database"
	"backend-portfolio/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllQualifications(c *gin.Context) {
	var qualifications []models.Qualification
	if err := database.GetDB().Find(&qualifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch qualifications"})
		return
	}

	c.JSON(http.StatusOK, qualifications)
}

func GetQualification(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var qualification models.Qualification
	if err := database.GetDB().First(&qualification, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Qualification not found"})
		return
	}

	c.JSON(http.StatusOK, qualification)
}

func CreateQualification(c *gin.Context) {
	var qualification models.Qualification
	if err := c.ShouldBindJSON(&qualification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.GetDB().Create(&qualification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create qualification"})
		return
	}

	c.JSON(http.StatusCreated, qualification)
}

func UpdateQualification(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var qualification models.Qualification
	if err := database.GetDB().First(&qualification, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Qualification not found"})
		return
	}

	if err := c.ShouldBindJSON(&qualification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.GetDB().Save(&qualification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update qualification"})
		return
	}

	c.JSON(http.StatusOK, qualification)
}

func DeleteQualification(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := database.GetDB().Delete(&models.Qualification{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete qualification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Qualification deleted successfully"})
}