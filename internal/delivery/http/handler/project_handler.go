package handler

import (
	"github.com/fablelie/trello-clone-backend/internal/domain"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	projectUsecase domain.ProjectUsecase
}

// NewProjectHandler creates a new instance of project handler
func NewProjectHandler(u domain.ProjectUsecase) *ProjectHandler {
	return &ProjectHandler{
		projectUsecase: u,
	}
}

// Create handles the request to create a new project
func (h *ProjectHandler) Create(c fiber.Ctx) error {
	// Get actor_id directly from Locals (already parsed as uuid.UUID by Middleware)
	actorID := c.Locals("actor_id").(uuid.UUID)

	// Parse project name from request body
	type request struct {
		Name string `json:"project_name"`
	}
	var req request
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	// Call Usecase to handle project creation and default columns
	err := h.projectUsecase.CreateProject(req.Name, actorID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "project created successfully",
	})
}

// AddMember handles adding a new member to the project
func (h *ProjectHandler) AddMember(c fiber.Ctx) error {
	actorID := c.Locals("actor_id").(uuid.UUID)

	var req domain.ProjectMember
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	// Usecase will verify if the actorID has Admin role for this project
	err := h.projectUsecase.AddMember(actorID, &req)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "member added successfully"})
}
