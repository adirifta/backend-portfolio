package repository

import (
	"backend-portfolio/models"

	"gorm.io/gorm"
)

type skillRepo struct{ db *gorm.DB }

// NewSkillRepository returns a GORM-backed SkillRepository.
func NewSkillRepository(db *gorm.DB) SkillRepository {
	return &skillRepo{db: db}
}

func (r *skillRepo) FindAll() ([]models.Skill, error) {
	var items []models.Skill
	if err := r.db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *skillRepo) FindByID(id uint) (*models.Skill, error) {
	var item models.Skill
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *skillRepo) Create(skill *models.Skill) error {
	return r.db.Create(skill).Error
}

func (r *skillRepo) Save(skill *models.Skill) error {
	return r.db.Save(skill).Error
}

func (r *skillRepo) Delete(id uint) error {
	return r.db.Delete(&models.Skill{}, id).Error
}

func (r *skillRepo) FindAllGroupedByCategory() (map[string][]models.Skill, error) {
	var skills []models.Skill
	if err := r.db.Order("category, score DESC").Find(&skills).Error; err != nil {
		return nil, err
	}

	grouped := make(map[string][]models.Skill)
	for _, s := range skills {
		grouped[s.Category] = append(grouped[s.Category], s)
	}
	return grouped, nil
}
