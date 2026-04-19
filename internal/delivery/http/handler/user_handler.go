package handler

import (
	"github.com/fablelie/trello-clone-backend/internal/domain"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	UserUsecase domain.UserUsecase
}

// NewUserHandler creates a request handler for user
func NewUserHandler(u domain.UserUsecase) *UserHandler {
	return &UserHandler{
		UserUsecase: u,
	}
}

// Register receives JSON data to register a new user
func (h *UserHandler) Register(c fiber.Ctx) error {
	// Create a struct to parse the request body
	type request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req request
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	// Pass the data to the usecase layer
	err := h.UserUsecase.Register(req.Email, req.Password, req.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user created successfully"})
}

// Login receives credentials and returns a JWT Token
func (h *UserHandler) Login(c fiber.Ctx) error {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req request
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	// Validate credentials and get token from usecase
	token, err := h.UserUsecase.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}
