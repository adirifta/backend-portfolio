package main

import (
	"backend-portfolio/config"
	"backend-portfolio/database"
	"backend-portfolio/handlers"
	"backend-portfolio/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	database.InitDB(cfg)

	// Setup router
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://adirdk.com", "http://localhost:3000"}, // Your React app address
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Public routes
	r.POST("/api/login", handlers.Login)
	r.POST("/api/reset-admin", handlers.ResetAdminPassword)
	r.POST("/api/create-user", handlers.CreateUser)

	// Public routes
	r.GET("/api/about", handlers.GetAbout)
	r.GET("/api/portfolio", handlers.GetAllPortfolio)
	r.GET("/api/portfolio/:id", handlers.GetPortfolio)
	r.GET("/api/skills", handlers.GetAllSkills)
	r.GET("/api/skills/:id", handlers.GetSkill)
	r.GET("/api/qualifications", handlers.GetAllQualifications)
	r.GET("/api/qualifications/:id", handlers.GetQualification)

	// Protected routes (admin only)
	auth := r.Group("/api/admin")
	auth.Use(middleware.AuthMiddleware())
	{
		auth.POST("/about", handlers.CreateOrUpdateAbout)
		auth.POST("/portfolio", handlers.CreatePortfolio)
		auth.PUT("/portfolio/:id", handlers.UpdatePortfolio)
		auth.DELETE("/portfolio/:id", handlers.DeletePortfolio)
		auth.POST("/skills", handlers.CreateSkill)
		auth.PUT("/skills/:id", handlers.UpdateSkill)
		auth.DELETE("/skills/:id", handlers.DeleteSkill)
		auth.POST("/qualifications", handlers.CreateQualification)
		auth.PUT("/qualifications/:id", handlers.UpdateQualification)
		auth.DELETE("/qualifications/:id", handlers.DeleteQualification)
	}

	// Start server
	r.Run(":" + cfg.Port)
}
