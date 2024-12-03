package domain

import (
	"time"

	"github.com/google/uuid"
)

// RegisterResponse represents the registration response
type RegisterResponse struct {
	User *User `json:"user"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Status    string    `json:"status"`
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}
