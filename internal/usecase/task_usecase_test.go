package usecase_test

import (
	"errors"
	"testing"

	"github.com/fablelie/trello-clone-backend/internal/domain"
	"github.com/fablelie/trello-clone-backend/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 1. Create Mock for TaskRepository
type mockTaskRepo struct {
	mock.Mock
}

func (m *mockTaskRepo) Create(task *domain.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *mockTaskRepo) Update(task *domain.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *mockTaskRepo) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockTaskRepo) GetByID(id uuid.UUID) (*domain.Task, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Task), args.Error(1)
}

func (m *mockTaskRepo) GetByProjectID(projectID uuid.UUID) ([]domain.Task, error) {
	args := m.Called(projectID)
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *mockTaskRepo) AddMember(taskID uuid.UUID, userID uuid.UUID) error {
	args := m.Called(taskID, userID)
	return args.Error(0)
}

func (m *mockTaskRepo) AddMembers(taskID uuid.UUID, userIDs []uuid.UUID) error {
	args := m.Called(taskID, userIDs)
	return args.Error(0)
}

func (m *mockTaskRepo) RemoveMember(taskID uuid.UUID, userID uuid.UUID) error {
	args := m.Called(taskID, userID)
	return args.Error(0)
}

func (m *mockTaskRepo) IsMember(taskID uuid.UUID, userID uuid.UUID) (bool, error) {
	args := m.Called(taskID, userID)
	return args.Bool(0), args.Error(1)
}

// ==========================================
// Test Cases for CreateTask
// ==========================================

func TestCreateTask_Success(t *testing.T) {
	tRepo := new(mockTaskRepo)
	pRepo := new(mockProjectRepo)
	uRepo := new(mockUserRepo)
	u := usecase.NewTaskUsecase(tRepo, pRepo, uRepo)

	actorID := uuid.New()
	projectID := uuid.New()

	member := &domain.ProjectMember{ProjectID: projectID, UserID: actorID, Role: "Member"}
	pRepo.On("GetMember", projectID, actorID).Return(member, nil)

	tRepo.On("Create", mock.Anything).Return(nil)

	task := &domain.Task{
		ProjectID:   projectID,
		Subject:     "Test Task",
		Description: "Hello unit test",
		Column:      &domain.Column{ColumnID: uuid.New()},
	}

	err := u.CreateTask(actorID, task)

	assert.NoError(t, err)
	tRepo.AssertExpectations(t)
}

func TestCreateTask_NotAMember(t *testing.T) {
	tRepo := new(mockTaskRepo)
	pRepo := new(mockProjectRepo)
	uRepo := new(mockUserRepo)
	u := usecase.NewTaskUsecase(tRepo, pRepo, uRepo)

	actorID := uuid.New()
	projectID := uuid.New()

	pRepo.On("GetMember", projectID, actorID).Return(nil, errors.New("not found"))

	task := &domain.Task{
		ProjectID: projectID,
		Subject:   "Test Task",
	}

	err := u.CreateTask(actorID, task)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied")
	tRepo.AssertNotCalled(t, "Create", mock.Anything)
}

// ==========================================
// Test Cases for UpdateTask
// ==========================================

func TestUpdateTask_Success_ByCreator(t *testing.T) {
	// Arrange
	tRepo := new(mockTaskRepo)
	pRepo := new(mockProjectRepo)
	uRepo := new(mockUserRepo)
	u := usecase.NewTaskUsecase(tRepo, pRepo, uRepo)

	actorID := uuid.New()
	taskID := uuid.New()
	projectID := uuid.New()

	existingTask := &domain.Task{
		TaskID:    taskID,
		ProjectID: projectID,
		Subject:   "Old Task",
		Assigner:  &domain.User{ID: actorID},
	}

	tRepo.On("GetByID", taskID).Return(existingTask, nil)

	tRepo.On("Update", mock.Anything).Return(nil)

	updateTask := &domain.Task{
		TaskID:  taskID,
		Subject: "Updated Task Name",
	}

	err := u.UpdateTask(actorID, updateTask)

	assert.NoError(t, err)
	tRepo.AssertExpectations(t)
}

func TestUpdateTask_PermissionDenied_NotCreatorNorAdmin(t *testing.T) {
	tRepo := new(mockTaskRepo)
	pRepo := new(mockProjectRepo)
	uRepo := new(mockUserRepo)
	u := usecase.NewTaskUsecase(tRepo, pRepo, uRepo)

	actorID := uuid.New()
	taskID := uuid.New()
	projectID := uuid.New()
	creatorID := uuid.New()

	existingTask := &domain.Task{
		TaskID:    taskID,
		ProjectID: projectID,
		Assigner:  &domain.User{ID: creatorID},
	}

	tRepo.On("GetByID", taskID).Return(existingTask, nil)

	// Mock actorID is not a Admin just a Member in project
	member := &domain.ProjectMember{ProjectID: projectID, UserID: actorID, Role: "Member"}
	pRepo.On("GetMember", projectID, actorID).Return(member, nil)

	updateTask := &domain.Task{
		TaskID: taskID,
	}

	err := u.UpdateTask(actorID, updateTask)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied")
	tRepo.AssertNotCalled(t, "Update", mock.Anything)
}
