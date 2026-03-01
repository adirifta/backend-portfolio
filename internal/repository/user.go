package repository

import (
	"backend-portfolio/models"

	"gorm.io/gorm"
)

type userRepo struct{ db *gorm.DB }

// NewUserRepository returns a GORM-backed UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) UpdatePassword(username string, hashedPassword string) error {
	return r.db.Model(&models.User{}).
		Where("username = ?", username).
		Update("password", hashedPassword).Error
}
