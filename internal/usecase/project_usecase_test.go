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

type mockProjectRepo struct {
	mock.Mock
}

func (m *mockProjectRepo) Create(project *domain.Project) error {
	args := m.Called(project)
	return args.Error(0)
}

func (m *mockProjectRepo) GetByID(id uuid.UUID) (*domain.Project, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Project), args.Error(1)
}

func (m *mockProjectRepo) GetByUserID(userID uuid.UUID) ([]domain.Project, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.Project), args.Error(1)
}

func (m *mockProjectRepo) CreateColumn(column *domain.Column) error {
	args := m.Called(column)
	return args.Error(0)
}

func (m *mockProjectRepo) UpdateColumn(column *domain.Column) error {
	args := m.Called(column)
	return args.Error(0)
}

func (m *mockProjectRepo) DeleteColumn(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockProjectRepo) AddMember(member *domain.ProjectMember) error {
	args := m.Called(member)
	return args.Error(0)
}

func (m *mockProjectRepo) AddMembers(members []domain.ProjectMember) error {
	args := m.Called(members)
	return args.Error(0)
}

func (m *mockProjectRepo) GetMember(projectID uuid.UUID, userID uuid.UUID) (*domain.ProjectMember, error) {
	args := m.Called(projectID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ProjectMember), args.Error(1)
}

func (m *mockProjectRepo) UpdateMember(member *domain.ProjectMember) error {
	args := m.Called(member)
	return args.Error(0)
}

func (m *mockProjectRepo) RemoveMember(projectID, userID uuid.UUID) error {
	args := m.Called(projectID, userID)
	return args.Error(0)
}

type mockProjectUserRepo struct {
	mock.Mock
}

func (m *mockProjectUserRepo) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockProjectUserRepo) GetByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockProjectUserRepo) GetByID(id uuid.UUID) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// ==========================================
// Test Cases for CreateProject
// ==========================================

func TestCreateProject_Success(t *testing.T) {
	pRepo := new(mockProjectRepo)
	uRepo := new(mockProjectUserRepo)
	u := usecase.NewProjectUsecase(pRepo, uRepo)

	ownerID := uuid.New()
	projectName := "New Board"

	pRepo.On("Create", mock.Anything).Return(nil)

	err := u.CreateProject(projectName, ownerID)

	assert.NoError(t, err)
	pRepo.AssertExpectations(t)
}

// ==========================================
// Test Cases for AddMembers (check permission Admin)
// ==========================================

func TestAddMembers_PermissionDenied(t *testing.T) {
	pRepo := new(mockProjectRepo)
	uRepo := new(mockProjectUserRepo)
	u := usecase.NewProjectUsecase(pRepo, uRepo)

	actorID := uuid.New()
	projectID := uuid.New()

	nonAdminMember := &domain.ProjectMember{
		ProjectID: projectID,
		UserID:    actorID,
		Role:      "Member",
		Position:  "Developer",
	}
	pRepo.On("GetMember", projectID, actorID).Return(nonAdminMember, nil)

	req := []domain.AddMemberRequest{
		{Email: "friend@test.com", Role: "Member", Position: "Designer"},
	}

	err := u.AddMembers(actorID, projectID, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied")

	pRepo.AssertNotCalled(t, "AddMembers", mock.Anything)
}

func TestAddMembers_UserNotFound(t *testing.T) {
	pRepo := new(mockProjectRepo)
	uRepo := new(mockProjectUserRepo)
	u := usecase.NewProjectUsecase(pRepo, uRepo)

	actorID := uuid.New()
	projectID := uuid.New()

	adminMember := &domain.ProjectMember{
		ProjectID: projectID,
		UserID:    actorID,
		Role:      "Admin",
	}
	pRepo.On("GetMember", projectID, actorID).Return(adminMember, nil)

	uRepo.On("GetByEmail", "notfound@test.com").Return(nil, errors.New("user not found"))

	req := []domain.AddMemberRequest{
		{Email: "notfound@test.com", Role: "Member", Position: "Developer"},
	}

	// Act
	err := u.AddMembers(actorID, projectID, req)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	pRepo.AssertNotCalled(t, "AddMembers", mock.Anything)
}
