package server

import (
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"templatenextjs_golangserver/src/handlers"
	appmiddleware "templatenextjs_golangserver/src/middleware"
	"templatenextjs_golangserver/src/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

// WebuiFS is set by the main entry point to expose Next.js static assets
var WebuiFS fs.FS

// PublicAddress is the address the full site (API + web UI) listens on, from PORT.
func PublicAddress() string {
	appPort := os.Getenv("PORT")
	if appPort == "" {
		appPort = "8080"
	}
	if !strings.HasPrefix(appPort, ":") {
		appPort = ":" + appPort
	}
	return appPort
}

// NewApp constructs and configures the Fiber API application.
func NewApp() *fiber.App {
	app := fiber.New()

	app.Use("/api/v1", appmiddleware.VerboseLogger())

	// Register API endpoints
	app.Get("/api/v1/hello", handlers.GetMessageHandler)
	app.Get("/api/v1/health", handlers.HealthCheckHandler)

	return app
}

// RunServer starts either the dev reverse proxy (DEV=true) or the production
// single-binary Fiber server serving both /api/v1 and embedded Next.js assets.
func RunServer() {
	appPort := PublicAddress()
	app := NewApp()

	if os.Getenv("DEV") == "true" {
		apiPort := os.Getenv("API_PORT")
		if apiPort == "" {
			apiPort = "8082"
		}
		apiPortRaw := strings.TrimPrefix(apiPort, ":")
		apiPortWithColon := ":" + apiPortRaw

		webPort := os.Getenv("WEB_PORT")
		if webPort == "" {
			webPort = os.Getenv("NEXTJS_PORT")
		}
		if webPort == "" {
			webPort = "3000"
		}
		webPortRaw := strings.TrimPrefix(webPort, ":")

		// In dev mode, start the Fiber API server on the internal port in background
		go func() {
			utils.LogInfo("Starting Fiber API server on internal port %s...", apiPortWithColon)
			if err := app.Listen(apiPortWithColon); err != nil {
				utils.LogErr("Fiber API server failed: %v", err)
			}
		}()

		nextjsURL, err := url.Parse("http://localhost:" + webPortRaw)
		if err != nil {
			utils.LogErr("Failed to parse Next.js target URL: %v", err)
		}
		fiberURL, err := url.Parse("http://localhost:" + apiPortRaw)
		if err != nil {
			utils.LogErr("Failed to parse Fiber target URL: %v", err)
		}

		nextjsProxy := httputil.NewSingleHostReverseProxy(nextjsURL)
		fiberProxy := httputil.NewSingleHostReverseProxy(fiberURL)

		devMux := http.NewServeMux()
		devMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api/v1") {
				fiberProxy.ServeHTTP(w, r)
			} else {
				nextjsProxy.ServeHTTP(w, r)
			}
		})

		utils.LogInfo("Starting Dev Proxy server on port %s (routing API to %s and Web to :%s)...", appPort, apiPortWithColon, webPortRaw)
		if err := http.ListenAndServe(appPort, devMux); err != nil {
			utils.LogErr("Dev Proxy server failed: %v", err)
		}
	} else {
		// Production mode: run Fiber directly on main appPort with embedded filesystem
		if WebuiFS != nil {
			app.Use(static.New("", static.Config{
				FS:         WebuiFS,
				IndexNames: []string{"index.html"},
				NotFoundHandler: func(c fiber.Ctx) error {
					data, err := fs.ReadFile(WebuiFS, "404/index.html")
					if err == nil {
						c.Type("html")
						return c.Status(fiber.StatusNotFound).Send(data)
					}
					return c.Status(fiber.StatusNotFound).SendString("Frontend assets not loaded")
				},
			}))
		} else {
			app.Use(func(c fiber.Ctx) error {
				return c.Status(fiber.StatusNotFound).SendString("Frontend assets not loaded")
			})
		}

		utils.LogInfo("Starting Production Fiber server on port %s...", appPort)
		log.Fatal(app.Listen(appPort))
	}
}
