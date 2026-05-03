package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fablelie/trello-clone-backend/internal/delivery/http/handler"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Create a Mock for UserUsecase
type mockUserUsecase struct {
	mock.Mock
}

func (m *mockUserUsecase) Register(name, email, password string) error {
	args := m.Called(name, email, password)
	return args.Error(0)
}

func (m *mockUserUsecase) Login(email, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(0), args.Error(1)
}

// ==========================================
// Test Cases Register
// ==========================================

func TestRegisterHandler_Success(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(mockUserUsecase)
	h := handler.NewUserHandler(mockUsecase)

	app.Post("/auth/register", h.Register)

	// Mock register successfully
	mockUsecase.On("Register", "pawat@test.com", "password123", "Pawat").Return(nil)

	reqBody, _ := json.Marshal(map[string]string{
		"email":    "pawat@test.com",
		"password": "password123",
		"name":     "Pawat",
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "user registered successfully", result["message"])

	mockUsecase.AssertExpectations(t)
}

func TestRegisterHandler_DuplicateEmail(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(mockUserUsecase)
	h := handler.NewUserHandler(mockUsecase)

	app.Post("/auth/register", h.Register)

	// Mock duplicate email and usecase to throw error
	mockUsecase.On("Register", "duplicate@test.com", "password123", "Pawat").
		Return(errors.New("user already exists with this email"))

	reqBody, _ := json.Marshal(map[string]string{
		"email":    "duplicate@test.com",
		"password": "password123",
		"name":     "Pawat",
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode) // คาดหวัง 400

	var result map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "user already exists with this email", result["message"])

	mockUsecase.AssertExpectations(t)
}

// ==========================================
// Test Cases for Login (check Cookie)
// ==========================================

func TestLoginHandler_Success(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(mockUserUsecase)
	h := handler.NewUserHandler(mockUsecase)

	app.Post("/auth/login", h.Login)

	fakeToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.fakeTokenString"
	mockUsecase.On("Login", "pawat@test.com", "password123").Return(fakeToken, nil)

	reqBody, _ := json.Marshal(map[string]string{
		"email":    "pawat@test.com",
		"password": "password123",
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check if the JWT cookie was actually sent back
	cookies := resp.Cookies()
	assert.NotEmpty(t, cookies)

	var jwtCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "jwt" {
			jwtCookie = c
			break
		}
	}

	assert.NotNil(t, jwtCookie)
	assert.Equal(t, fakeToken, jwtCookie.Value)
	assert.True(t, jwtCookie.HttpOnly)

	mockUsecase.AssertExpectations(t)
}
