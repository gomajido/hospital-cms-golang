package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gomajido/hospital-cms-golang/internal/module/auth/domain"
	"github.com/gomajido/hospital-cms-golang/internal/response"
	"github.com/google/uuid"
)

type AuthMiddleware struct {
	usecase domain.AuthUsecase
}

func NewAuthMiddleware(usecase domain.AuthUsecase) *AuthMiddleware {
	return &AuthMiddleware{
		usecase: usecase,
	}
}

// Protected validates the token and allows access if valid
func (m *AuthMiddleware) Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from header
		header := c.Get("Authorization")
		if header == "" {
			err := errors.New("missing authorization header")
			return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(err))
		}

		// Check bearer scheme
		parts := strings.Split(header, " ")
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			err := errors.New("invalid authorization format")
			return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(err))
		}

		// Split ID and token
		credentials := strings.Split(parts[1], "|")
		if len(credentials) != 2 {
			err := errors.New("invalid token format")
			return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(err))
		}

		tokenIDStr := credentials[0]
		token := credentials[1]

		if tokenIDStr == "" || token == "" {
			err := errors.New("missing token ID or token")
			return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(err))
		}

		// Parse token ID to UUID
		tokenID, err := uuid.Parse(tokenIDStr)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(errors.New("invalid token ID format")))
		}

		// Validate token
		err = m.usecase.ValidateUserToken(c.Context(), tokenIDStr, token)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(err))
		}

		// Get token from database
		userToken, err := m.usecase.GetUserTokenByID(c.Context(), tokenID)
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(err))
		}

		// Set user token in context
		c.Locals("user_token", userToken)

		return c.Next()
	}
}

// HasAbility checks if the user has the required ability
func (m *AuthMiddleware) HasAbility(ability string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userToken, ok := c.Locals("user_token").(*domain.UserToken)
		if !ok {
			err := errors.New("missing user token")
			return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(err))
		}

		// Check if user has the required ability
		hasAbility := false
		for _, a := range userToken.Ability {
			if a == ability {
				hasAbility = true
				break
			}
		}

		if !hasAbility {
			err := errors.New("insufficient permissions")
			return c.Status(fiber.StatusForbidden).JSON(response.ErrForbidden.WithError(err))
		}

		return c.Next()
	}
}

// HasAnyAbility checks if the user has any of the required abilities
func (m *AuthMiddleware) HasAnyAbility(abilities ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userToken, ok := c.Locals("user_token").(*domain.UserToken)
		if !ok {
			err := errors.New("missing user token")
			return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(err))
		}

		// Check if user has any of the required abilities
		hasAbility := false
		for _, required := range abilities {
			for _, userAbility := range userToken.Ability {
				if required == userAbility {
					hasAbility = true
					break
				}
			}
			if hasAbility {
				break
			}
		}

		if !hasAbility {
			err := errors.New("insufficient permissions")
			return c.Status(fiber.StatusForbidden).JSON(response.ErrForbidden.WithError(err))
		}

		return c.Next()
	}
}

// HasAllAbilities checks if the user has all of the required abilities
func (m *AuthMiddleware) HasAllAbilities(abilities ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userToken, ok := c.Locals("user_token").(*domain.UserToken)
		if !ok {
			err := errors.New("missing user token")
			return c.Status(fiber.StatusUnprocessableEntity).JSON(response.ErrUnprocessableEntity.WithError(err))
		}

		// Check if user has all required abilities
		for _, required := range abilities {
			hasAbility := false
			for _, userAbility := range userToken.Ability {
				if required == userAbility {
					hasAbility = true
					break
				}
			}
			if !hasAbility {
				err := errors.New("insufficient permissions")
				return c.Status(fiber.StatusForbidden).JSON(response.ErrForbidden.WithError(err))
			}
		}

		return c.Next()
	}
}
