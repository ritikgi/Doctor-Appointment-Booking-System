package service

import (
	"auth/pkg/model"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	CreateFunc      func(user *model.User) error
	FindByEmailFunc func(email string) (*model.User, error)
}

func (m *mockUserRepo) Create(user *model.User) error {
	return m.CreateFunc(user)
}
func (m *mockUserRepo) FindByEmail(email string) (*model.User, error) {
	return m.FindByEmailFunc(email)
}

func TestRegister_Success(t *testing.T) {
	svc := &AuthService{UserRepo: &mockUserRepo{
		CreateFunc: func(user *model.User) error {
			if user.Email == "test@example.com" && user.Role == "doctor" {
				return nil
			}
			return errors.New("unexpected user")
		},
	}}
	err := svc.Register("Test", "test@example.com", "password", "doctor")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRegister_Fail(t *testing.T) {
	svc := &AuthService{UserRepo: &mockUserRepo{
		CreateFunc: func(user *model.User) error {
			return errors.New("db error")
		},
	}}
	err := svc.Register("Test", "test@example.com", "password", "doctor")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLogin_Success(t *testing.T) {
	// password is "password"
	hash, _ := svcHash("password")
	svc := &AuthService{UserRepo: &mockUserRepo{
		FindByEmailFunc: func(email string) (*model.User, error) {
			return &model.User{ID: 1, Email: email, PasswordHash: hash, Role: "doctor"}, nil
		},
	}}
	// Set JWT_SECRET for test
	t.Setenv("JWT_SECRET", "testsecret")
	_, err := svc.Login("test@example.com", "password")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestLogin_Fail_WrongPassword(t *testing.T) {
	hash, _ := svcHash("password")
	svc := &AuthService{UserRepo: &mockUserRepo{
		FindByEmailFunc: func(email string) (*model.User, error) {
			return &model.User{ID: 1, Email: email, PasswordHash: hash, Role: "doctor"}, nil
		},
	}}
	_, err := svc.Login("test@example.com", "wrongpass")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLogin_Fail_UserNotFound(t *testing.T) {
	svc := &AuthService{UserRepo: &mockUserRepo{
		FindByEmailFunc: func(email string) (*model.User, error) {
			return nil, errors.New("not found")
		},
	}}
	_, err := svc.Login("notfound@example.com", "password")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// helper for password hash
func svcHash(pw string) (string, error) {
	cost := 10
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), cost)
	return string(hash), err
}
