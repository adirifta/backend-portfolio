package middleware

import "github.com/gin-gonic/gin"

// SecurityHeaders returns a Gin middleware that sets common security response headers.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("X-XSS-Protection", "1; mode=block")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

		// Prevent caching of API responses
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			h.Set("Cache-Control", "no-store, no-cache, must-revalidate")
			h.Set("Pragma", "no-cache")
		}

		c.Next()
	}
}
