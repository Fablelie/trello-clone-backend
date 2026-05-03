package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fablelie/trello-clone-backend/internal/delivery/http/handler"
	"github.com/fablelie/trello-clone-backend/internal/domain"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Create a Mock for ProjectUsecase
type mockProjectUsecase struct {
	mock.Mock
}

func (m *mockProjectUsecase) CreateProject(name string, ownerID uuid.UUID) error {
	args := m.Called(name, ownerID)
	return args.Error(0)
}

func (m *mockProjectUsecase) UpdateColumn(actorID uuid.UUID, column *domain.Column) error {
	args := m.Called(actorID, column)
	return args.Error(0)
}

func (m *mockProjectUsecase) AddMembers(actorID uuid.UUID, projectID uuid.UUID, members []domain.AddMemberRequest) error {
	args := m.Called(actorID, projectID, members)
	return args.Error(0)
}

func (m *mockProjectUsecase) RemoveMember(actorID uuid.UUID, projectID uuid.UUID, userID uuid.UUID) error {
	args := m.Called(actorID, projectID, userID)
	return args.Error(0)
}

func (m *mockProjectUsecase) GetMyProjects(actorID uuid.UUID) ([]domain.Project, error) {
	args := m.Called(actorID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Project), args.Error(1)
}

func (m *mockProjectUsecase) GetProjectByID(actorID uuid.UUID, projectID uuid.UUID) (*domain.Project, error) {
	args := m.Called(actorID, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

// ==========================================
// Test Cases
// ==========================================

func TestCreateProjectHandler_Success(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(mockProjectUsecase)
	h := handler.NewProjectHandler(mockUsecase)

	app.Post("/projects", func(c fiber.Ctx) error {
		c.Locals("actor_id", uuid.MustParse("76ca249b-2773-47b9-abea-e2703bcec4c4"))
		return h.Create(c)
	})

	mockUsecase.On("CreateProject", "My First Board", mock.Anything).Return(nil)

	reqBody, _ := json.Marshal(map[string]string{
		"project_name": "My First Board",
	})

	// Cretae HTTP Request
	req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var result map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "project created successfully", result["message"])

	mockUsecase.AssertExpectations(t)
}

func TestCreateProjectHandler_InternalError(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(mockProjectUsecase)
	h := handler.NewProjectHandler(mockUsecase)

	app.Post("/projects", func(c fiber.Ctx) error {
		c.Locals("actor_id", uuid.New())
		return h.Create(c)
	})

	mockUsecase.On("CreateProject", mock.Anything, mock.Anything).Return(errors.New("db connection timeout"))

	reqBody, _ := json.Marshal(map[string]string{
		"project_name": "Fail Board",
	})

	req := httptest.NewRequest(http.MethodPost, "/projects", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	mockUsecase.AssertExpectations(t)
}

func TestAddMembersHandler_Success(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(mockProjectUsecase)
	h := handler.NewProjectHandler(mockUsecase)

	// Set URL Path to match with router (/api/v1/projects/:id/members)
	app.Post("/projects/:id/members", func(c fiber.Ctx) error {
		c.Locals("actor_id", uuid.MustParse("76ca249b-2773-47b9-abea-e2703bcec4c4"))
		return h.AddMembers(c)
	})

	projectID := uuid.New()
	reqBody, _ := json.Marshal([]domain.AddMemberRequest{
		{Email: "dev1@test.com", Role: "Member", Position: "Developer"},
		{Email: "dev2@test.com", Role: "Admin", Position: "DevOps"},
	})

	mockUsecase.On("AddMembers", mock.Anything, projectID, mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/projects/"+projectID.String()+"/members", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "members added successfully", result["message"])

	mockUsecase.AssertExpectations(t)
}
