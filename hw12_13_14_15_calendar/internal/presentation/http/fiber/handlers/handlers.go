package handlers

import (
	"github.com/gofiber/fiber/v3"
)

// HelloWorld is a handler for hello world.
func HelloWorld(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("Hello World!")
}

// HealthCheck is a handler for health checking.
func HealthCheck(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("ok")
}
