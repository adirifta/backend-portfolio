// Package main provides a convenience entrypoint for local development.
// For production, use: go run ./cmd/api
//
// Usage:
//
//	go run main.go              (local dev)
//	go run ./cmd/api            (production entrypoint)
//	docker compose up -d        (Docker)
package main

import (
	"log"

	"backend-portfolio/config"
	"backend-portfolio/database"
	"backend-portfolio/internal/auth"
	"backend-portfolio/internal/geoip"
	"backend-portfolio/internal/handler"
	"backend-portfolio/internal/repository"
	"backend-portfolio/internal/router"
)

func main() {
	cfg := config.LoadConfig()

	// Database
	database.InitDB(cfg)
	db := database.GetDB()

	// Services
	jwtSvc := auth.NewJWTService(
		cfg.JWTSecret, cfg.JWTRefreshSecret,
		cfg.AccessTokenExpiry, cfg.RefreshTokenExpiry,
	)
	geoSvc := geoip.NewService()

	// Repositories
	users := repository.NewUserRepository(db)
	abouts := repository.NewAboutRepository(db)
	portfolios := repository.NewPortfolioRepository(db)
	skills := repository.NewSkillRepository(db)
	qualifications := repository.NewQualificationRepository(db)
	visitors := repository.NewVisitorRepository(db)

	// Handler (all dependencies injected)
	h := handler.New(db, users, abouts, portfolios, skills, qualifications, visitors, jwtSvc, geoSvc)

	// Router
	r := router.SetupRouter(h, jwtSvc, cfg.AllowedOrigins)

	log.Printf("🎯 Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
