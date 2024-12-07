package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gomajido/hospital-cms-golang/config"
	"github.com/gomajido/hospital-cms-golang/internal/module/auth/domain"
	"github.com/gomajido/hospital-cms-golang/pkg/app_log"
	"github.com/google/uuid"
)

type authUsecase struct {
	repo domain.AuthRepository
	cfg  *config.Config
}

// NewAuthUsecase creates a new auth usecase instance
func NewAuthUsecase(repo domain.AuthRepository, cfg *config.Config) domain.AuthUsecase {
	return &authUsecase{
		repo: repo,
		cfg:  cfg,
	}
}

// generateToken generates a random token string
func (a *authUsecase) generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Register handles user registration
func (a *authUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.RegisterResponse, error) {
	// Check if user with the same email already exists
	existingUser, err := a.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		app_log.Errorf("Failed to check existing user: %v", err)
		return nil, err
	}

	if existingUser != nil {
		app_log.Errorf("User already exists with email: %s", req.Email)
		return nil, errors.New("user with this email already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		app_log.Errorf("Failed to hash password: %v", err)
		return nil, err
	}

	// Create the user
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Phone:    req.Phone,
	}

	err = a.repo.CreateUser(ctx, user)
	if err != nil {
		app_log.Errorf("Failed to create user: %v", err)
		return nil, err
	}

	// Get the default role (member)
	roles, err := a.repo.GetRolesByNames(ctx, []string{"member"})
	if err != nil {
		app_log.Errorf("Failed to get member role: %v", err)
		return nil, err
	}

	if len(roles) == 0 {
		app_log.Error("Default role 'member' not found")
		return nil, errors.New("default role 'member' not found")
	}

	// Convert role to UUID slice
	roleIDs := []uuid.UUID{roles[0].ID}

	// Assign the member role to the user
	err = a.repo.AssignRolesToUser(ctx, user.ID, roleIDs)
	if err != nil {
		app_log.Errorf("Failed to assign member role to user: %v", err)
		return nil, err
	}

	// Get the updated user with roles
	updatedUser, err := a.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		app_log.Errorf("Failed to get updated user: %v", err)
		return nil, err
	}

	return &domain.RegisterResponse{
		User: updatedUser,
	}, nil
}

// Login handles user authentication
func (a *authUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	// Get user by email
	user, err := a.repo.GetUserByEmail(ctx, req.Email)

	if err != nil || user == nil {
		app_log.Errorf("User not found: %v", err)
		return nil, errors.New("invalid credentials")
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		app_log.Errorf("User not found: %v", err)
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token, err := a.generateToken()
	if err != nil {
		app_log.Errorf("Failed to generate token: %v", err)
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Get user roles
	roles, err := a.repo.GetUserRoles(ctx, user.ID)
	if err != nil {
		app_log.Errorf("Failed to get user roles: %v", err)
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Convert roles to abilities
	var abilities []string
	for _, role := range roles {
		abilities = append(abilities, role.Name)
	}

	// Create user token
	expiry := time.Now().Add(24 * time.Hour)
	userToken := &domain.UserToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		Ability:   abilities,
		ExpiredAt: expiry,
	}

	// Save token
	createdToken, err := a.repo.CreateUserToken(ctx, userToken)
	if err != nil {
		app_log.Errorf("Failed to create user token: %v", err)
		return nil, fmt.Errorf("failed to create user token: %w", err)
	}

	return &domain.LoginResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Phone:     user.Phone,
		Status:    user.Status,
		Token:     createdToken.ID.String() + "|" + token,
		ExpiredAt: createdToken.ExpiredAt,
	}, nil
}

// Logout invalidates the user's token
func (a *authUsecase) Logout(ctx context.Context, tokenIDStr string) error {
	tokenID, err := uuid.Parse(tokenIDStr)
	if err != nil {
		return fmt.Errorf("invalid token ID format: %w", err)
	}
	return a.InvalidateUserToken(ctx, tokenID)
}

// UpdateUser updates user information
func (a *authUsecase) UpdateUser(ctx context.Context, user *domain.User) error {
	return a.repo.UpdateUser(ctx, user)
}

// AssignRoles assigns roles to a user
func (a *authUsecase) AssignRoles(ctx context.Context, userID uuid.UUID, roleNames []string) error {
	// Get roles by names
	roles, err := a.repo.GetRolesByNames(ctx, roleNames)
	if err != nil {
		return err
	}

	// Convert role names to UUIDs
	roleIDs := make([]uuid.UUID, len(roles))
	for i, role := range roles {
		roleIDs[i] = role.ID
	}

	return a.repo.AssignRolesToUser(ctx, userID, roleIDs)
}

// GetUserRoles retrieves all roles assigned to a user
func (a *authUsecase) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]domain.Role, error) {
	return a.repo.GetUserRoles(ctx, userID)
}

// GetUserTokenByID retrieves a user token by ID
func (a *authUsecase) GetUserTokenByID(ctx context.Context, tokenID uuid.UUID) (*domain.UserToken, error) {
	userToken, err := a.repo.GetUserTokenByID(ctx, tokenID)
	if err != nil {
		app_log.Errorf("Failed to get user token: %v", err)
		return nil, fmt.Errorf("failed to get user token: %w", err)
	}

	if userToken == nil {
		return nil, errors.New("token not found")
	}

	return userToken, nil
}

// ValidateUserToken validates a user token
func (a *authUsecase) ValidateUserToken(ctx context.Context, tokenIDStr string, token string) error {
	tokenID, err := uuid.Parse(tokenIDStr)
	if err != nil {
		return fmt.Errorf("invalid token ID format: %w", err)
	}

	userToken, err := a.repo.GetUserTokenByID(ctx, tokenID)
	if err != nil {
		return fmt.Errorf("failed to get user token: %w", err)
	}

	if userToken == nil {
		return errors.New("token not found")
	}

	if userToken.Token != token {
		return errors.New("invalid token")
	}

	if time.Now().After(userToken.ExpiredAt) {
		return errors.New("token expired")
	}

	return nil
}

// InvalidateUserToken invalidates a user token
func (a *authUsecase) InvalidateUserToken(ctx context.Context, tokenID uuid.UUID) error {
	return a.repo.InvalidateUserToken(ctx, tokenID)
}

// InvalidateUserTokens invalidates all tokens for a user
func (a *authUsecase) InvalidateUserTokens(ctx context.Context, userID uuid.UUID) error {
	return a.repo.InvalidateUserTokens(ctx, userID)
}

// GetUserByID retrieves a user by their ID
func (a *authUsecase) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	return a.repo.GetUserByID(ctx, userID)
}
