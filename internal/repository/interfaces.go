// Package repository defines interfaces for data access.
// Implementations use GORM but handlers depend only on these interfaces (DIP).
package repository

import "backend-portfolio/models"

// UserRepository abstracts user data access.
type UserRepository interface {
	FindByID(id uint) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	Create(user *models.User) error
	UpdatePassword(username string, hashedPassword string) error
}

// AboutRepository abstracts about data access.
type AboutRepository interface {
	FindFirst() (*models.About, error)
	FindByID(id uint) (*models.About, error)
	Create(about *models.About) error
	Save(about *models.About) error
}

// PortfolioRepository abstracts portfolio data access.
type PortfolioRepository interface {
	FindAll() ([]models.Portfolio, error)
	FindByID(id uint) (*models.Portfolio, error)
	Create(portfolio *models.Portfolio) error
	Save(portfolio *models.Portfolio) error
	Delete(id uint) error
	CreateMedia(media *models.PortfolioMedia) error
	FindMedia(portfolioID, mediaID uint) (*models.PortfolioMedia, error)
	DeleteMedia(media *models.PortfolioMedia) error
	// Transaction wraps the callback in a DB transaction.
	WithTransaction(fn func(tx PortfolioRepository) error) error
}

// SkillRepository abstracts skill data access.
type SkillRepository interface {
	FindAll() ([]models.Skill, error)
	FindByID(id uint) (*models.Skill, error)
	Create(skill *models.Skill) error
	Save(skill *models.Skill) error
	Delete(id uint) error
	FindAllGroupedByCategory() (map[string][]models.Skill, error)
}

// QualificationRepository abstracts qualification data access.
type QualificationRepository interface {
	FindAll() ([]models.Qualification, error)
	FindByID(id uint) (*models.Qualification, error)
	Create(qualification *models.Qualification) error
	Save(qualification *models.Qualification) error
	Delete(id uint) error
}
