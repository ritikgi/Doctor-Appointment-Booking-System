package service

import (
	"auth/pkg/model"
	"auth/pkg/util"
	"os"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
}

type AuthService struct {
	UserRepo UserRepo
}

// Register a new user (hashes password)
func (s *AuthService) Register(name, email, password, role string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := &model.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
	}
	return s.UserRepo.Create(user)
}

// Authenticate user and return JWT if successful
func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.UserRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", err
	}
	// Generate JWT token
	return util.GenerateJWT(user.ID, user.Email, user.Role, os.Getenv("JWT_SECRET"))
}
