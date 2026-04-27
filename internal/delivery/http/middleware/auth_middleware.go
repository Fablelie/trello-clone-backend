package middleware

import (
	"fmt"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
)

// AuthMiddleware uses jwtware to validate Token integrity
func AuthMiddleware(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
		Extractor:  extractors.FromCookie("jwt"),
		ErrorHandler: func(c fiber.Ctx, err error) error {
			fmt.Printf("JWT Validation Error: %v\n", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized: " + err.Error(),
			})
		},
	})
}
