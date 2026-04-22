package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func UserContextMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		userToken := c.Locals("user").(*jwt.Token)
		claims := userToken.Claims.(jwt.MapClaims)

		// เก็บ ID ที่เป็น UUID ไว้ให้เลย
		c.Locals("actor_id", uuid.MustParse(claims["user_id"].(string)))
		return c.Next()
	}
}
