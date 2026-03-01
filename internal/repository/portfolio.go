package repository

import (
	"backend-portfolio/models"

	"gorm.io/gorm"
)

type portfolioRepo struct{ db *gorm.DB }

// NewPortfolioRepository returns a GORM-backed PortfolioRepository.
func NewPortfolioRepository(db *gorm.DB) PortfolioRepository {
	return &portfolioRepo{db: db}
}

func (r *portfolioRepo) FindAll() ([]models.Portfolio, error) {
	var items []models.Portfolio
	if err := r.db.Preload("MediaFiles").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *portfolioRepo) FindByID(id uint) (*models.Portfolio, error) {
	var item models.Portfolio
	if err := r.db.Preload("MediaFiles").First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *portfolioRepo) Create(portfolio *models.Portfolio) error {
	return r.db.Create(portfolio).Error
}

func (r *portfolioRepo) Save(portfolio *models.Portfolio) error {
	return r.db.Save(portfolio).Error
}

func (r *portfolioRepo) Delete(id uint) error {
	return r.db.Delete(&models.Portfolio{}, id).Error
}

func (r *portfolioRepo) CreateMedia(media *models.PortfolioMedia) error {
	return r.db.Create(media).Error
}

func (r *portfolioRepo) FindMedia(portfolioID, mediaID uint) (*models.PortfolioMedia, error) {
	var media models.PortfolioMedia
	if err := r.db.Where("portfolio_id = ? AND id = ?", portfolioID, mediaID).First(&media).Error; err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *portfolioRepo) DeleteMedia(media *models.PortfolioMedia) error {
	return r.db.Delete(media).Error
}

// WithTransaction runs fn inside a database transaction.
// The fn receives a new PortfolioRepository backed by the transaction.
func (r *portfolioRepo) WithTransaction(fn func(tx PortfolioRepository) error) error {
	return r.db.Transaction(func(gormTx *gorm.DB) error {
		txRepo := &portfolioRepo{db: gormTx}
		return fn(txRepo)
	})
}
