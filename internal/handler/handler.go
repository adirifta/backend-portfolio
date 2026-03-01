// Package handler provides HTTP handlers with injected dependencies.
//
// All handlers are methods on the Handler struct, which receives its
// dependencies (repositories, auth services) via the constructor.
// This follows the Dependency Inversion Principle: handlers depend on
// repository interfaces, not concrete GORM implementations.
package handler

import (
	"net/http"
	"strconv"

	"backend-portfolio/internal/auth"
	"backend-portfolio/internal/middleware"
	"backend-portfolio/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler holds all dependencies required by HTTP handlers.
type Handler struct {
	db             *gorm.DB
	users          repository.UserRepository
	abouts         repository.AboutRepository
	portfolios     repository.PortfolioRepository
	skills         repository.SkillRepository
	qualifications repository.QualificationRepository

	jwt    *auth.JWTService
	cookie *auth.CookieManager
	csrf   *middleware.CSRFService
}

// New creates a Handler with all required dependencies.
func New(
	db *gorm.DB,
	users repository.UserRepository,
	abouts repository.AboutRepository,
	portfolios repository.PortfolioRepository,
	skills repository.SkillRepository,
	qualifications repository.QualificationRepository,
	jwt *auth.JWTService,
	cookie *auth.CookieManager,
	csrf *middleware.CSRFService,
) *Handler {
	return &Handler{
		db:             db,
		users:          users,
		abouts:         abouts,
		portfolios:     portfolios,
		skills:         skills,
		qualifications: qualifications,
		jwt:            jwt,
		cookie:         cookie,
		csrf:           csrf,
	}
}

// parseIDParam extracts and validates an unsigned integer ":id" path parameter.
func parseIDParam(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return 0, err
	}
	return uint(id), nil
}

// parseNamedIDParam extracts and validates an unsigned integer path parameter by name.
func parseNamedIDParam(c *gin.Context, name string) (uint, error) {
	id, err := strconv.ParseUint(c.Param(name), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid " + name})
		return 0, err
	}
	return uint(id), nil
}
