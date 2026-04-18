package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders returns a Gin middleware that sets common security response headers.
// Includes CSP for XSS protection.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("X-XSS-Protection", "1; mode=block")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		// Content Security Policy for XSS protection
		// Restricts inline scripts and external resource loading
		h.Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'")

		// Prevent caching of API responses
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			h.Set("Cache-Control", "no-store, no-cache, must-revalidate")
			h.Set("Pragma", "no-cache")
		}

		c.Next()
	}
}
