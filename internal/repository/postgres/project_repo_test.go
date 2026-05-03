package postgres_test

import (
	"testing"

	"github.com/fablelie/trello-clone-backend/internal/domain"
	postgresRepo "github.com/fablelie/trello-clone-backend/internal/repository/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ==========================================
// Project Repository Tests
// ==========================================

func TestProjectRepository_Create_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)

	// Arrange: Create owner user first
	ownerUser := &domain.User{
		ID:       uuid.New(),
		Name:     "Project Owner",
		Email:    "owner_" + uuid.New().String() + "@example.com",
		Password: "hashed_password",
	}
	err := userRepo.Create(ownerUser)
	assert.NoError(t, err)

	projectID := uuid.New()
	testProject := &domain.Project{
		ProjectID:   projectID,
		ProjectName: "Test Project",
		Members: []domain.ProjectMember{
			{
				UserID: ownerUser.ID,
				Role:   "Admin",
			},
		},
	}

	// Act
	err = projectRepo.Create(testProject)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, projectID, testProject.ProjectID)
	assert.NotZero(t, testProject.CreatedAt)

	// Cleanup
	testDB.Exec("DELETE FROM project_members WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM columns WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	testDB.Exec("DELETE FROM users WHERE id = ?", ownerUser.ID)
}

func TestProjectRepository_Create_WithDefaultColumns(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)

	// Arrange
	ownerUser := &domain.User{
		ID:       uuid.New(),
		Name:     "Project Owner",
		Email:    "owner_" + uuid.New().String() + "@example.com",
		Password: "hashed_password",
	}
	err := userRepo.Create(ownerUser)
	assert.NoError(t, err)

	projectID := uuid.New()
	testProject := &domain.Project{
		ProjectID:   projectID,
		ProjectName: "Test Project With Columns",
		Members: []domain.ProjectMember{
			{
				UserID: ownerUser.ID,
				Role:   "Admin",
			},
		},
	}

	// Act
	err = projectRepo.Create(testProject)
	assert.NoError(t, err)

	// Act: Fetch project to verify default columns
	fetchedProject, err := projectRepo.GetByID(projectID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, fetchedProject)
	assert.Equal(t, 3, len(fetchedProject.Columns))
	assert.Equal(t, "Todo", fetchedProject.Columns[0].Name)
	assert.Equal(t, "Doing", fetchedProject.Columns[1].Name)
	assert.Equal(t, "Done", fetchedProject.Columns[2].Name)

	// Cleanup
	testDB.Exec("DELETE FROM project_members WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM columns WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	testDB.Exec("DELETE FROM users WHERE id = ?", ownerUser.ID)
}

func TestProjectRepository_GetByID_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)

	// Arrange
	ownerUser := &domain.User{
		ID:       uuid.New(),
		Name:     "Owner",
		Email:    "owner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err := userRepo.Create(ownerUser)
	assert.NoError(t, err)

	projectID := uuid.New()
	testProject := &domain.Project{
		ProjectID:   projectID,
		ProjectName: "My Project",
		Members: []domain.ProjectMember{
			{UserID: ownerUser.ID, Role: "Admin"},
		},
	}
	err = projectRepo.Create(testProject)
	assert.NoError(t, err)

	// Act
	fetchedProject, err := projectRepo.GetByID(projectID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, fetchedProject)
	assert.Equal(t, projectID, fetchedProject.ProjectID)
	assert.Equal(t, "My Project", fetchedProject.ProjectName)
	assert.NotEmpty(t, fetchedProject.Members)
	assert.NotEmpty(t, fetchedProject.Columns)

	// Cleanup
	testDB.Exec("DELETE FROM project_members WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM columns WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	testDB.Exec("DELETE FROM users WHERE id = ?", ownerUser.ID)
}

func TestProjectRepository_GetByID_NotFound(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	projectRepo := postgresRepo.NewProjectRepository(testDB)

	// Act
	project, err := projectRepo.GetByID(uuid.New())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, project)
}

func TestProjectRepository_GetByUserID_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)

	// Arrange
	user := &domain.User{
		ID:       uuid.New(),
		Name:     "User",
		Email:    "user_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err := userRepo.Create(user)
	assert.NoError(t, err)

	// Create multiple projects for this user
	projectID1 := uuid.New()
	project1 := &domain.Project{
		ProjectID:   projectID1,
		ProjectName: "Project 1",
		Members: []domain.ProjectMember{
			{UserID: user.ID, Role: "Admin"},
		},
	}
	err = projectRepo.Create(project1)
	assert.NoError(t, err)

	projectID2 := uuid.New()
	project2 := &domain.Project{
		ProjectID:   projectID2,
		ProjectName: "Project 2",
		Members: []domain.ProjectMember{
			{UserID: user.ID, Role: "Member"},
		},
	}
	err = projectRepo.Create(project2)
	assert.NoError(t, err)

	// Act
	projects, err := projectRepo.GetByUserID(user.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, projects)
	assert.GreaterOrEqual(t, len(projects), 2)

	// Cleanup
	testDB.Exec("DELETE FROM project_members WHERE project_id IN (?, ?)", projectID1, projectID2)
	testDB.Exec("DELETE FROM columns WHERE project_id IN (?, ?)", projectID1, projectID2)
	testDB.Exec("DELETE FROM projects WHERE id IN (?, ?)", projectID1, projectID2)
	testDB.Exec("DELETE FROM users WHERE id = ?", user.ID)
}

func TestProjectRepository_GetByUserID_Empty(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	projectRepo := postgresRepo.NewProjectRepository(testDB)

	// Act
	projects, err := projectRepo.GetByUserID(uuid.New())

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, projects)
}

// ==========================================
// Column Management Tests
// ==========================================

func TestProjectRepository_CreateColumn_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)

	// Arrange
	owner := &domain.User{
		ID:       uuid.New(),
		Name:     "Owner",
		Email:    "owner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err := userRepo.Create(owner)
	assert.NoError(t, err)

	projectID := uuid.New()
	project := &domain.Project{
		ProjectID:   projectID,
		ProjectName: "Project",
		Members: []domain.ProjectMember{
			{UserID: owner.ID, Role: "Admin"},
		},
	}
	err = projectRepo.Create(project)
	assert.NoError(t, err)

	column := &domain.Column{
		ProjectID: projectID,
		Name:      "In Progress",
		Order:     4,
		Color:     "#FF6B6B",
	}

	// Act
	err = projectRepo.CreateColumn(column)

	// Assert
	assert.NoError(t, err)
	assert.NotZero(t, column.ColumnID)

	// Cleanup
	testDB.Exec("DELETE FROM columns WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM project_members WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	testDB.Exec("DELETE FROM users WHERE id = ?", owner.ID)
}

func TestProjectRepository_UpdateColumn_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)

	// Arrange
	owner := &domain.User{
		ID:       uuid.New(),
		Name:     "Owner",
		Email:    "owner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err := userRepo.Create(owner)
	assert.NoError(t, err)

	projectID := uuid.New()
	project := &domain.Project{
		ProjectID:   projectID,
		ProjectName: "Project",
		Members: []domain.ProjectMember{
			{UserID: owner.ID, Role: "Admin"},
		},
	}
	err = projectRepo.Create(project)
	assert.NoError(t, err)

	column := &domain.Column{
		ProjectID: projectID,
		Name:      "Original Name",
		Order:     4,
		Color:     "#FF6B6B",
	}
	err = projectRepo.CreateColumn(column)
	assert.NoError(t, err)

	// Act: Update column
	column.Name = "Updated Name"
	column.Color = "#4ECDC4"
	column.Order = 5
	err = projectRepo.UpdateColumn(column)

	// Assert
	assert.NoError(t, err)

	// Cleanup
	testDB.Exec("DELETE FROM columns WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM project_members WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	testDB.Exec("DELETE FROM users WHERE id = ?", owner.ID)
}

func TestProjectRepository_DeleteColumn_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)

	// Arrange
	owner := &domain.User{
		ID:       uuid.New(),
		Name:     "Owner",
		Email:    "owner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err := userRepo.Create(owner)
	assert.NoError(t, err)

	projectID := uuid.New()
	project := &domain.Project{
		ProjectID:   projectID,
		ProjectName: "Project",
		Members: []domain.ProjectMember{
			{UserID: owner.ID, Role: "Admin"},
		},
	}
	err = projectRepo.Create(project)
	assert.NoError(t, err)

	column := &domain.Column{
		ProjectID: projectID,
		Name:      "To Delete",
		Order:     4,
		Color:     "#FF6B6B",
	}
	err = projectRepo.CreateColumn(column)
	assert.NoError(t, err)

	columnID := column.ColumnID

	// Act
	err = projectRepo.DeleteColumn(columnID)

	// Assert
	assert.NoError(t, err)

	// Cleanup
	testDB.Exec("DELETE FROM columns WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM project_members WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	testDB.Exec("DELETE FROM users WHERE id = ?", owner.ID)
}

// ==========================================
// Member Management Tests
// ==========================================

func TestProjectRepository_AddMembers_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)

	// Arrange
	owner := &domain.User{
		ID:       uuid.New(),
		Name:     "Owner",
		Email:    "owner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err := userRepo.Create(owner)
	assert.NoError(t, err)

	member := &domain.User{
		ID:       uuid.New(),
		Name:     "Member",
		Email:    "member_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err = userRepo.Create(member)
	assert.NoError(t, err)

	projectID := uuid.New()
	project := &domain.Project{
		ProjectID:   projectID,
		ProjectName: "Project",
		Members: []domain.ProjectMember{
			{UserID: owner.ID, Role: "Admin"},
		},
	}
	err = projectRepo.Create(project)
	assert.NoError(t, err)

	// Act
	newMembers := []domain.ProjectMember{
		{
			ProjectID: projectID,
			UserID:    member.ID,
			Role:      "Member",
			Position:  "Developer",
		},
	}
	err = projectRepo.AddMembers(newMembers)

	// Assert
	assert.NoError(t, err)

	// Verify member was added
	fetchedMember, err := projectRepo.GetMember(projectID, member.ID)
	assert.NoError(t, err)
	assert.NotNil(t, fetchedMember)
	assert.Equal(t, member.ID, fetchedMember.UserID)

	// Cleanup
	testDB.Exec("DELETE FROM project_members WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM columns WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	testDB.Exec("DELETE FROM users WHERE id IN (?, ?)", owner.ID, member.ID)
}

func TestProjectRepository_GetMember_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)

	// Arrange
	owner := &domain.User{
		ID:       uuid.New(),
		Name:     "Owner",
		Email:    "owner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err := userRepo.Create(owner)
	assert.NoError(t, err)

	projectID := uuid.New()
	project := &domain.Project{
		ProjectID:   projectID,
		ProjectName: "Project",
		Members: []domain.ProjectMember{
			{UserID: owner.ID, Role: "Admin"},
		},
	}
	err = projectRepo.Create(project)
	assert.NoError(t, err)

	// Act
	member, err := projectRepo.GetMember(projectID, owner.ID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, member)
	assert.Equal(t, owner.ID, member.UserID)
	assert.Equal(t, projectID, member.ProjectID)

	// Cleanup
	testDB.Exec("DELETE FROM project_members WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM columns WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	testDB.Exec("DELETE FROM users WHERE id = ?", owner.ID)
}

func TestProjectRepository_UpdateMember_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)

	// Arrange
	owner := &domain.User{
		ID:       uuid.New(),
		Name:     "Owner",
		Email:    "owner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err := userRepo.Create(owner)
	assert.NoError(t, err)

	projectID := uuid.New()
	project := &domain.Project{
		ProjectID:   projectID,
		ProjectName: "Project",
		Members: []domain.ProjectMember{
			{UserID: owner.ID, Role: "Admin"},
		},
	}
	err = projectRepo.Create(project)
	assert.NoError(t, err)

	// Act: Update member
	updatedMember := &domain.ProjectMember{
		ProjectID: projectID,
		UserID:    owner.ID,
		Role:      "Owner",
		Position:  "Lead",
	}
	err = projectRepo.UpdateMember(updatedMember)

	// Assert
	assert.NoError(t, err)

	// Verify update
	fetchedMember, err := projectRepo.GetMember(projectID, owner.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Owner", fetchedMember.Role)
	assert.Equal(t, "Lead", fetchedMember.Position)

	// Cleanup
	testDB.Exec("DELETE FROM project_members WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM columns WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	testDB.Exec("DELETE FROM users WHERE id = ?", owner.ID)
}

func TestProjectRepository_RemoveMember_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)

	// Arrange
	owner := &domain.User{
		ID:       uuid.New(),
		Name:     "Owner",
		Email:    "owner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err := userRepo.Create(owner)
	assert.NoError(t, err)

	member := &domain.User{
		ID:       uuid.New(),
		Name:     "Member",
		Email:    "member_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err = userRepo.Create(member)
	assert.NoError(t, err)

	projectID := uuid.New()
	project := &domain.Project{
		ProjectID:   projectID,
		ProjectName: "Project",
		Members: []domain.ProjectMember{
			{UserID: owner.ID, Role: "Admin"},
		},
	}
	err = projectRepo.Create(project)
	assert.NoError(t, err)

	// Add a member
	err = projectRepo.AddMembers([]domain.ProjectMember{
		{ProjectID: projectID, UserID: member.ID, Role: "Member"},
	})
	assert.NoError(t, err)

	// Act: Remove the member
	err = projectRepo.RemoveMember(projectID, member.ID)

	// Assert
	assert.NoError(t, err)

	// Verify member is removed
	fetched, err := projectRepo.GetMember(projectID, member.ID)
	assert.Error(t, err)
	assert.Nil(t, fetched)

	// Cleanup
	testDB.Exec("DELETE FROM project_members WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM columns WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	testDB.Exec("DELETE FROM users WHERE id IN (?, ?)", owner.ID, member.ID)
}
