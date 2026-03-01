package main

import (
	"log"
	"os"

	"backend-portfolio/config"
	"backend-portfolio/database"
	"backend-portfolio/internal/auth"
	"backend-portfolio/internal/handler"
	"backend-portfolio/internal/middleware"
	"backend-portfolio/internal/repository"
	"backend-portfolio/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {
	// Production mode when running on Cloud Run
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Println("🔧 Loading configuration...")
	cfg := config.LoadConfig()

	log.Println("🗄️ Initializing database...")
	database.InitDB(cfg)
	db := database.GetDB()

	// Services
	jwtSvc := auth.NewJWTService(
		cfg.JWTSecret, cfg.JWTRefreshSecret,
		cfg.AccessTokenExpiry, cfg.RefreshTokenExpiry,
	)
	cookieMgr := auth.NewCookieManager(
		cfg.CookieDomain, cfg.CookieSecure, cfg.CookieSameSite,
		cfg.AccessTokenExpiry, cfg.RefreshTokenExpiry,
	)
	csrfSvc := middleware.NewCSRFService(jwtSvc.AccessSecret(), cookieMgr)

	// Repositories
	users := repository.NewUserRepository(db)
	abouts := repository.NewAboutRepository(db)
	portfolios := repository.NewPortfolioRepository(db)
	skills := repository.NewSkillRepository(db)
	qualifications := repository.NewQualificationRepository(db)

	// Handler (all dependencies injected)
	h := handler.New(db, users, abouts, portfolios, skills, qualifications, jwtSvc, cookieMgr, csrfSvc)

	// Router
	log.Println("🚀 Setting up router...")
	r := router.SetupRouter(h, jwtSvc, csrfSvc, cfg.AllowedOrigins)

	// Start server
	port := cfg.Port
	if os.Getenv("K_SERVICE") != "" {
		log.Println("🌐 Running on Cloud Run environment 🚀")
	} else {
		log.Println("💻 Running locally")
	}

	log.Printf("🎯 Server starting on port %s", port)
	log.Printf("📍 Health check: http://localhost:%s/health", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
