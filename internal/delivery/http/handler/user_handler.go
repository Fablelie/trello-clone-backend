package handler

import (
	"time"

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

func (h *UserHandler) Login(c fiber.Ctx) error {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req request
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	// Call usecase to get the token
	token, err := h.UserUsecase.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": err.Error()})
	}

	// Create and send cookie to the browser
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "login success"})
}
