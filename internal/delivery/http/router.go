package http

import (
	"github.com/fablelie/trello-clone-backend/internal/delivery/http/handler"
	"github.com/fablelie/trello-clone-backend/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v3"
)

// SetupRouter manages all API routes
func SetupRouter(app *fiber.App, userHandler *handler.UserHandler, projectHandler *handler.ProjectHandler, taskHandler *handler.TaskHandler, secret string) {
	// Create a group for API api
	api := app.Group("/api/v1")

	// Auth Routes
	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	protected := api.Group("", middleware.AuthMiddleware(secret), middleware.UserContextMiddleware())

	// project
	projects := protected.Group("/projects")
	projects.Get("/", projectHandler.GetAll)
	projects.Get("/:id", projectHandler.GetByID)
	projects.Post("/", projectHandler.Create)
	projects.Post("/:id/members", projectHandler.AddMembers)

	// task
	tasks := protected.Group("/tasks")
	tasks.Post("/", taskHandler.Create)
	tasks.Patch("/:id/move", taskHandler.MoveTask)
	tasks.Post("/:id/assign", taskHandler.AssignMember)
	tasks.Delete("/:id", taskHandler.DeleteTask)
}
