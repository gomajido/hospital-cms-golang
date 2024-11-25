package domain

import (
	"context"
)

// AuthUsecase defines the interface for auth-related business logic
type AuthUsecase interface {
	// Registration and authentication
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	Logout(ctx context.Context, userID string, token string) error
	GetUserTokenByIDAndToken(ctx context.Context, userID string, token string) (*UserToken, error)
	GetUserByID(ctx context.Context, userID string) (*User, error)

	// User management
	UpdateUser(ctx context.Context, user *User) error

	// Role management
	AssignRoles(ctx context.Context, userID string, roleNames []string) error
	GetUserRoles(ctx context.Context, userID string) ([]Role, error)
}
