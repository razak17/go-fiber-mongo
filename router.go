package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/razak17/go-fiber-mongo/handlers"
)

func generateApp() *fiber.App {
	app := fiber.New()

	api := app.Group("/api") // /api

	v1 := api.Group("/v1") // /api/v1

	// health check route
	v1.Get("/health", handlers.HealthCheckHandler)

	// library group and routes
	libGroup := v1.Group("/library") // /api/v1/library
	libGroup.Get("/", handlers.GetLibraries)
	libGroup.Post("/", handlers.CreateLibraryHandler)

	return app
}
