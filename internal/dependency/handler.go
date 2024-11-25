package dependency

import (
	"context"

	"github.com/gomajido/hospital-cms-golang/config"
	"github.com/gomajido/hospital-cms-golang/internal/middleware"
	"github.com/gomajido/hospital-cms-golang/internal/module/auth/domain"
	"github.com/gomajido/hospital-cms-golang/internal/module/auth/handler"
)

type ApplicationHandler struct {
	AuthHandler    domain.AuthHandler
	AuthMiddleware *middleware.AuthMiddleware
}

func InitHandlers(ctx context.Context, cfg *config.Config, service *AppUsecase) *ApplicationHandler {
	authMiddleware := middleware.NewAuthMiddleware(service.AuthUsecase)

	return &ApplicationHandler{
		AuthHandler:    handler.NewAuthHandler(service.AuthUsecase),
		AuthMiddleware: authMiddleware,
	}
}
