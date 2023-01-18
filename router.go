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
	libGroup.Get("/:id", handlers.GetLibrary)
	libGroup.Post("/", handlers.CreateLibraryHandler)
	libGroup.Put("/:id", handlers.UpdateLibrary)
	libGroup.Delete("/:id", handlers.DeleteLibrary)
	libGroup.Get("/books/:id", handlers.GetLibraryBooks)

	bookGroup := v1.Group("/book") // /api/v1/library
	bookGroup.Get("/", handlers.GetBooks)
	bookGroup.Get("/:id", handlers.GetBook)
	bookGroup.Post("/", handlers.CreateBook)
	bookGroup.Put("/:id", handlers.UpdateBook)
	bookGroup.Delete("/:id", handlers.DeleteBook)

	return app
}
