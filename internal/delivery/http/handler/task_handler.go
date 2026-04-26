package handler

import (
	"github.com/fablelie/trello-clone-backend/internal/domain"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type TaskHandler struct {
	taskUsecase domain.TaskUsecase
}

// NewTaskHandler creates a new instance of task handler
func NewTaskHandler(u domain.TaskUsecase) *TaskHandler {
	return &TaskHandler{
		taskUsecase: u,
	}
}

// Create handles the request to create a new task
func (h *TaskHandler) Create(c fiber.Ctx) error {
	actorID := c.Locals("actor_id").(uuid.UUID)

	var req domain.Task
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	// Call Usecase (membership check is handled inside)
	err := h.taskUsecase.CreateTask(actorID, &req)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "task created successfully"})
}

// MoveTask handles updating the task's column (status)
func (h *TaskHandler) MoveTask(c fiber.Ctx) error {
	actorID := c.Locals("actor_id").(uuid.UUID)

	taskID, _ := uuid.Parse(c.Params("id"))

	type request struct {
		TargetColumnID uuid.UUID `json:"target_column_id"`
	}
	var req request
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	err := h.taskUsecase.MoveTask(actorID, taskID, req.TargetColumnID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "task moved successfully"})
}

// AssignMember handles assigning a user to a task
func (h *TaskHandler) AssignMember(c fiber.Ctx) error {
	actorID := c.Locals("actor_id").(uuid.UUID)
	taskID, _ := uuid.Parse(c.Params("id"))

	type request struct {
		UserID uuid.UUID `json:"user_id"`
	}
	var req request
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	h.taskUsecase.AssignMember(actorID, taskID, req.UserID)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "member assigned"})
}

// DeleteTask handles task removal
func (h *TaskHandler) DeleteTask(c fiber.Ctx) error {
	actorID := c.Locals("actor_id").(uuid.UUID)
	taskID, _ := uuid.Parse(c.Params("id"))

	err := h.taskUsecase.DeleteTask(actorID, taskID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "task deleted"})
}
