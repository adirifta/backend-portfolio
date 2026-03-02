// Package router provides a single SetupRouter function that wires all
// middleware, public routes, auth routes, and admin routes.
//
// Both the root main.go (dev) and cmd/api/main.go (prod) call this function
// so route definitions are never duplicated (DRY).
package router

import (
	"time"

	"backend-portfolio/internal/auth"
	"backend-portfolio/internal/handler"
	"backend-portfolio/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter creates and configures a Gin engine with all application routes.
func SetupRouter(
	h *handler.Handler,
	jwtSvc *auth.JWTService,
	csrfSvc *middleware.CSRFService,
	allowedOrigins []string,
) *gin.Engine {
	r := gin.Default()

	// ── Global Middleware ──────────────────────────────────────
	// CORS must be applied BEFORE other middleware so preflight
	// OPTIONS requests are answered immediately with the correct headers.
	r.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization",
			"Accept", "X-Requested-With", "Content-Disposition",
			"X-XSRF-TOKEN",
		},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(middleware.SecurityHeaders())

	// ── Static files ──────────────────────────────────────────
	r.Static("/uploads", "./uploads")

	// ── Public routes ─────────────────────────────────────────
	r.GET("/api/about", h.GetAbout)
	r.GET("/api/portfolio", h.GetAllPortfolio)
	r.GET("/api/portfolio/:id", h.GetPortfolio)
	r.GET("/api/skills", h.GetAllSkills)
	r.GET("/api/skills/:id", h.GetSkill)
	r.GET("/api/skills/category", h.GetSkillsByCategory)
	r.GET("/api/qualifications", h.GetAllQualifications)
	r.GET("/api/qualifications/:id", h.GetQualification)

	// ── Auth routes ───────────────────────────────────────────
	r.POST("/api/auth/login", h.Login)
	r.POST("/api/auth/logout", h.Logout)
	r.POST("/api/auth/refresh", h.RefreshToken)
	r.GET("/api/auth/me", middleware.AuthMiddleware(jwtSvc), h.GetMe)
	r.POST("/api/login", h.Login) // backward compatibility

	// ── Protected routes (admin only — cookie auth + CSRF) ───
	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(jwtSvc))
	admin.Use(csrfSvc.Middleware())
	{
		// User management
		admin.POST("/create-user", h.CreateUser)
		admin.POST("/reset-admin", h.ResetAdminPassword)

		// About
		admin.POST("/about", h.CreateOrUpdateAbout)
		admin.PUT("/about/:id", h.UpdateAbout)

		// Portfolio
		admin.POST("/portfolio", h.CreatePortfolio)
		admin.PUT("/portfolio/:id", h.UpdatePortfolio)
		admin.DELETE("/portfolio/:id", h.DeletePortfolio)
		admin.DELETE("/portfolio-media/:portfolio_id/:media_id", h.DeletePortfolioMedia)

		// Skills
		admin.POST("/skills", h.CreateSkill)
		admin.PUT("/skills/:id", h.UpdateSkill)
		admin.DELETE("/skills/:id", h.DeleteSkill)

		// Qualifications
		admin.POST("/qualifications", h.CreateQualification)
		admin.PUT("/qualifications/:id", h.UpdateQualification)
		admin.DELETE("/qualifications/:id", h.DeleteQualification)
	}

	// ── Health & Info ─────────────────────────────────────────
	r.GET("/health", h.HealthCheck)
	r.GET("/", h.Info)

	return r
}
