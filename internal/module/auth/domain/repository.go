package domain

import (
	"context"
)

// AuthRepository defines the interface for auth-related database operations
type AuthRepository interface {
	// User operations
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id string) error

	// Role operations
	GetRolesByNames(ctx context.Context, names []string) ([]Role, error)
	GetRoleByID(ctx context.Context, id string) (*Role, error)
	AssignRolesToUser(ctx context.Context, userID string, roleIDs []string) error
	GetUserRoles(ctx context.Context, userID string) ([]Role, error)

	// Token management
	CreateUserToken(ctx context.Context, token *UserToken) error
	GetUserTokenByIDAndToken(ctx context.Context, userID string, token string) (*UserToken, error)
	InvalidateUserToken(ctx context.Context, userID string, token string) error
	InvalidateUserTokens(ctx context.Context, userID string) error
}
