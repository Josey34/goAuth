package usecase

import (
	"goauth/entity"
	"goauth/repository"
)

type UserUsecase interface {
	GetProfile(userID string) (*entity.User, error)
	UpdateProfile(userID, name string) (*entity.User, error)
}

type UserUsecaseImpl struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &UserUsecaseImpl{
		userRepo: userRepo,
	}
}

func (u *UserUsecaseImpl) GetProfile(userID string) (*entity.User, error) {
	return u.userRepo.FindByID(userID)
}

func (u *UserUsecaseImpl) UpdateProfile(userID, name string) (*entity.User, error) {
	return u.userRepo.Update(userID, name)
}
