package domain

import "github.com/gofiber/fiber/v2"

// AuthHandler defines the interface for auth-related HTTP handlers
type AuthHandler interface {
	// Registration and authentication
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error

	// User management
	GetUserByID(c *fiber.Ctx) error
	UpdateUser(c *fiber.Ctx) error

	// Role management
	AssignRoles(c *fiber.Ctx) error
	GetUserRoles(c *fiber.Ctx) error
}
