package service

import (
	"user/pkg/model"
	"user/pkg/repository"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

// Get user profile by ID
func (s *UserService) GetProfile(id uint) (*model.User, error) {
	return s.UserRepo.FindByID(id)
}

// Update user profile
func (s *UserService) UpdateProfile(user *model.User) error {
	return s.UserRepo.Update(user)
}

// ListByRole returns all users with the given role (e.g., 'doctor')
func (s *UserService) ListByRole(role string) ([]*model.User, error) {
	return s.UserRepo.FindByRole(role)
}
