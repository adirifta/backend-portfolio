package handler

import (
	"backend-portfolio/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllQualifications returns all qualifications.
func (h *Handler) GetAllQualifications(c *gin.Context) {
	qualifications, err := h.qualifications.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch qualifications"})
		return
	}
	c.JSON(http.StatusOK, qualifications)
}

// GetQualification returns a single qualification by ID.
func (h *Handler) GetQualification(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		return
	}

	qualification, err := h.qualifications.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Qualification not found"})
		return
	}
	c.JSON(http.StatusOK, qualification)
}

// CreateQualification creates a new qualification.
func (h *Handler) CreateQualification(c *gin.Context) {
	var qualification models.Qualification
	if err := c.ShouldBindJSON(&qualification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.qualifications.Create(&qualification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create qualification"})
		return
	}
	c.JSON(http.StatusCreated, qualification)
}

// UpdateQualification updates an existing qualification.
func (h *Handler) UpdateQualification(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		return
	}

	qualification, err := h.qualifications.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Qualification not found"})
		return
	}

	if err := c.ShouldBindJSON(qualification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.qualifications.Save(qualification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update qualification"})
		return
	}
	c.JSON(http.StatusOK, qualification)
}

// DeleteQualification removes a qualification by ID.
func (h *Handler) DeleteQualification(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		return
	}

	if err := h.qualifications.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete qualification"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Qualification deleted successfully"})
}
