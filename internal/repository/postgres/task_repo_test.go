package postgres_test

import (
	"testing"

	"github.com/fablelie/trello-clone-backend/internal/domain"
	postgresRepo "github.com/fablelie/trello-clone-backend/internal/repository/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// ==========================================
// Task Repository Tests
// ==========================================

func TestTaskRepository_Create_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)
	taskRepo := postgresRepo.NewTaskRepository(testDB)

	// Arrange: Create owner, project, and column
	owner := &domain.User{
		ID:       uuid.New(),
		Name:     "Owner",
		Email:    "owner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err := userRepo.Create(owner)
	assert.NoError(t, err)

	assigner := &domain.User{
		ID:       uuid.New(),
		Name:     "Assigner",
		Email:    "assigner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err = userRepo.Create(assigner)
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

	// Get the first column (Todo)
	fetchedProject, err := projectRepo.GetByID(projectID)
	assert.NoError(t, err)
	columnID := fetchedProject.Columns[0].ColumnID

	// Create task
	taskID := uuid.New()
	task := &domain.Task{
		TaskID:      taskID,
		ProjectID:   projectID,
		Subject:     "New Task",
		Description: "Task Description",
		Assigner:    &domain.User{ID: assigner.ID},
		Column:      &domain.Column{ColumnID: columnID},
	}

	// Act
	err = taskRepo.Create(task)

	// Assert
	assert.NoError(t, err)
	assert.NotZero(t, task.CreatedAt)

	// Cleanup
	testDB.Exec("DELETE FROM task_members WHERE task_schema_id = ?", taskID)
	testDB.Exec("DELETE FROM tasks WHERE id = ?", taskID)
	testDB.Exec("DELETE FROM project_members WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM columns WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	testDB.Exec("DELETE FROM users WHERE id IN (?, ?)", owner.ID, assigner.ID)
}

func TestTaskRepository_GetByID_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)
	taskRepo := postgresRepo.NewTaskRepository(testDB)

	// Arrange: Setup
	owner := &domain.User{
		ID:       uuid.New(),
		Name:     "Owner",
		Email:    "owner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err := userRepo.Create(owner)
	assert.NoError(t, err)

	assigner := &domain.User{
		ID:       uuid.New(),
		Name:     "Assigner",
		Email:    "assigner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err = userRepo.Create(assigner)
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

	fetchedProject, err := projectRepo.GetByID(projectID)
	assert.NoError(t, err)
	columnID := fetchedProject.Columns[0].ColumnID

	taskID := uuid.New()
	task := &domain.Task{
		TaskID:      taskID,
		ProjectID:   projectID,
		Subject:     "Task to Fetch",
		Description: "Description",
		Assigner:    &domain.User{ID: assigner.ID},
		Column:      &domain.Column{ColumnID: columnID},
	}
	err = taskRepo.Create(task)
	assert.NoError(t, err)

	// Act
	fetchedTask, err := taskRepo.GetByID(taskID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, fetchedTask)
	assert.Equal(t, taskID, fetchedTask.TaskID)
	assert.Equal(t, "Task to Fetch", fetchedTask.Subject)
	assert.Equal(t, "Description", fetchedTask.Description)
	assert.NotNil(t, fetchedTask.Assigner)
	assert.Equal(t, assigner.ID, fetchedTask.Assigner.ID)

	// Cleanup
	testDB.Exec("DELETE FROM task_members WHERE task_schema_id = ?", taskID)
	testDB.Exec("DELETE FROM tasks WHERE id = ?", taskID)
	testDB.Exec("DELETE FROM project_members WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM columns WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	testDB.Exec("DELETE FROM users WHERE id IN (?, ?)", owner.ID, assigner.ID)
}

func TestTaskRepository_GetByID_NotFound(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	taskRepo := postgresRepo.NewTaskRepository(testDB)

	// Act
	task, err := taskRepo.GetByID(uuid.New())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, task)
}

func TestTaskRepository_Update_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)
	taskRepo := postgresRepo.NewTaskRepository(testDB)

	// Arrange: Setup
	owner := &domain.User{
		ID:       uuid.New(),
		Name:     "Owner",
		Email:    "owner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err := userRepo.Create(owner)
	assert.NoError(t, err)

	assigner := &domain.User{
		ID:       uuid.New(),
		Name:     "Assigner",
		Email:    "assigner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err = userRepo.Create(assigner)
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

	fetchedProject, err := projectRepo.GetByID(projectID)
	assert.NoError(t, err)
	columnID1 := fetchedProject.Columns[0].ColumnID
	columnID2 := fetchedProject.Columns[1].ColumnID

	taskID := uuid.New()
	task := &domain.Task{
		TaskID:      taskID,
		ProjectID:   projectID,
		Subject:     "Original Subject",
		Description: "Original Description",
		Assigner:    &domain.User{ID: assigner.ID},
		Column:      &domain.Column{ColumnID: columnID1},
	}
	err = taskRepo.Create(task)
	assert.NoError(t, err)

	// Act: Update task
	task.Subject = "Updated Subject"
	task.Description = "Updated Description"
	task.Column.ColumnID = columnID2
	err = taskRepo.Update(task)

	// Assert
	assert.NoError(t, err)

	// Verify update
	updated, err := taskRepo.GetByID(taskID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Subject", updated.Subject)
	assert.Equal(t, "Updated Description", updated.Description)

	// Cleanup
	testDB.Exec("DELETE FROM task_members WHERE task_schema_id = ?", taskID)
	testDB.Exec("DELETE FROM tasks WHERE id = ?", taskID)
	testDB.Exec("DELETE FROM project_members WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM columns WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	testDB.Exec("DELETE FROM users WHERE id IN (?, ?)", owner.ID, assigner.ID)
}

func TestTaskRepository_Delete_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)
	taskRepo := postgresRepo.NewTaskRepository(testDB)

	// Arrange: Setup
	owner := &domain.User{
		ID:       uuid.New(),
		Name:     "Owner",
		Email:    "owner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err := userRepo.Create(owner)
	assert.NoError(t, err)

	assigner := &domain.User{
		ID:       uuid.New(),
		Name:     "Assigner",
		Email:    "assigner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err = userRepo.Create(assigner)
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

	fetchedProject, err := projectRepo.GetByID(projectID)
	assert.NoError(t, err)
	columnID := fetchedProject.Columns[0].ColumnID

	taskID := uuid.New()
	task := &domain.Task{
		TaskID:      taskID,
		ProjectID:   projectID,
		Subject:     "Task to Delete",
		Description: "Description",
		Assigner:    &domain.User{ID: assigner.ID},
		Column:      &domain.Column{ColumnID: columnID},
	}
	err = taskRepo.Create(task)
	assert.NoError(t, err)

	// Act
	err = taskRepo.Delete(taskID)

	// Assert
	assert.NoError(t, err)

	// Verify deletion
	deleted, err := taskRepo.GetByID(taskID)
	assert.Error(t, err)
	assert.Nil(t, deleted)

	// Cleanup
	testDB.Exec("DELETE FROM task_members WHERE task_schema_id = ?", taskID)
	testDB.Exec("DELETE FROM project_members WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM columns WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	testDB.Exec("DELETE FROM users WHERE id IN (?, ?)", owner.ID, assigner.ID)
}

func TestTaskRepository_GetByProjectID_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)
	taskRepo := postgresRepo.NewTaskRepository(testDB)

	// Arrange: Setup
	owner := &domain.User{
		ID:       uuid.New(),
		Name:     "Owner",
		Email:    "owner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err := userRepo.Create(owner)
	assert.NoError(t, err)

	assigner := &domain.User{
		ID:       uuid.New(),
		Name:     "Assigner",
		Email:    "assigner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err = userRepo.Create(assigner)
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

	fetchedProject, err := projectRepo.GetByID(projectID)
	assert.NoError(t, err)
	columnID := fetchedProject.Columns[0].ColumnID

	// Create multiple tasks
	taskID1 := uuid.New()
	task1 := &domain.Task{
		TaskID:      taskID1,
		ProjectID:   projectID,
		Subject:     "Task 1",
		Description: "Description 1",
		Assigner:    &domain.User{ID: assigner.ID},
		Column:      &domain.Column{ColumnID: columnID},
	}
	err = taskRepo.Create(task1)
	assert.NoError(t, err)

	taskID2 := uuid.New()
	task2 := &domain.Task{
		TaskID:      taskID2,
		ProjectID:   projectID,
		Subject:     "Task 2",
		Description: "Description 2",
		Assigner:    &domain.User{ID: assigner.ID},
		Column:      &domain.Column{ColumnID: columnID},
	}
	err = taskRepo.Create(task2)
	assert.NoError(t, err)

	// Act
	tasks, err := taskRepo.GetByProjectID(projectID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tasks)
	assert.GreaterOrEqual(t, len(tasks), 2)

	// Cleanup
	testDB.Exec("DELETE FROM task_members WHERE task_schema_id IN (?, ?)", taskID1, taskID2)
	testDB.Exec("DELETE FROM tasks WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM project_members WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM columns WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	testDB.Exec("DELETE FROM users WHERE id IN (?, ?)", owner.ID, assigner.ID)
}

func TestTaskRepository_GetByProjectID_Empty(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	taskRepo := postgresRepo.NewTaskRepository(testDB)

	// Act
	tasks, err := taskRepo.GetByProjectID(uuid.New())

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, tasks)
}

// ==========================================
// Task Members Management Tests
// ==========================================

func TestTaskRepository_AddMembers_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)
	taskRepo := postgresRepo.NewTaskRepository(testDB)

	// Arrange: Setup
	owner := &domain.User{
		ID:       uuid.New(),
		Name:     "Owner",
		Email:    "owner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err := userRepo.Create(owner)
	assert.NoError(t, err)

	assigner := &domain.User{
		ID:       uuid.New(),
		Name:     "Assigner",
		Email:    "assigner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err = userRepo.Create(assigner)
	assert.NoError(t, err)

	member1 := &domain.User{
		ID:       uuid.New(),
		Name:     "Member 1",
		Email:    "member1_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err = userRepo.Create(member1)
	assert.NoError(t, err)

	member2 := &domain.User{
		ID:       uuid.New(),
		Name:     "Member 2",
		Email:    "member2_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err = userRepo.Create(member2)
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

	fetchedProject, err := projectRepo.GetByID(projectID)
	assert.NoError(t, err)
	columnID := fetchedProject.Columns[0].ColumnID

	taskID := uuid.New()
	task := &domain.Task{
		TaskID:      taskID,
		ProjectID:   projectID,
		Subject:     "Task",
		Description: "Description",
		Assigner:    &domain.User{ID: assigner.ID},
		Column:      &domain.Column{ColumnID: columnID},
	}
	err = taskRepo.Create(task)
	assert.NoError(t, err)

	// Act
	err = taskRepo.AddMembers(taskID, []uuid.UUID{member1.ID, member2.ID})

	// Assert
	assert.NoError(t, err)

	// Verify members were added
	isMember1, err := taskRepo.IsMember(taskID, member1.ID)
	assert.NoError(t, err)
	assert.True(t, isMember1)

	isMember2, err := taskRepo.IsMember(taskID, member2.ID)
	assert.NoError(t, err)
	assert.True(t, isMember2)

	// Cleanup
	testDB.Exec("DELETE FROM task_members WHERE task_schema_id = ?", taskID)
	testDB.Exec("DELETE FROM tasks WHERE id = ?", taskID)
	testDB.Exec("DELETE FROM project_members WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM columns WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	testDB.Exec("DELETE FROM users WHERE id IN (?, ?, ?, ?)", owner.ID, assigner.ID, member1.ID, member2.ID)
}

func TestTaskRepository_IsMember_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)
	taskRepo := postgresRepo.NewTaskRepository(testDB)

	// Arrange: Setup
	owner := &domain.User{
		ID:       uuid.New(),
		Name:     "Owner",
		Email:    "owner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err := userRepo.Create(owner)
	assert.NoError(t, err)

	assigner := &domain.User{
		ID:       uuid.New(),
		Name:     "Assigner",
		Email:    "assigner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err = userRepo.Create(assigner)
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

	fetchedProject, err := projectRepo.GetByID(projectID)
	assert.NoError(t, err)
	columnID := fetchedProject.Columns[0].ColumnID

	taskID := uuid.New()
	task := &domain.Task{
		TaskID:      taskID,
		ProjectID:   projectID,
		Subject:     "Task",
		Description: "Description",
		Assigner:    &domain.User{ID: assigner.ID},
		Column:      &domain.Column{ColumnID: columnID},
	}
	err = taskRepo.Create(task)
	assert.NoError(t, err)

	// Add member to task
	err = taskRepo.AddMembers(taskID, []uuid.UUID{member.ID})
	assert.NoError(t, err)

	// Act
	isMember, err := taskRepo.IsMember(taskID, member.ID)

	// Assert
	assert.NoError(t, err)
	assert.True(t, isMember)

	// Cleanup
	testDB.Exec("DELETE FROM task_members WHERE task_schema_id = ?", taskID)
	testDB.Exec("DELETE FROM tasks WHERE id = ?", taskID)
	testDB.Exec("DELETE FROM project_members WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM columns WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	testDB.Exec("DELETE FROM users WHERE id IN (?, ?, ?)", owner.ID, assigner.ID, member.ID)
}

func TestTaskRepository_IsMember_NotMember(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)
	taskRepo := postgresRepo.NewTaskRepository(testDB)

	// Arrange: Setup
	owner := &domain.User{
		ID:       uuid.New(),
		Name:     "Owner",
		Email:    "owner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err := userRepo.Create(owner)
	assert.NoError(t, err)

	assigner := &domain.User{
		ID:       uuid.New(),
		Name:     "Assigner",
		Email:    "assigner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err = userRepo.Create(assigner)
	assert.NoError(t, err)

	nonMember := &domain.User{
		ID:       uuid.New(),
		Name:     "Non Member",
		Email:    "nonmember_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err = userRepo.Create(nonMember)
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

	fetchedProject, err := projectRepo.GetByID(projectID)
	assert.NoError(t, err)
	columnID := fetchedProject.Columns[0].ColumnID

	taskID := uuid.New()
	task := &domain.Task{
		TaskID:      taskID,
		ProjectID:   projectID,
		Subject:     "Task",
		Description: "Description",
		Assigner:    &domain.User{ID: assigner.ID},
		Column:      &domain.Column{ColumnID: columnID},
	}
	err = taskRepo.Create(task)
	assert.NoError(t, err)

	// Act
	isMember, err := taskRepo.IsMember(taskID, nonMember.ID)

	// Assert
	assert.NoError(t, err)
	assert.False(t, isMember)

	// Cleanup
	testDB.Exec("DELETE FROM task_members WHERE task_schema_id = ?", taskID)
	testDB.Exec("DELETE FROM tasks WHERE id = ?", taskID)
	testDB.Exec("DELETE FROM project_members WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM columns WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	testDB.Exec("DELETE FROM users WHERE id IN (?, ?, ?)", owner.ID, assigner.ID, nonMember.ID)
}

func TestTaskRepository_RemoveMember_Success(t *testing.T) {
	if testDB == nil {
		t.Skip("Skipping integration test: DB connection not established")
	}

	userRepo := postgresRepo.NewUserRepository(testDB)
	projectRepo := postgresRepo.NewProjectRepository(testDB)
	taskRepo := postgresRepo.NewTaskRepository(testDB)

	// Arrange: Setup
	owner := &domain.User{
		ID:       uuid.New(),
		Name:     "Owner",
		Email:    "owner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err := userRepo.Create(owner)
	assert.NoError(t, err)

	assigner := &domain.User{
		ID:       uuid.New(),
		Name:     "Assigner",
		Email:    "assigner_" + uuid.New().String() + "@example.com",
		Password: "password",
	}
	err = userRepo.Create(assigner)
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

	fetchedProject, err := projectRepo.GetByID(projectID)
	assert.NoError(t, err)
	columnID := fetchedProject.Columns[0].ColumnID

	taskID := uuid.New()
	task := &domain.Task{
		TaskID:      taskID,
		ProjectID:   projectID,
		Subject:     "Task",
		Description: "Description",
		Assigner:    &domain.User{ID: assigner.ID},
		Column:      &domain.Column{ColumnID: columnID},
	}
	err = taskRepo.Create(task)
	assert.NoError(t, err)

	// Add member to task
	err = taskRepo.AddMembers(taskID, []uuid.UUID{member.ID})
	assert.NoError(t, err)

	// Verify member is added
	isMember, err := taskRepo.IsMember(taskID, member.ID)
	assert.NoError(t, err)
	assert.True(t, isMember)

	// Act
	err = taskRepo.RemoveMember(taskID, member.ID)

	// Assert
	assert.NoError(t, err)

	// Verify member is removed
	isMember, err = taskRepo.IsMember(taskID, member.ID)
	assert.NoError(t, err)
	assert.False(t, isMember)

	// Cleanup
	testDB.Exec("DELETE FROM task_members WHERE task_schema_id = ?", taskID)
	testDB.Exec("DELETE FROM tasks WHERE id = ?", taskID)
	testDB.Exec("DELETE FROM project_members WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM columns WHERE project_id = ?", projectID)
	testDB.Exec("DELETE FROM projects WHERE id = ?", projectID)
	testDB.Exec("DELETE FROM users WHERE id IN (?, ?, ?)", owner.ID, assigner.ID, member.ID)
}
