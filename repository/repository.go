package repository

import "goauth/entity"

type UserRepository interface {
	Insert(email, name, passwordHash string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByID(id string) (*entity.User, error)
	Update(id, name string) (*entity.User, error)
}
