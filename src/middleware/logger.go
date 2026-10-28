package middleware

import (
	"time"

	"templatenextjs_golangserver/src/utils"

	"github.com/gofiber/fiber/v3"
)

// VerboseLogger returns a Fiber middleware that logs every request together
// with its response: method, path, status code, latency, and the (redacted,
// truncated) request and response bodies. Credential fields such as password
// and token are masked by utils.SanitizeBody before printing.
//
// Mount it with a prefix (e.g. app.Use("/api/v1", middleware.VerboseLogger()))
// so only API traffic is logged, not static frontend assets.
func VerboseLogger() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()

		// Copy the request body before the handler consumes it. c.Body()
		// returns a buffer that may be reused by fasthttp, so snapshot it.
		reqBody := append([]byte(nil), c.Body()...)

		// Run the actual handler chain.
		err := c.Next()

		latency := time.Since(start)
		status := c.Response().StatusCode()
		respBody := c.Response().Body()

		utils.LogHTTP(c.Method(), c.OriginalURL(), status, latency, reqBody, respBody)

		return err
	}
}
