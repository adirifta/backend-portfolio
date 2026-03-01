// Package middleware provides HTTP middleware for the Gin router.
package middleware

import (
	"backend-portfolio/internal/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware returns a Gin middleware that validates the access token
// from the HTTP-only cookie (preferred) or Authorization header (fallback).
// Only "admin" role is allowed through.
func AuthMiddleware(jwt *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. Prefer HTTP-only cookie
		if cookie, err := c.Cookie(auth.AccessTokenCookieName); err == nil && cookie != "" {
			tokenString = cookie
		}

		// 2. Fallback: Authorization Bearer header (for non-browser API clients)
		if tokenString == "" {
			if h := c.GetHeader("Authorization"); h != "" {
				tokenString = strings.TrimPrefix(h, "Bearer ")
			}
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		claims, err := jwt.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if claims.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// OptionalAuthMiddleware extracts user info from the token if present,
// but does not block the request if absent.
func OptionalAuthMiddleware(jwt *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		if cookie, err := c.Cookie(auth.AccessTokenCookieName); err == nil && cookie != "" {
			tokenString = cookie
		}
		if tokenString == "" {
			if h := c.GetHeader("Authorization"); h != "" {
				tokenString = strings.TrimPrefix(h, "Bearer ")
			}
		}

		if tokenString != "" {
			if claims, err := jwt.ValidateAccessToken(tokenString); err == nil {
				c.Set("user_id", claims.UserID)
				c.Set("role", claims.Role)
			}
		}

		c.Next()
	}
}
