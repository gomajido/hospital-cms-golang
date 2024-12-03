package router

import (
	"github.com/gofiber/fiber/v2"

	"github.com/gomajido/hospital-cms-golang/internal/middleware"
	"github.com/gomajido/hospital-cms-golang/internal/module/doctor/handler"
)

// RegisterDoctorRoutes registers all doctor routes
func RegisterDoctorRoutes(router fiber.Router, h *handler.DoctorHandler, authMiddleware *middleware.AuthMiddleware) {
	doctors := router.Group("/doctors")

	// Public routes
	doctors.Get("", h.List)
	doctors.Get("/:id", h.GetByID)

	// Protected routes for admins only
	doctors.Use(authMiddleware.Protected())
	doctors.Use(authMiddleware.HasAbility("admin"))
	doctors.Post("", h.Create)
	doctors.Put("/:id", h.Update)
	doctors.Delete("/:id", h.Delete)

	// Schedule routes
	doctors.Post("/:id/schedules", h.CreateSchedule)
	doctors.Get("/:id/schedules", h.GetSchedules)
	doctors.Put("/schedules/:id", h.UpdateSchedule)
	doctors.Delete("/schedules/:id", h.DeleteSchedule)

	// Reschedule routes
	doctors.Post("/schedules/:id/reschedules", h.CreateReschedule)
	doctors.Get("/schedules/:id/reschedules", h.GetReschedules)
	doctors.Put("/reschedules/:id", h.UpdateReschedule)
	doctors.Delete("/reschedules/:id", h.DeleteReschedule)
}
