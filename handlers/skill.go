package handlers

import (
	"net/http"
	"backend-portfolio/database"
	"backend-portfolio/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllSkills(c *gin.Context) {
	var skills []models.Skill
	if err := database.GetDB().Find(&skills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch skills"})
		return
	}

	c.JSON(http.StatusOK, skills)
}

func GetSkill(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var skill models.Skill
	if err := database.GetDB().First(&skill, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	c.JSON(http.StatusOK, skill)
}

func CreateSkill(c *gin.Context) {
	var skill models.Skill
	if err := c.ShouldBindJSON(&skill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi score
	if skill.Score < 0 || skill.Score > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Score must be between 0 and 100"})
		return
	}

	// Validasi level
	validLevels := map[string]bool{
		"Beginner":     true,
		"Intermediate": true,
		"Advanced":     true,
		"Expert":       true,
	}
	if !validLevels[skill.Level] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Level must be one of: Beginner, Intermediate, Advanced, Expert"})
		return
	}

	if err := database.GetDB().Create(&skill).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create skill"})
		return
	}

	c.JSON(http.StatusCreated, skill)
}

func UpdateSkill(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var skill models.Skill
	if err := database.GetDB().First(&skill, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	var updateData models.Skill
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi score jika diupdate
	if updateData.Score != 0 && (updateData.Score < 0 || updateData.Score > 100) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Score must be between 0 and 100"})
		return
	}

	// Validasi level jika diupdate
	if updateData.Level != "" {
		validLevels := map[string]bool{
			"Beginner":     true,
			"Intermediate": true,
			"Advanced":     true,
			"Expert":       true,
		}
		if !validLevels[updateData.Level] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Level must be one of: Beginner, Intermediate, Advanced, Expert"})
			return
		}
	}

	// Update fields
	if updateData.Name != "" {
		skill.Name = updateData.Name
	}
	if updateData.Level != "" {
		skill.Level = updateData.Level
	}
	if updateData.Score != 0 {
		skill.Score = updateData.Score
	}
	if updateData.Category != "" {
		skill.Category = updateData.Category
	}
	if updateData.Icon != "" {
		skill.Icon = updateData.Icon
	}

	if err := database.GetDB().Save(&skill).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update skill"})
		return
	}

	c.JSON(http.StatusOK, skill)
}

func DeleteSkill(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := database.GetDB().Delete(&models.Skill{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete skill"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill deleted successfully"})
}

func GetSkillsByCategory(c *gin.Context) {
	var skills []models.Skill
	if err := database.GetDB().Order("category, score DESC").Find(&skills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch skills"})
		return
	}

	// Group by category
	skillsByCategory := make(map[string][]models.Skill)
	for _, skill := range skills {
		skillsByCategory[skill.Category] = append(skillsByCategory[skill.Category], skill)
	}

	c.JSON(http.StatusOK, skillsByCategory)
}