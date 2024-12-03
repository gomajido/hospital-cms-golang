package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gomajido/hospital-cms-golang/internal/module/auth/domain"
	"github.com/gomajido/hospital-cms-golang/internal/response"
	"github.com/google/uuid"
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
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest)
	}

	if errors := req.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithErrorInfo(errors))
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
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest)
	}

	if errors := req.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithErrorInfo(errors))
	}

	resp, err := h.usecase.Login(c.Context(), &req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.ErrUnauthorized.WithError(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(resp))
}

// Logout handles user logout
func (h *authHandler) Logout(c *fiber.Ctx) error {
	if err := h.usecase.Logout(c.Context(), c.Locals("user_token").(string)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok)
}

// GetUserByID retrieves a user by their ID
func (h *authHandler) GetUserByID(c *fiber.Ctx) error {
	userIDStr := c.Params("id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(fmt.Errorf("user id is required")))
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(fmt.Errorf("invalid user id format")))
	}

	user, err := h.usecase.GetUserByID(c.Context(), userID)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrRecordNotFound.WithError(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(user))
}

// UpdateUser handles updating user details
func (h *authHandler) UpdateUser(c *fiber.Ctx) error {
	// Get user token from context
	userToken, ok := c.Locals("user_token").(*domain.UserToken)
	if !ok {
		err := fmt.Errorf("missing user token")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(err))
	}

	var req domain.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest.WithError(err))
	}

	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithErrorInfo(err))
	}

	// Get current user data
	user, err := h.usecase.GetUserByID(c.Context(), userToken.UserID)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrRecordNotFound.WithError(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	// Update user fields
	user.Name = req.Name
	user.Phone = req.Phone

	// Update user in database
	err = h.usecase.UpdateUser(c.Context(), user)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrRecordNotFound.WithError(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok)
}

// AssignRoles assigns roles to a user
func (h *authHandler) AssignRoles(c *fiber.Ctx) error {
	userIDStr := c.Params("id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(fmt.Errorf("user id is required")))
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(fmt.Errorf("invalid user id format")))
	}

	var req domain.AssignRolesRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrBadRequest)
	}

	if errors := req.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithErrorInfo(errors))
	}

	if err := h.usecase.AssignRoles(c.Context(), userID, req.RoleNames); err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrRecordNotFound.WithError(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok)
}

// GetUserRoles retrieves all roles assigned to a user
func (h *authHandler) GetUserRoles(c *fiber.Ctx) error {
	userIDStr := c.Params("id")
	if userIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(fmt.Errorf("user id is required")))
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrInvalidParam.WithError(fmt.Errorf("invalid user id format")))
	}

	roles, err := h.usecase.GetUserRoles(c.Context(), userID)
	if err != nil {
		if err.Error() == "user not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrRecordNotFound.WithError(err))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrInternalServer.WithError(err))
	}

	return c.Status(fiber.StatusOK).JSON(response.Ok.WithData(roles))
}
