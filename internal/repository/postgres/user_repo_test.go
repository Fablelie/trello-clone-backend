package postgres_test

import (
	"os"
	"testing"

	postgresRepo "github.com/fablelie/trello-clone-backend/internal/repository/postgres"

	"github.com/fablelie/trello-clone-backend/internal/domain"
	"github.com/fablelie/trello-clone-backend/internal/infrastructure/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var testDB *gorm.DB

// TestMain runs before all test cases to prepare the DB connection
func TestMain(m *testing.M) {
	_ = godotenv.Load("../../../../.env")

	if os.Getenv("DB_HOST") != "" {
		testDB = database.NewPostgresDB(
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)
	}

	code := m.Run()

	os.Exit(code)
}

// ==========================================
// User Repository Tests
// ==========================================

func TestUserRepository_Create_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	repo := postgresRepo.NewUserRepository(testDB)

	// Arrange
	testUser := &domain.User{
		ID:       uuid.New(),
		Name:     "John Doe",
		Email:    "john_" + uuid.New().String() + "@example.com",
		Password: "hashed_password_123",
	}

	// Act
	err := repo.Create(testUser)

	// Assert
	assert.NoError(t, err)
	assert.NotZero(t, testUser.CreatedAt)

	// Cleanup
	testDB.Exec("DELETE FROM users WHERE id = ?", testUser.ID)
}

func TestUserRepository_Create_And_GetByEmail_Integration(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	repo := postgresRepo.NewUserRepository(testDB)

	testUser := &domain.User{
		ID:       uuid.New(),
		Name:     "Test Integration User",
		Email:    "integration_test_" + uuid.New().String() + "@example.com",
		Password: "hashed_password_123",
	}

	err := repo.Create(testUser)

	assert.NoError(t, err)

	fetchedUser, err := repo.GetByEmail(testUser.Email)

	assert.NoError(t, err)
	assert.NotNil(t, fetchedUser)
	assert.Equal(t, testUser.ID, fetchedUser.ID)
	assert.Equal(t, testUser.Name, fetchedUser.Name)
	assert.Equal(t, testUser.Email, fetchedUser.Email)
	assert.Equal(t, testUser.Password, fetchedUser.Password)

	// Cleanup
	testDB.Exec("DELETE FROM users WHERE id = ?", testUser.ID)
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	repo := postgresRepo.NewUserRepository(testDB)

	// Act
	user, err := repo.GetByEmail("nonexistent_" + uuid.New().String() + "@example.com")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserRepository_GetByID_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	repo := postgresRepo.NewUserRepository(testDB)

	// Arrange: Create a user first
	testUser := &domain.User{
		ID:       uuid.New(),
		Name:     "Jane Doe",
		Email:    "jane_" + uuid.New().String() + "@example.com",
		Password: "hashed_password_456",
	}

	err := repo.Create(testUser)
	assert.NoError(t, err)

	// Act
	fetchedUser, err := repo.GetByID(testUser.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, fetchedUser)
	assert.Equal(t, testUser.ID, fetchedUser.ID)
	assert.Equal(t, testUser.Name, fetchedUser.Name)
	assert.Equal(t, testUser.Email, fetchedUser.Email)

	// Cleanup
	testDB.Exec("DELETE FROM users WHERE id = ?", testUser.ID)
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	repo := postgresRepo.NewUserRepository(testDB)

	// Act
	user, err := repo.GetByID(uuid.New())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, user)
}
