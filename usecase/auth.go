package usecase

import (
	"fmt"
	"goauth/entity"
	"goauth/errors"
	"goauth/repository"
	"goauth/service"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(email, password, name string) (*entity.User, error)
	Login(email, password string) (*entity.User, string, string, error)
	Refresh(refreshToken string) (string, error)
}

type AuthUsecaseImpl struct {
	userRepo     repository.UserRepository
	bcryptCost   int
	tokenService service.TokenService
}

func NewAuthUsecase(userRepo repository.UserRepository, bcryptCost int, tokenService service.TokenService) AuthUsecase {
	return &AuthUsecaseImpl{
		userRepo:     userRepo,
		bcryptCost:   bcryptCost,
		tokenService: tokenService,
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

func (a *AuthUsecaseImpl) Login(email, password string) (*entity.User, string, string, error) {
	trimEmail := strings.TrimSpace(email)
	lowerEmail := strings.ToLower(trimEmail)

	userFound, err := a.userRepo.FindByEmail(lowerEmail)
	if err != nil {
		return nil, "", "", errors.AuthError{Message: "Invalid email or password"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.PasswordHash), []byte(password)); err != nil {
		return nil, "", "", errors.AuthError{Message: "Invalid email or password"}
	}

	accessToken, err := a.tokenService.GenerateAccess(userFound.ID, userFound.Role)
	if err != nil {
		return nil, "", "", errors.AuthError{Message: "Failed to generate access token"}
	}

	refreshToken, err := a.tokenService.GenerateRefresh(userFound.ID, userFound.Role)
	if err != nil {
		return nil, "", "", errors.AuthError{Message: "Failed to generate refresh token"}
	}

	return userFound, accessToken, refreshToken, nil
}

func (a *AuthUsecaseImpl) Refresh(refreshToken string) (string, error) {
	claims, err := a.tokenService.Validate(refreshToken)
	if err != nil {
		return "", errors.AuthError{Message: "Invalid refresh token"}
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return "", errors.AuthError{Message: "Invalid token type"}
	}

	userID := fmt.Sprintf("%v", claims["sub"])
	role := fmt.Sprintf("%v", claims["role"])

	accessToken, err := a.tokenService.GenerateAccess(userID, role)
	if err != nil {
		return "", errors.AuthError{Message: "Failed to generate access token"}
	}

	return accessToken, nil
}
