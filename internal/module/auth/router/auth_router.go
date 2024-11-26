package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gomajido/hospital-cms-golang/internal/middleware"
	"github.com/gomajido/hospital-cms-golang/internal/module/auth/domain"
)

// RegisterAuthRoutes registers all auth routes
func RegisterAuthRoutes(router fiber.Router, handler domain.AuthHandler, authMiddleware *middleware.AuthMiddleware) {
	auth := router.Group("/auth")

	// Public routes
	auth.Post("/register", handler.Register)
	auth.Post("/login", handler.Login)

	// Protected routes
	auth.Use(authMiddleware.Protected())
	auth.Post("/logout", handler.Logout)

	// User management routes
	users := auth.Group("/users")
	users.Get("/me", handler.GetUserByID)
	users.Put("/me", handler.UpdateUser)

	// Role management routes
	// users.Post("/:id/roles", handler.AssignRoles)
	// users.Get("/:id/roles", handler.GetUserRoles)
}
