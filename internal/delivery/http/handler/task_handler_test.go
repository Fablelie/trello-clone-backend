package handler_test

import (
	"bytes"
	"encoding/json"
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

// Create a Mock for TaskUsecase
type mockTaskUsecase struct {
	mock.Mock
}

func (m *mockTaskUsecase) CreateTask(actorID uuid.UUID, task *domain.Task) error {
	args := m.Called(actorID, task)
	return args.Error(0)
}

func (m *mockTaskUsecase) GetTasksByProject(actorID uuid.UUID, projectID uuid.UUID) ([]domain.Task, error) {
	args := m.Called(actorID, projectID)
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *mockTaskUsecase) MoveTask(actorID uuid.UUID, taskID uuid.UUID, targetColumnID uuid.UUID) error {
	args := m.Called(actorID, taskID, targetColumnID)
	return args.Error(0)
}

func (m *mockTaskUsecase) UpdateTask(actorID uuid.UUID, task *domain.Task) error {
	args := m.Called(actorID, task)
	return args.Error(0)
}

func (m *mockTaskUsecase) DeleteTask(actorID uuid.UUID, taskID uuid.UUID) error {
	args := m.Called(actorID, taskID)
	return args.Error(0)
}

func (m *mockTaskUsecase) AssignMembers(actorID uuid.UUID, taskID uuid.UUID, emails []string) error {
	args := m.Called(actorID, taskID, emails)
	return args.Error(0)
}

func (m *mockTaskUsecase) RemoveMember(actorID uuid.UUID, taskID uuid.UUID, email string) error {
	args := m.Called(actorID, taskID, email)
	return args.Error(0)
}

func (m *mockTaskUsecase) CheckUserInTask(taskID uuid.UUID, userID uuid.UUID) (bool, error) {
	args := m.Called(taskID, userID)
	return args.Bool(0), args.Error(1)
}

// ==========================================
// Test Cases for TaskHandler
// ==========================================

func TestMoveTaskHandler_Success(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(mockTaskUsecase)
	h := handler.NewTaskHandler(mockUsecase)

	// Set Route: /api/v1/projects/:project_id/tasks/:id/move
	app.Patch("/projects/:project_id/tasks/:id/move", func(c fiber.Ctx) error {
		c.Locals("actor_id", uuid.MustParse("76ca249b-2773-47b9-abea-e2703bcec4c4"))
		return h.MoveTask(c)
	})

	projectID := uuid.New()
	taskID := uuid.New()
	targetColumnID := uuid.New()

	mockUsecase.On("MoveTask", mock.Anything, taskID, targetColumnID).Return(nil)

	reqBody, _ := json.Marshal(map[string]string{
		"target_column_id": targetColumnID.String(),
	})

	url := "/projects/" + projectID.String() + "/tasks/" + taskID.String() + "/move"
	req := httptest.NewRequest(http.MethodPatch, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Act
	resp, err := app.Test(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.Equal(t, "task moved successfully", result["message"])

	mockUsecase.AssertExpectations(t)
}
