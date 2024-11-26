package handler

import (
	"errors"
	"strings"
	"unicode"

	"github.com/gofiber/fiber/v2"
	"github.com/gomajido/hospital-cms-golang/internal/module/auth/domain"
	"github.com/gomajido/hospital-cms-golang/internal/response"
)

type authHandler struct {
	usecase domain.AuthUsecase
}

// NewAuthHandler creates a new auth handler instance
func NewAuthHandler(usecase domain.AuthUsecase) domain.AuthHandler {
	return &authHandler{
		usecase: usecase,
	}
}

// Register handles user registration
func (h *authHandler) Register(c *fiber.Ctx) error {
	var req domain.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(err))
	}

	// Validate required fields
	var validationErrors []string

	// Email validation
	if req.Email == "" {
		validationErrors = append(validationErrors, "email is required")
	} else if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		validationErrors = append(validationErrors, "invalid email format")
	}

	// Password validation
	if req.Password == "" {
		validationErrors = append(validationErrors, "password is required")
	} else if len(req.Password) < 8 {
		validationErrors = append(validationErrors, "password must be at least 8 characters long")
	} else {
		var (
			hasUpper   bool
			hasLower   bool
			hasNumber  bool
			hasSpecial bool
		)
		for _, char := range req.Password {
			switch {
			case unicode.IsUpper(char):
				hasUpper = true
			case unicode.IsLower(char):
				hasLower = true
			case unicode.IsNumber(char):
				hasNumber = true
			case unicode.IsPunct(char) || unicode.IsSymbol(char):
				hasSpecial = true
			}
		}
		if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
			validationErrors = append(validationErrors, "password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
		}
	}

	// Name validation
	if req.Name == "" {
		validationErrors = append(validationErrors, "name is required")
	} else if len(req.Name) < 2 {
		validationErrors = append(validationErrors, "name must be at least 2 characters long")
	}

	// Phone validation (if provided)
	if req.Phone != "" {
		// Remove any non-digit characters
		phone := strings.Map(func(r rune) rune {
			if unicode.IsDigit(r) {
				return r
			}
			return -1
		}, req.Phone)

		if len(phone) < 10 || len(phone) > 15 {
			validationErrors = append(validationErrors, "phone number must be between 10 and 15 digits")
		}
		// Update the request with cleaned phone number
		req.Phone = phone
	}

	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(errors.New(strings.Join(validationErrors, "; "))))
	}

	resp, err := h.usecase.Register(c.Context(), &req)
	if err != nil {
		if err.Error() == "user already exists" {
			return c.Status(fiber.StatusConflict).JSON(response.ErrBadRequest.WithError(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(errors.New("something went wrong")))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Ok.WithData(resp))
}

// Login handles user authentication
func (h *authHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(err))
	}

	// Validate request
	if req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(errors.New("missing required fields")))
	}

	resp, err := h.usecase.Login(c.Context(), &req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrUnauthorized.WithError(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(errors.New("something went wrong")))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Ok.WithData(resp))
}

// Logout handles user logout
func (h *authHandler) Logout(c *fiber.Ctx) error {
	userToken, ok := c.Locals("user_token").(*domain.UserToken)
	if !ok {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(errors.New("missing user token")))
	}

	if err := h.usecase.Logout(c.Context(), userToken.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(errors.New("something went wrong")))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(nil))
}

// GetUserByID retrieves a user by their ID
func (h *authHandler) GetUserByID(c *fiber.Ctx) error {
	userToken, ok := c.Locals("user_token").(*domain.UserToken)
	if !ok {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(errors.New("missing user token")))
	}

	user, err := h.usecase.GetUserByID(c.Context(), userToken.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(errors.New("something went wrong")))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(user))
}

// UpdateUser updates user information
func (h *authHandler) UpdateUser(c *fiber.Ctx) error {
	userToken, ok := c.Locals("user_token").(*domain.UserToken)
	if !ok {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(errors.New("missing user token")))
	}

	// First get the existing user
	existingUser, err := h.usecase.GetUserByID(c.Context(), userToken.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(errors.New("failed to get user")))
	}

	// Parse the update request
	var updateReq struct {
		Name   string `json:"name"`
		Phone  string `json:"phone"`
		Status string `json:"status"`
	}
	if err := c.BodyParser(&updateReq); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(err))
	}

	// Update only the allowed fields
	existingUser.Name = updateReq.Name
	existingUser.Phone = updateReq.Phone
	existingUser.Status = updateReq.Status

	if err := h.usecase.UpdateUser(c.Context(), existingUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(errors.New("failed to update user")))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(existingUser))
}

// AssignRoles assigns roles to a user
func (h *authHandler) AssignRoles(c *fiber.Ctx) error {
	var req struct {
		RoleNames []string `json:"role_names"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(err))
	}

	// Get user ID from params
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(errors.New("missing user ID")))
	}

	if err := h.usecase.AssignRoles(c.Context(), userID, req.RoleNames); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(errors.New("something went wrong")))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(nil))
}

// GetUserRoles retrieves all roles assigned to a user
func (h *authHandler) GetUserRoles(c *fiber.Ctx) error {
	// Get user ID from params
	userID := c.Params("id")
	if userID == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(errors.New("missing user ID")))
	}

	roles, err := h.usecase.GetUserRoles(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(errors.New("something went wrong")))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(roles))
}
