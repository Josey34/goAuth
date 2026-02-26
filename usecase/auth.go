package usecase

import (
	"goauth/entity"
	"goauth/errors"
	"goauth/repository"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(email, password, name string) (*entity.User, error)
}

type AuthUsecaseImpl struct {
	userRepo   repository.UserRepository
	bcryptCost int
}

func NewAuthUsecase(userRepo repository.UserRepository, bcryptCost int) AuthUsecase {
	return &AuthUsecaseImpl{
		userRepo:   userRepo,
		bcryptCost: bcryptCost,
	}
}

func (a *AuthUsecaseImpl) Register(email, password, name string) (*entity.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" || password == "" || name == "" {
		return nil, errors.ValidationError{Message: "Email, password and name are required"}
	}

	if _, err := a.userRepo.FindByEmail(email); err == nil {
		return nil, errors.ConflictError{Message: "Email already exists"}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), a.bcryptCost)
	if err != nil {
		return nil, errors.ValidationError{Message: "Failed to hash password"}
	}

	createdUser, err := a.userRepo.Insert(email, name, string(passwordHash))
	if err != nil {
		return nil, errors.AuthError{Message: "Failed to create user"}
	}

	return createdUser, nil
}
