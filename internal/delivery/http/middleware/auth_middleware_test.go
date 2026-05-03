package middleware_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/fablelie/trello-clone-backend/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ==========================================
// Auth Middleware Tests
// ==========================================

func TestAuthMiddleware_ValidToken_Success(t *testing.T) {
	// Setup
	app := fiber.New()
	secret := "test-secret-key-very-secure"

	// Create a valid JWT token
	userID := uuid.New()
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	// Setup routes
	app.Post("/protected", middleware.AuthMiddleware(secret), func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
	})

	// Create request with valid token in cookie
	req := createRequestWithCookie(tokenString)

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)

	if assert.NotNil(t, resp) {
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	}
}

func TestAuthMiddleware_MissingToken_Unauthorized(t *testing.T) {
	// Setup
	app := fiber.New()
	secret := "test-secret-key-very-secure"

	// Setup routes without token
	app.Post("/protected", middleware.AuthMiddleware(secret), func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
	})

	// Create request WITHOUT token cookie
	req := createRequest()

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)

	if assert.NotNil(t, resp) {
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	}
}

func TestAuthMiddleware_InvalidToken_Unauthorized(t *testing.T) {
	// Setup
	app := fiber.New()
	secret := "test-secret-key-very-secure"
	invalidToken := "invalid.token.here"

	// Setup routes
	app.Post("/protected", middleware.AuthMiddleware(secret), func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
	})

	// Create request with invalid token
	req := createRequestWithCookie(invalidToken)

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAuthMiddleware_ExpiredToken_Unauthorized(t *testing.T) {
	// Setup
	app := fiber.New()
	secret := "test-secret-key-very-secure"

	// Create an expired JWT token
	userID := uuid.New()
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	// Setup routes
	app.Post("/protected", middleware.AuthMiddleware(secret), func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
	})

	// Create request with expired token
	req := createRequestWithCookie(tokenString)

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAuthMiddleware_WrongSecret_Unauthorized(t *testing.T) {
	// Setup
	app := fiber.New()
	correctSecret := "correct-secret"
	wrongSecret := "wrong-secret"

	// Create token signed with correct secret
	userID := uuid.New()
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(correctSecret))
	assert.NoError(t, err)

	// Setup routes with wrong secret
	app.Post("/protected", middleware.AuthMiddleware(wrongSecret), func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
	})

	// Create request with token signed by different secret
	req := createRequestWithCookie(tokenString)

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestAuthMiddleware_MalformedToken_Unauthorized(t *testing.T) {
	// Setup
	app := fiber.New()
	secret := "test-secret-key-very-secure"
	malformedToken := "malformed.token"

	// Setup routes
	app.Post("/protected", middleware.AuthMiddleware(secret), func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
	})

	// Create request with malformed token
	req := createRequestWithCookie(malformedToken)

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// ==========================================
// User Context Middleware Tests
// ==========================================

func TestUserContextMiddleware_ValidUserContext_Success(t *testing.T) {
	// Setup
	app := fiber.New()
	secret := "test-secret-key-very-secure"
	userID := uuid.New()

	// Create a valid JWT token
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	// Setup routes with both middleware
	app.Post("/protected",
		middleware.AuthMiddleware(secret),
		middleware.UserContextMiddleware(),
		func(c fiber.Ctx) error {
			actorID := c.Locals("actor_id")
			if actorID == nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "actor_id not set"})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"actor_id": actorID})
		},
	)

	// Create request
	req := createRequestWithCookie(tokenString)

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestUserContextMiddleware_NoUserContext_BadRequest(t *testing.T) {
	// Setup
	app := fiber.New()

	// Setup route without auth middleware (so user context is nil)
	app.Post("/protected",
		middleware.UserContextMiddleware(),
		func(c fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	// Create request without token
	req := createRequest()

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestUserContextMiddleware_InvalidClaimsFormat_Unauthorized(t *testing.T) {
	// Setup
	app := fiber.New()
	secret := "test-secret-key-very-secure"

	// Create token with non-MapClaims format (if possible)
	// For this test, we'll create a token with missing user_id
	claims := jwt.MapClaims{
		"email": "test@example.com", // Missing user_id
		"exp":   time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	// Setup routes
	app.Post("/protected",
		middleware.AuthMiddleware(secret),
		middleware.UserContextMiddleware(),
		func(c fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	// Create request
	req := createRequestWithCookie(tokenString)

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestUserContextMiddleware_MissingUserID_Unauthorized(t *testing.T) {
	// Setup
	app := fiber.New()
	secret := "test-secret-key-very-secure"

	// Create token without user_id
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	// Setup routes
	app.Post("/protected",
		middleware.AuthMiddleware(secret),
		middleware.UserContextMiddleware(),
		func(c fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	// Create request
	req := createRequestWithCookie(tokenString)

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestUserContextMiddleware_InvalidUserIDFormat_Recovery(t *testing.T) {
	// Setup
	app := fiber.New()
	secret := "test-secret-key-very-secure"

	app.Post("/protected",
		middleware.AuthMiddleware(secret),
		middleware.UserContextMiddleware(),
		func(c fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "success"})
		},
	)

	claims := jwt.MapClaims{
		"user_id": "not-a-valid-uuid",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	req := createRequestWithCookie(tokenString)

	resp, err := app.Test(req)

	assert.NoError(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	}
}

func TestUserContextMiddleware_ActorIDStoredCorrectly(t *testing.T) {
	// Setup
	app := fiber.New()
	secret := "test-secret-key-very-secure"
	expectedUserID := uuid.New()

	// Create a valid JWT token with user_id
	claims := jwt.MapClaims{
		"user_id": expectedUserID.String(),
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	// Variable to capture actor_id
	var capturedActorID uuid.UUID

	// Setup routes
	app.Post("/protected",
		middleware.AuthMiddleware(secret),
		middleware.UserContextMiddleware(),
		func(c fiber.Ctx) error {
			actorID, ok := c.Locals("actor_id").(uuid.UUID)
			if !ok {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "actor_id type assertion failed"})
			}
			capturedActorID = actorID
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"actor_id": actorID.String()})
		},
	)

	// Create request
	req := createRequestWithCookie(tokenString)

	// Execute
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedUserID, capturedActorID)
}

// ==========================================
// Helper Functions
// ==========================================

// createRequest creates a basic POST request
func createRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "/protected", nil)
	req.Host = "example.com"
	req.Header.Set("User-Agent", "fiber-test")
	return req
}

// createRequestWithCookie creates a POST request with JWT cookie
func createRequestWithCookie(token string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, "/protected", nil)

	req.Host = "example.com"

	// Setup Cookie for JWT token
	req.AddCookie(&http.Cookie{
		Name:  "jwt",
		Value: token,
	})

	return req
}
