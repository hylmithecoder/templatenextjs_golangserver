package middleware

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v3"
)

// TestVerboseLoggerPassthrough verifies the middleware forwards the request to
// the handler and leaves the response status and body intact.
func TestVerboseLoggerPassthrough(t *testing.T) {
	app := fiber.New()
	app.Use("/api/v1", VerboseLogger())
	app.Post("/api/v1/echo", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status": "success",
			"echo":   string(c.Body()),
		})
	})

	payload := `{"username":"admin","password":"s3cr3t"}`
	req := httptest.NewRequest("POST", "/api/v1/echo", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusCreated {
		t.Errorf("expected status 201, got %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	var body map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["status"] != "success" {
		t.Errorf("expected success status, got %v", body["status"])
	}
	if echo, _ := body["echo"].(string); echo != payload {
		t.Errorf("request body corrupted by middleware: got %q want %q", echo, payload)
	}
}

// TestVerboseLoggerPropagatesError ensures handler errors still surface through
// the middleware.
func TestVerboseLoggerPropagatesError(t *testing.T) {
	app := fiber.New()
	app.Use("/api/v1", VerboseLogger())
	app.Get("/api/v1/boom", func(c fiber.Ctx) error {
		return fiber.NewError(fiber.StatusInternalServerError, "boom")
	})

	req := httptest.NewRequest("GET", "/api/v1/boom", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", resp.StatusCode)
	}
}
