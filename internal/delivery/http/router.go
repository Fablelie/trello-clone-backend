package http

import (
	"github.com/fablelie/trello-clone-backend/internal/delivery/http/handler"
	"github.com/fablelie/trello-clone-backend/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v3"
)

// SetupRouter manages all API routes
func SetupRouter(app *fiber.App, userHandler *handler.UserHandler, projectHandler *handler.ProjectHandler, secret string) {
	// Create a group for API api
	api := app.Group("/api/v1")

	// Auth Routes
	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	protected := api.Group("/", middleware.AuthMiddleware(secret), middleware.UserContextMiddleware())

	projects := protected.Group("/projects")
	projects.Post("/", projectHandler.Create)
	projects.Post("/members", projectHandler.AddMember)
}
