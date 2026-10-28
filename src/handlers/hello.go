package handlers

import (
	"github.com/gofiber/fiber/v3"
)

// GetMessageHandler returns the welcome message for /api/v1/hello.
func GetMessageHandler(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Welcome to your Golang Fiber API!",
	})
}

// HealthCheckHandler returns service health status for /api/v1/health.
func HealthCheckHandler(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
	})
}
