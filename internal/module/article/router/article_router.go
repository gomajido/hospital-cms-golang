package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/gomajido/hospital-cms-golang/internal/middleware"
	"github.com/gomajido/hospital-cms-golang/internal/module/article/handler"
)

// RegisterArticleRoutes registers all article routes
func RegisterArticleRoutes(router fiber.Router, h *handler.ArticleHandler, authMiddleware *middleware.AuthMiddleware) {
	articles := router.Group("/articles")

	// Public routes
	articles.Get("", h.List)
	articles.Get("/:id", h.GetByID)
	articles.Get("/slug/:slug", h.GetBySlug)

	// Protected routes for admins only
	articles.Use(authMiddleware.Protected())
	articles.Use(authMiddleware.HasAbility("admin"))
	articles.Post("", h.Create)
	articles.Put("/:id", h.Update)
	articles.Delete("/:id", h.Delete)
}
