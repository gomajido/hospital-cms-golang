package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gomajido/hospital-cms-golang/config"
	"github.com/gomajido/hospital-cms-golang/internal/dependency"
	appointmentRouter "github.com/gomajido/hospital-cms-golang/internal/module/appointment/router"
	articleRouter "github.com/gomajido/hospital-cms-golang/internal/module/article/router"
	authRouter "github.com/gomajido/hospital-cms-golang/internal/module/auth/router"
	doctorRouter "github.com/gomajido/hospital-cms-golang/internal/module/doctor/router"
)

type Router struct {
	ApplicationHandler *dependency.ApplicationHandler
	HttpConfig         config.HttpConfig
}

func Run(r *Router) error {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Apexa",
		AppName:       "Apexa Core API",
		BodyLimit:     100 * 1024 * 1024, // 10mb
	})
	app.Use(recover.New())

	// Health check route
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON("pong")
	})

	app.Get("/pong", func(c *fiber.Ctx) error {
		return c.JSON("ping")
	})

	// API routes
	v1 := app.Group("/api/v1")

	// Register auth routes
	authRouter.RegisterAuthRoutes(v1, r.ApplicationHandler.AuthHandler, r.ApplicationHandler.AuthMiddleware)

	// Register article routes
	articleRouter.RegisterArticleRoutes(v1, r.ApplicationHandler.ArticleHandler, r.ApplicationHandler.AuthMiddleware)

	// Register doctor routes
	doctorRouter.RegisterDoctorRoutes(v1, r.ApplicationHandler.DoctorHandler, r.ApplicationHandler.AuthMiddleware)

	// Register appointment routes
	appointmentRouter.RegisterAppointmentRoutes(v1, r.ApplicationHandler.AppointmentHandler, r.ApplicationHandler.AuthMiddleware)

	err := app.Listen(r.HttpConfig.Address)
	if err != nil {
		return err
	}
	return nil
}
