package handlers

import "github.com/gofiber/fiber/v2"

func HealthCheckHandler(c *fiber.Ctx) error {
	return c.SendString("OK")
}

func TestHandler(c *fiber.Ctx) error {
	return c.SendString("Hello World")
}
