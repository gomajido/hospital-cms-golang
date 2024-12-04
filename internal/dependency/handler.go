package dependency

import (
	"context"

	"github.com/gomajido/hospital-cms-golang/config"
	"github.com/gomajido/hospital-cms-golang/internal/middleware"
	appointmentHandler "github.com/gomajido/hospital-cms-golang/internal/module/appointment/handler"
	articleHandler "github.com/gomajido/hospital-cms-golang/internal/module/article/handler"
	"github.com/gomajido/hospital-cms-golang/internal/module/auth/domain"
	"github.com/gomajido/hospital-cms-golang/internal/module/auth/handler"
	doctorHandler "github.com/gomajido/hospital-cms-golang/internal/module/doctor/handler"
)

type ApplicationHandler struct {
	AuthHandler        domain.AuthHandler
	AuthMiddleware     *middleware.AuthMiddleware
	ArticleHandler     *articleHandler.ArticleHandler
	DoctorHandler      *doctorHandler.DoctorHandler
	AppointmentHandler *appointmentHandler.AppointmentHandler
}

func InitHandlers(ctx context.Context, cfg *config.Config, service *AppUsecase) *ApplicationHandler {
	authMiddleware := middleware.NewAuthMiddleware(service.AuthUsecase)

	return &ApplicationHandler{
		AuthHandler:        handler.NewAuthHandler(service.AuthUsecase),
		AuthMiddleware:     authMiddleware,
		ArticleHandler:     articleHandler.NewArticleHandler(service.ArticleUsecase),
		DoctorHandler:      doctorHandler.NewDoctorHandler(service.DoctorUsecase),
		AppointmentHandler: appointmentHandler.NewAppointmentHandler(service.AppointmentUsecase),
	}
}
