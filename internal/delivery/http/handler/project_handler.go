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
	// 1. ดึง user_id จาก Locals (ที่ได้จาก AuthMiddleware)
	// 1. Extract user_id from Locals (set by AuthMiddleware)
	actorIDStr := c.Locals("user_id").(string)
	actorID, err := uuid.Parse(actorIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid user id in token"})
	}

	// 2. รับข้อมูลชื่อโปรเจกต์จาก Body
	// 2. Parse project name from request body
	type request struct {
		Name string `json:"project_name"`
	}
	var req request
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	// 3. ส่งไปให้ Usecase จัดการสร้างโปรเจกต์
	// 3. Send to Usecase to handle project creation
	err = h.projectUsecase.CreateProject(req.Name, actorID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "project created successfully",
	})
}

// AddMember handles adding a new member to the project
// AddMember จัดการเพิ่มสมาชิกใหม่เข้าไปในโปรเจกต์
func (h *ProjectHandler) AddMember(c fiber.Ctx) error {
	actorIDStr := c.Locals("user_id").(string)
	actorID := uuid.MustParse(actorIDStr)

	var req domain.ProjectMember
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	// เรียกใช้ Usecase (ซึ่งจะมีการเช็คสิทธิ์ Admin ข้างใน)
	// Call Usecase (which includes Admin role check inside)
	err := h.projectUsecase.AddMember(actorID, &req)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "member added successfully"})
}
