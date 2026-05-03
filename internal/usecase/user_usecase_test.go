package usecase_test

import (
	"errors"
	"testing"

	"github.com/fablelie/trello-clone-backend/internal/domain"
	"github.com/fablelie/trello-clone-backend/internal/usecase"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// Create a Mock for UserRepository
type mockUserRepo struct {
	mock.Mock
}

func (m *mockUserRepo) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *mockUserRepo) GetByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserRepo) GetByID(id uuid.UUID) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// ==========================================
// Test Cases for Register
// ==========================================

func TestRegister_Success(t *testing.T) {
	repo := new(mockUserRepo)
	jwtSecret := "test_secret_key"
	u := usecase.NewUserUsecase(repo, jwtSecret)

	repo.On("GetByEmail", "new@test.com").Return(nil, errors.New("record not found"))
	repo.On("Create", mock.Anything).Return(nil)

	err := u.Register("new@test.com", "password123", "Pawat")

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	repo := new(mockUserRepo)
	jwtSecret := "test_secret_key"
	u := usecase.NewUserUsecase(repo, jwtSecret)

	existingUser := &domain.User{
		ID:    uuid.New(),
		Email: "existing@test.com",
	}

	// Mock had exist this email in record
	repo.On("GetByEmail", "existing@test.com").Return(existingUser, nil)

	err := u.Register("existing@test.com", "password123", "Pawat")

	assert.Error(t, err)
	assert.Equal(t, "email already exists", err.Error())
	repo.AssertNotCalled(t, "Create", mock.Anything)
}

// ==========================================
// Test Cases for Login
// ==========================================

func TestLogin_Success(t *testing.T) {
	repo := new(mockUserRepo)
	jwtSecret := "test_secret_key"
	u := usecase.NewUserUsecase(repo, jwtSecret)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	dbUser := &domain.User{
		ID:       uuid.New(),
		Email:    "login@test.com",
		Password: string(hashedPassword),
	}

	repo.On("GetByEmail", "login@test.com").Return(dbUser, nil)

	token, err := u.Login("login@test.com", "password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	repo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	repo := new(mockUserRepo)
	jwtSecret := "test_secret_key"
	u := usecase.NewUserUsecase(repo, jwtSecret)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), bcrypt.DefaultCost)
	dbUser := &domain.User{
		ID:       uuid.New(),
		Email:    "login@test.com",
		Password: string(hashedPassword),
	}

	repo.On("GetByEmail", "login@test.com").Return(dbUser, nil)

	token, err := u.Login("login@test.com", "wrong_password")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "invalid email or password", err.Error())
}
