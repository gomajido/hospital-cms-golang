package handler

import (
	"errors"
	"strings"

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

	// Validate request
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(errors.New("missing required fields")))
	}

	resp, err := h.usecase.Register(c.Context(), &req)
	if err != nil {
		if err.Error() == "user already exists" {
			return c.Status(fiber.StatusConflict).JSON(response.ErrBadRequest.WithError(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
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
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Ok.WithData(resp))
}

// Logout handles user logout
func (h *authHandler) Logout(c *fiber.Ctx) error {
	// Get token from header
	header := c.Get("Authorization")
	if header == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(errors.New("missing authorization header")))
	}

	// Check bearer scheme
	parts := strings.Split(header, " ")
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(errors.New("invalid authorization format")))
	}

	// Split ID and token
	credentials := strings.Split(parts[1], "|")
	if len(credentials) != 2 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(errors.New("invalid token format")))
	}

	userID := credentials[0]
	token := credentials[1]

	if err := h.usecase.Logout(c.Context(), userID, token); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(nil))
}

// GetUserByID retrieves a user by their ID
func (h *authHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(errors.New("missing user ID")))
	}

	user, err := h.usecase.GetUserByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	if user == nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ErrRecordNotFound.WithError(errors.New("user not found")))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(user))
}

// UpdateUser updates user information
func (h *authHandler) UpdateUser(c *fiber.Ctx) error {
	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(err))
	}

	// Get user ID from params
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(errors.New("missing user ID")))
	}
	user.ID = id

	if err := h.usecase.UpdateUser(c.Context(), &user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(user))
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
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
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
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(roles))
}
