package repository

import (
	"database/sql"
	"fmt"
	"goauth/entity"
	"goauth/errors"
	"log"
	"time"

	"github.com/google/uuid"
)

type SQLiteUserRepo struct {
	db *sql.DB
}

func NewSQLiteUserRepo(db *sql.DB) UserRepository {
	return &SQLiteUserRepo{db: db}
}
func (r *SQLiteUserRepo) Insert(email, name, passwordHash string) (*entity.User, error) {
	userId := uuid.New().String()
	now := time.Now()

	_, err := r.db.Exec(
		"INSERT INTO users (id, email, username, password_hash, role, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		userId, email, name, passwordHash, "user", now,
	)
	if err != nil {
		log.Printf("DEBUG: Insert failed - %v", err)
		return nil, fmt.Errorf("Error when insert user: %w", err)
	}

	returnedUser := &entity.User{}
	if err := r.db.QueryRow(
		"SELECT id, email, username, password_hash, role, created_at FROM users WHERE id = ?",
		userId,
	).Scan(&returnedUser.ID, &returnedUser.Email, &returnedUser.Name, &returnedUser.PasswordHash, &returnedUser.Role, &returnedUser.CreatedAt); err != nil {
		log.Printf("DEBUG: Select failed - %v", err)
		return nil, fmt.Errorf("Error selecting user: %w", err)
	}

	return returnedUser, nil
}

func (r *SQLiteUserRepo) FindByEmail(email string) (*entity.User, error) {
	user := &entity.User{}

	if err := r.db.QueryRow(
		"SELECT id, email, username, password_hash, role, created_at FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Role, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFoundError{Message: "user not found"}

		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	return user, nil
}

func (r *SQLiteUserRepo) FindByID(id string) (*entity.User, error) {
	user := &entity.User{}

	if err := r.db.QueryRow(
		"SELECT id, email, username, password_hash, role, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.Role, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.NotFoundError{Message: "user not found"}
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	return user, nil
}

func (r *SQLiteUserRepo) Update(id, name string) (*entity.User, error) {
	_, err := r.db.Exec(
		"UPDATE users SET username = ? WHERE id = ?",
		name, id,
	)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return r.FindByID(id)
}
