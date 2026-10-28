package server

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v3"
)

func TestPublicAddress(t *testing.T) {
	orig := os.Getenv("PORT")
	defer os.Setenv("PORT", orig)

	os.Unsetenv("PORT")
	if got := PublicAddress(); got != ":8080" {
		t.Errorf("PublicAddress() default = %q, want \":8080\"", got)
	}

	os.Setenv("PORT", "9090")
	if got := PublicAddress(); got != ":9090" {
		t.Errorf("PublicAddress() with 9090 = %q, want \":9090\"", got)
	}

	os.Setenv("PORT", ":3001")
	if got := PublicAddress(); got != ":3001" {
		t.Errorf("PublicAddress() with :3001 = %q, want \":3001\"", got)
	}
}

func TestNewAppRoutes(t *testing.T) {
	app := NewApp()

	for _, path := range []string{"/api/v1/hello", "/api/v1/health"} {
		req := httptest.NewRequest("GET", path, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("GET %s failed: %v", path, err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Errorf("GET %s status = %d, want 200", path, resp.StatusCode)
		}
	}
}
