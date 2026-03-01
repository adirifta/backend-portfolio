package handler

import (
	"backend-portfolio/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Valid skill levels.
var validLevels = map[string]bool{
	"Beginner":     true,
	"Intermediate": true,
	"Advanced":     true,
	"Expert":       true,
}

// GetAllSkills returns all skills.
func (h *Handler) GetAllSkills(c *gin.Context) {
	skills, err := h.skills.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch skills"})
		return
	}
	c.JSON(http.StatusOK, skills)
}

// GetSkill returns a single skill by ID.
func (h *Handler) GetSkill(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		return
	}

	skill, err := h.skills.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}
	c.JSON(http.StatusOK, skill)
}

// CreateSkill creates a new skill with validation.
func (h *Handler) CreateSkill(c *gin.Context) {
	var skill models.Skill
	if err := c.ShouldBindJSON(&skill); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if skill.Score < 0 || skill.Score > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Score must be between 0 and 100"})
		return
	}
	if !validLevels[skill.Level] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Level must be one of: Beginner, Intermediate, Advanced, Expert"})
		return
	}

	if err := h.skills.Create(&skill); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create skill"})
		return
	}
	c.JSON(http.StatusCreated, skill)
}

// UpdateSkill updates an existing skill.
func (h *Handler) UpdateSkill(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		return
	}

	skill, err := h.skills.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	var update models.Skill
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate score if provided
	if update.Score != 0 && (update.Score < 0 || update.Score > 100) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Score must be between 0 and 100"})
		return
	}

	// Validate level if provided
	if update.Level != "" && !validLevels[update.Level] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Level must be one of: Beginner, Intermediate, Advanced, Expert"})
		return
	}

	// Apply non-zero fields
	if update.Name != "" {
		skill.Name = update.Name
	}
	if update.Level != "" {
		skill.Level = update.Level
	}
	if update.Score != 0 {
		skill.Score = update.Score
	}
	if update.Category != "" {
		skill.Category = update.Category
	}
	if update.Icon != "" {
		skill.Icon = update.Icon
	}

	if err := h.skills.Save(skill); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update skill"})
		return
	}
	c.JSON(http.StatusOK, skill)
}

// DeleteSkill removes a skill by ID.
func (h *Handler) DeleteSkill(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		return
	}

	if err := h.skills.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete skill"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Skill deleted successfully"})
}

// GetSkillsByCategory returns skills grouped by category.
func (h *Handler) GetSkillsByCategory(c *gin.Context) {
	grouped, err := h.skills.FindAllGroupedByCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch skills"})
		return
	}
	c.JSON(http.StatusOK, grouped)
}
