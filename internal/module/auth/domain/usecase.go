package domain

import (
	"context"
)

// AuthUsecase defines the interface for auth-related business logic
type AuthUsecase interface {
	// User operations
	Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	Logout(ctx context.Context, tokenID string) error
	UpdateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, userID string) (*User, error)

	// Role operations
	AssignRoles(ctx context.Context, userID string, roleNames []string) error
	GetUserRoles(ctx context.Context, userID string) ([]Role, error)

	// Token operations
	GetUserTokenByID(ctx context.Context, tokenID string) (*UserToken, error)
	ValidateUserToken(ctx context.Context, tokenID string, token string) error
	InvalidateUserToken(ctx context.Context, tokenID string) error
	InvalidateUserTokens(ctx context.Context, userID string) error
}
