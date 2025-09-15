package handlers

import (
	"net/http"
	"backend-portfolio/database"
	"backend-portfolio/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllPortfolio(c *gin.Context) {
	var portfolios []models.Portfolio
	if err := database.GetDB().Find(&portfolios).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch portfolio items"})
		return
	}

	c.JSON(http.StatusOK, portfolios)
}

func GetPortfolio(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var portfolio models.Portfolio
	if err := database.GetDB().First(&portfolio, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio item not found"})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

func CreatePortfolio(c *gin.Context) {
	var portfolio models.Portfolio
	if err := c.ShouldBindJSON(&portfolio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.GetDB().Create(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create portfolio item"})
		return
	}

	c.JSON(http.StatusCreated, portfolio)
}

func UpdatePortfolio(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var portfolio models.Portfolio
	if err := database.GetDB().First(&portfolio, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio item not found"})
		return
	}

	if err := c.ShouldBindJSON(&portfolio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.GetDB().Save(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update portfolio item"})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

func DeletePortfolio(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := database.GetDB().Delete(&models.Portfolio{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete portfolio item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Portfolio item deleted successfully"})
}