package repository

import (
	"backend-portfolio/models"

	"gorm.io/gorm"
)

type aboutRepo struct{ db *gorm.DB }

// NewAboutRepository returns a GORM-backed AboutRepository.
func NewAboutRepository(db *gorm.DB) AboutRepository {
	return &aboutRepo{db: db}
}

func (r *aboutRepo) FindFirst() (*models.About, error) {
	var about models.About
	if err := r.db.First(&about).Error; err != nil {
		return nil, err
	}
	return &about, nil
}

func (r *aboutRepo) FindByID(id uint) (*models.About, error) {
	var about models.About
	if err := r.db.First(&about, id).Error; err != nil {
		return nil, err
	}
	return &about, nil
}

func (r *aboutRepo) Create(about *models.About) error {
	return r.db.Create(about).Error
}

func (r *aboutRepo) Save(about *models.About) error {
	return r.db.Save(about).Error
}
