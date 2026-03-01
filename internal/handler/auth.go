package handler

import (
	"backend-portfolio/internal/auth"
	"backend-portfolio/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login authenticates a user and sets HTTP-only cookies with access + refresh tokens.
func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.users.FindByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, err := h.jwt.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
		return
	}

	refreshToken, err := h.jwt.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate refresh token"})
		return
	}

	h.cookie.SetTokenCookies(c, accessToken, refreshToken)
	h.csrf.SetCSRFCookie(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// Logout clears the authentication and CSRF cookies.
func (h *Handler) Logout(c *gin.Context) {
	h.cookie.ClearTokenCookies(c)
	h.csrf.ClearCSRFCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// RefreshToken validates the refresh cookie and issues a new token pair (rotation).
func (h *Handler) RefreshToken(c *gin.Context) {
	cookie, err := c.Cookie(auth.RefreshTokenCookieName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found"})
		return
	}

	claims, err := h.jwt.ValidateRefreshToken(cookie)
	if err != nil {
		h.cookie.ClearTokenCookies(c)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Verify the user still exists
	user, err := h.users.FindByID(claims.UserID)
	if err != nil {
		h.cookie.ClearTokenCookies(c)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User no longer exists"})
		return
	}

	newAccess, err := h.jwt.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate access token"})
		return
	}

	newRefresh, err := h.jwt.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate refresh token"})
		return
	}

	h.cookie.SetTokenCookies(c, newAccess, newRefresh)
	h.csrf.SetCSRFCookie(c)

	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed successfully"})
}

// GetMe returns the currently authenticated user info.
func (h *Handler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
		return
	}

	user, err := h.users.FindByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// CreateUser creates a new user (admin only).
func (h *Handler) CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Role == "" {
		req.Role = "admin"
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := &models.User{Username: req.Username, Password: string(hashed), Role: req.Role}
	if err := h.users.Create(user); err != nil {
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

// ResetAdminPassword resets the admin password (requires X-Reset-Secret header).
func (h *Handler) ResetAdminPassword(c *gin.Context) {
	resetSecret := os.Getenv("RESET_SECRET")
	if resetSecret == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Reset not configured"})
		return
	}

	if c.GetHeader("X-Reset-Secret") != resetSecret {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed"})
		return
	}

	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "new_password is required (min 8 chars)"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err := h.users.UpdatePassword("admin", string(hashed)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin password reset successfully"})
}
