package handlers

import (
	"backend-portfolio/database"
	"backend-portfolio/models"
	"backend-portfolio/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.GetDB().Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// Handler untuk membuat admin user pertama kali (hanya untuk setup)
func CreateAdminUser(c *gin.Context) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	adminUser := models.User{
		Username: "admin",
		Password: string(hashedPassword),
		Role:     "admin",
	}

	if err := database.GetDB().Create(&adminUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create admin user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Admin user created successfully"})
}

// CreateUser handler untuk membuat user baru
func CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	if err := database.GetDB().Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// ResetAdminPassword handler untuk reset password admin
func ResetAdminPassword(c *gin.Context) {
	// Hanya boleh diakses di development atau dengan secret key
	if c.GetHeader("X-Reset-Secret") != "dev-reset-2024" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)

	// Update password admin
	if err := database.GetDB().Model(&models.User{}).
		Where("username = ?", "admin").
		Update("password", string(hashedPassword)).
		Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin password reset successfully to 'admin123'"})
}

// Tambahkan function UpdateAbout
func UpdateAbout(c *gin.Context) {
    id := c.Param("id")
    
    var about models.About
    if err := database.GetDB().First(&about, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "About information not found"})
        return
    }

    if err := c.ShouldBindJSON(&about); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := database.GetDB().Save(&about).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update about"})
        return
    }

    c.JSON(http.StatusOK, about)
}