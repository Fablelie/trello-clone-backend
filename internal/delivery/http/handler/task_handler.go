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

	projectID, err := uuid.Parse(c.Params("project_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid project id"})
	}

	var req domain.Task
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	req.ProjectID = projectID

	// Call Usecase (membership check is handled inside)
	if err := h.taskUsecase.CreateTask(actorID, &req); err != nil {
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

// UpdateTask handles updating the task supject, description, column (status)
func (h *TaskHandler) UpdateTask(c fiber.Ctx) error {
	actorID := c.Locals("actor_id").(uuid.UUID)
	taskID, _ := uuid.Parse(c.Params("id"))

	var req domain.Task
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	req.TaskID = taskID

	if err := h.taskUsecase.UpdateTask(actorID, &req); err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "task updated successfuly"})
}

// AssignMembers handles assigning a users to a task
func (h *TaskHandler) AssignMembers(c fiber.Ctx) error {
	actorID := c.Locals("actor_id").(uuid.UUID)
	taskID, _ := uuid.Parse(c.Params("id"))

	type request struct {
		Emails []string `json:"emails"`
	}
	var req request
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	err := h.taskUsecase.AssignMembers(actorID, taskID, req.Emails)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "members assigned successfully"})
}

// RemoveMember handles remove a user from a task
func (h *TaskHandler) RemoveMember(c fiber.Ctx) error {
	actorID := c.Locals("actor_id").(uuid.UUID)
	taskID, _ := uuid.Parse(c.Params("id"))

	type request struct {
		Email string `json:"email"`
	}
	var req request
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	err := h.taskUsecase.RemoveMember(actorID, taskID, req.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "member remover successfully"})
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

func (h *TaskHandler) GetByProjectID(c fiber.Ctx) error {
	actorID := c.Locals("actor_id").(uuid.UUID)
	projectID, err := uuid.Parse(c.Params("project_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid project id"})
	}

	tasks, err := h.taskUsecase.GetTasksByProject(actorID, projectID)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(tasks)
}
