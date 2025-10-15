package main

import (
	"backend-portfolio/config"
	"backend-portfolio/database"
	"backend-portfolio/handlers"
	"backend-portfolio/middleware"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Set production mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	database.InitDB(cfg)

	// Setup router
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://backend-portofolio.adirdk.com", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Public routes
	r.POST("/api/login", handlers.Login)
	r.POST("/api/reset-admin", handlers.ResetAdminPassword)
	r.POST("/api/create-user", handlers.CreateUser)
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
		auth.PUT("/about/:id", handlers.UpdateAbout)
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

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"message": "Server is running",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
		if port == "" {
			port = "8080"
		}
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
