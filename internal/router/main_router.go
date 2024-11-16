package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gomajido/hospital-cms-golang/config"
	"github.com/gomajido/hospital-cms-golang/internal/dependency"
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
		ServerHeader:  "Nabitu",
		AppName:       "Nabitu Core API",
		BodyLimit:     100 * 1024 * 1024, // 10mb
	})
	app.Use(recover.New())
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON("pong")
	})

	err := app.Listen(r.HttpConfig.Address)

	if err != nil {
		return err
	}
	return nil
}
