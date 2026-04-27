package middleware

import (
	"fmt"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func UserContextMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		fmt.Printf("Cookie 'jwt' received: %s\n", c.Cookies(("jwt")))

		user := jwtware.FromContext(c)
		fmt.Printf("User in Locals: %v\n", user)
		if user == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "unauthorized"})
		}

		claims, ok := user.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid token claims"})
		}

		userIdStr, ok := claims["user_id"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "user_id not found in token"})
		}

		// เก็บ ID ที่เป็น UUID ไว้ให้เลย
		c.Locals("actor_id", uuid.MustParse(userIdStr))
		return c.Next()
	}
}
