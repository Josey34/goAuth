package entity

import "time"

type User struct {
	ID           string
	Email        string
	Name         string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
}
