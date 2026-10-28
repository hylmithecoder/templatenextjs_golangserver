package handlers

import (
	"templatenextjs_golangserver/src/utils"

	"github.com/gofiber/fiber/v3"
)

// ErrorResponse builds a consistent error JSON body and logs it.
func ErrorResponse(c fiber.Ctx, status int, message string, err error) error {
	body := fiber.Map{
		"status":  "error",
		"message": message,
	}

	if err != nil {
		body["error"] = err.Error()
		utils.LogErr("%s: %v", message, err)
	}

	return c.Status(status).JSON(body)
}
