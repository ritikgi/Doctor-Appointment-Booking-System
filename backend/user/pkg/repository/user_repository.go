package repository

import (
	"user/pkg/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	return r.DB.Save(user).Error
}

// FindByRole returns all users with the given role (e.g., 'doctor')
func (r *UserRepository) FindByRole(role string) ([]*model.User, error) {
	var users []*model.User
	err := r.DB.Where("role = ?", role).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
