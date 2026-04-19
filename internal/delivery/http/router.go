package http

import (
	"github.com/fablelie/trello-clone-backend/internal/delivery/http/handler"
	"github.com/gofiber/fiber/v3"
)

// SetupRouter manages all API routes
func SetupRouter(app *fiber.App, userHandler *handler.UserHandler) {
	// Create a group for API v1
	v1 := app.Group("/api/v1")

	// Auth Routes
	auth := v1.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)
}
