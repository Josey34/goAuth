package mock_test

import (
	"goauth/entity"
	"goauth/errors"
)

type MockUserRepository struct {
	Users map[string]*entity.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users: make(map[string]*entity.User),
	}
}

func (m *MockUserRepository) Insert(email, name, passwordHash string) (*entity.User, error) {
	user := &entity.User{
		ID:           email,
		Email:        email,
		Name:         name,
		PasswordHash: passwordHash,
		Role:         "user",
	}
	m.Users[email] = user
	return user, nil
}

func (m *MockUserRepository) FindByEmail(email string) (*entity.User, error) {
	user, ok := m.Users[email]
	if !ok {
		return nil, errors.NotFoundError{Message: "user not found"}
	}
	return user, nil
}

func (m *MockUserRepository) FindByID(id string) (*entity.User, error) {
	user, ok := m.Users[id]
	if !ok {
		return nil, errors.NotFoundError{Message: "user not found"}
	}
	return user, nil
}

func (m *MockUserRepository) Update(id, name string) (*entity.User, error) {
	user, ok := m.Users[id]
	if !ok {
		return nil, errors.NotFoundError{Message: "user not found"}
	}
	user.Name = name
	return user, nil
}
