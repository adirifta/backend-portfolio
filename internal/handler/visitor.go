package handler

import (
	"backend-portfolio/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// TrackVisitor records a visitor hit (public — called by React frontend).
//
// The React frontend should call POST /api/visitors/track on every page load
// with a JSON body: { "path": "/about", "referer": "https://google.com" }
// If no body is sent, the path defaults to "/".
func (h *Handler) TrackVisitor(c *gin.Context) {
	var body struct {
		Path    string `json:"path"`
		Referer string `json:"referer"`
	}
	_ = c.ShouldBindJSON(&body) // silently ignore bind errors
	if body.Path == "" {
		body.Path = "/"
	}

	// Prefer referer from body (sent by React via document.referrer),
	// fall back to HTTP Referer header.
	referer := body.Referer
	if referer == "" {
		referer = c.Request.Referer()
	}

	// GeoIP lookup for country & city
	ip := c.ClientIP()
	var country, city string
	if h.geoip != nil {
		loc := h.geoip.Lookup(ip)
		country = loc.Country
		city = loc.City
	}

	visitor := &models.Visitor{
		IPAddress: ip,
		UserAgent: c.Request.UserAgent(),
		Path:      body.Path,
		Referer:   referer,
		Country:   country,
		City:      city,
		VisitedAt: time.Now(),
	}

	if err := h.visitors.Record(visitor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record visitor"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Visitor recorded",
	})
}

// GetVisitorStats returns aggregated visitor statistics (admin).
func (h *Handler) GetVisitorStats(c *gin.Context) {
	stats, err := h.visitors.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get visitor stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// GetVisitorStatsPublic returns a limited set of visitor statistics (public).
// Only exposes total and unique counts — no IP addresses or user agents.
func (h *Handler) GetVisitorStatsPublic(c *gin.Context) {
	stats, err := h.visitors.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get visitor stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"today":             stats.Today,
		"this_week":         stats.ThisWeek,
		"this_month":        stats.ThisMonth,
		"this_year":         stats.ThisYear,
		"total":             stats.Total,
		"unique_today":      stats.UniqueToday,
		"unique_this_week":  stats.UniqueWeek,
		"unique_this_month": stats.UniqueMonth,
		"unique_this_year":  stats.UniqueYear,
		"unique_total":      stats.UniqueTotal,
		"daily_chart":       stats.DailyChart,
		"monthly_chart":     stats.MonthlyChart,
		"top_pages":         stats.TopPages,
	})
}
