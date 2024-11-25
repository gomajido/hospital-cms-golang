package domain

import (
	"context"
)

// AuthUsecase defines the interface for auth-related business logic
type AuthUsecase interface {
	// Registration and authentication
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	
	// User management
	GetUserByID(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	
	// Role management
	AssignRoles(ctx context.Context, userID string, roleNames []string) error
	GetUserRoles(ctx context.Context, userID string) ([]Role, error)
}
