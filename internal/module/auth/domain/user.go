package domain

import (
	"time"
)

type User struct {
	ID              string     `json:"id"`
	Email           string     `json:"email"`
	Password        string     `json:"-"` // Never expose password in JSON
	Name            string     `json:"name"`
	Phone           string     `json:"phone"`
	Status          string     `json:"status"`
	EmailVerifiedAt *time.Time `json:"email_verified_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
	Roles           []Role     `json:"roles,omitempty"`
}

type Role struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
