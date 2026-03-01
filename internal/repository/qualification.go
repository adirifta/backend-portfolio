package repository

import (
	"backend-portfolio/models"

	"gorm.io/gorm"
)

type qualificationRepo struct{ db *gorm.DB }

// NewQualificationRepository returns a GORM-backed QualificationRepository.
func NewQualificationRepository(db *gorm.DB) QualificationRepository {
	return &qualificationRepo{db: db}
}

func (r *qualificationRepo) FindAll() ([]models.Qualification, error) {
	var items []models.Qualification
	if err := r.db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *qualificationRepo) FindByID(id uint) (*models.Qualification, error) {
	var item models.Qualification
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *qualificationRepo) Create(q *models.Qualification) error {
	return r.db.Create(q).Error
}

func (r *qualificationRepo) Save(q *models.Qualification) error {
	return r.db.Save(q).Error
}

func (r *qualificationRepo) Delete(id uint) error {
	return r.db.Delete(&models.Qualification{}, id).Error
}
