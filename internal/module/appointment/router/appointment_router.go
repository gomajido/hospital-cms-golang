package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gomajido/hospital-cms-golang/internal/middleware"
	"github.com/gomajido/hospital-cms-golang/internal/module/appointment/handler"
)

func RegisterAppointmentRoutes(router fiber.Router, appointmentHandler *handler.AppointmentHandler, authMiddleware *middleware.AuthMiddleware) {
	appointmentRouter := router.Group("/appointments")

	// Protected routes (require authentication)
	appointmentRouter.Use(authMiddleware.Protected())
	{
		// Create appointment
		appointmentRouter.Post("/", appointmentHandler.Create)

		// Get appointments for current user
		appointmentRouter.Get("/me", appointmentHandler.GetByUserID)

		// Get appointments for a doctor
		appointmentRouter.Get("/doctor/:doctor_id", appointmentHandler.GetByDoctorID)

		// Get specific appointment
		appointmentRouter.Get("/:id", appointmentHandler.GetByID)

		// Cancel appointment
		appointmentRouter.Post("/:id/cancel", appointmentHandler.Cancel)

		// Reschedule appointment
		appointmentRouter.Post("/:id/reschedule", appointmentHandler.Reschedule)

		// Check availability
		appointmentRouter.Post("/check-availability", appointmentHandler.CheckAvailability)
	}
}
