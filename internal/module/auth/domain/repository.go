package domain

import (
	"context"

	"github.com/google/uuid"
)

// AuthRepository defines the interface for auth-related database operations
type AuthRepository interface {
	// User operations
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error

	// Role operations
	GetRolesByNames(ctx context.Context, names []string) ([]Role, error)
	GetRoleByID(ctx context.Context, id uuid.UUID) (*Role, error)
	AssignRolesToUser(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]Role, error)

	// Token management
	CreateUserToken(ctx context.Context, token *UserToken) (*UserToken, error)
	GetUserTokenByID(ctx context.Context, tokenID uuid.UUID) (*UserToken, error)
	InvalidateUserToken(ctx context.Context, tokenID uuid.UUID) error
	InvalidateUserTokens(ctx context.Context, userID uuid.UUID) error
}
