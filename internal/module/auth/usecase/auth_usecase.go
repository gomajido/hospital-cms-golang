package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gomajido/hospital-cms-golang/config"
	"github.com/gomajido/hospital-cms-golang/internal/module/auth/domain"
	"github.com/gomajido/hospital-cms-golang/pkg/app_log"
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
	// Check if user already exists
	existingUser, err := a.repo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		app_log.Errorf("failed to hash password: %v", err)
		return nil, errors.New("failed to process password")
	}

	// Create user
	user := &domain.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Phone:    req.Phone,
	}

	err = a.repo.CreateUser(ctx, user)
	if err != nil {
		app_log.Errorf("failed to create user: %v", err)
		return nil, err
	}

	return &domain.RegisterResponse{
		User: user,
	}, nil
}

// Login handles user authentication
func (a *authUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	// Get user by email
	user, err := a.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		app_log.Errorf("failed to get user: %v", err)
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token, err := a.generateToken()
	if err != nil {
		app_log.Errorf("failed to generate token: %v", err)
		return nil, err
	}

	// Get user roles and extract abilities
	roles, err := a.repo.GetUserRoles(ctx, user.ID)
	if err != nil {
		app_log.Errorf("failed to get user roles: %v", err)
		return nil, err
	}

	var abilities []string
	for _, role := range roles {
		// Add role abilities to the token
		// This is where you would typically get abilities from your role definitions
		abilities = append(abilities, role.Name) // Using role names as abilities for now
	}

	// Create user token
	expiry := time.Now().Add(24 * time.Hour)
	userToken := &domain.UserToken{
		UserID:    user.ID,
		Token:     token,
		Ability:   abilities,
		ExpiredAt: expiry,
	}

	err = a.repo.CreateUserToken(ctx, userToken)
	if err != nil {
		app_log.Errorf("failed to create user token: %v", err)
		return nil, err
	}

	return &domain.LoginResponse{
		User:         user,
		AccessToken:  token,
		RefreshToken: token,
		TokenType:    "Bearer",
		ExpiresIn:    int64(expiry.Unix()),
		Ability:      userToken.Ability,
	}, nil
}

// Logout invalidates the user's token
func (a *authUsecase) Logout(ctx context.Context, userID string, token string) error {

	return a.repo.InvalidateUserToken(ctx, userID, token)
}

// UpdateUser updates user information
func (a *authUsecase) UpdateUser(ctx context.Context, user *domain.User) error {

	return a.repo.UpdateUser(ctx, user)
}

// AssignRoles assigns roles to a user
func (a *authUsecase) AssignRoles(ctx context.Context, userID string, roleNames []string) error {

	return a.repo.AssignRolesToUser(ctx, userID, roleNames)
}

// GetUserRoles retrieves all roles assigned to a user
func (a *authUsecase) GetUserRoles(ctx context.Context, userID string) ([]domain.Role, error) {

	return a.repo.GetUserRoles(ctx, userID)
}

// GetUserTokenByIDAndToken retrieves a user token by user ID and token
func (a *authUsecase) GetUserTokenByIDAndToken(ctx context.Context, userID string, token string) (*domain.UserToken, error) {
	return a.repo.GetUserTokenByIDAndToken(ctx, userID, token)
}

// GetUserByID retrieves a user by their ID
func (a *authUsecase) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	return a.repo.GetUserByID(ctx, userID)
}
