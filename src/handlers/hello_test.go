package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestGetMessageHandler(t *testing.T) {
	app := fiber.New()
	app.Get("/api/v1/hello", GetMessageHandler)

	req := httptest.NewRequest("GET", "/api/v1/hello", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	var body map[string]string
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["status"] != "success" {
		t.Errorf("expected status=success, got %q", body["status"])
	}
	if body["message"] != "Welcome to your Golang Fiber API!" {
		t.Errorf("unexpected message: %q", body["message"])
	}
}

func TestHealthCheckHandler(t *testing.T) {
	app := fiber.New()
	app.Get("/api/v1/health", HealthCheckHandler)

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	var body map[string]string
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("expected status=ok, got %q", body["status"])
	}
}

func TestErrorResponse(t *testing.T) {
	app := fiber.New()
	app.Get("/api/v1/err", func(c fiber.Ctx) error {
		return ErrorResponse(c, fiber.StatusBadRequest, "invalid input", errors.New("missing field"))
	})

	req := httptest.NewRequest("GET", "/api/v1/err", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	var body map[string]string
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["status"] != "error" || body["message"] != "invalid input" || body["error"] != "missing field" {
		t.Errorf("unexpected error payload: %+v", body)
	}
}
