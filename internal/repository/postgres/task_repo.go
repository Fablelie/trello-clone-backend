package postgres

import (
	"github.com/fablelie/trello-clone-backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type taskRepo struct {
	db *gorm.DB
}

// NewTaskRepository creates a new instance for task data management
func NewTaskRepository(db *gorm.DB) domain.TaskRepository {
	return &taskRepo{db: db}
}

// Create persists a new task to the database
func (r *taskRepo) Create(task *domain.Task) error {
	// Mapping Domain to Schema (Using ID from TaskSchema as TaskID)
	schema := TaskSchema{
		ID:          task.TaskID,
		ProjectID:   task.ProjectID,
		ColumnID:    task.Column.ColumnID,
		AssignerID:  task.Assigner.ID,
		Subject:     task.Subject,
		Description: task.Description,
	}

	if err := r.db.Create(&schema).Error; err != nil {
		return err
	}

	task.TaskID = schema.ID
	task.CreatedAt = schema.CreatedAt
	return nil
}

// Update handles editing task content
func (r *taskRepo) Update(task *domain.Task) error {
	updates := make(map[string]interface{})

	if task.Subject != "" {
		updates["subject"] = task.Subject
	}

	if task.Description != "" {
		updates["description"] = task.Description
	}

	if task.Column != nil && task.Column.ColumnID != uuid.Nil {
		updates["column_id"] = task.Column.ColumnID
	}

	if len(updates) == 0 {
		return nil
	}

	updates["updated_at"] = gorm.Expr("NOW()")

	return r.db.Model(&TaskSchema{}).Where("id = ?", task.TaskID).Updates(updates).Error
}

// Delete removes a task by ID
func (r *taskRepo) Delete(id uuid.UUID) error {
	return r.db.Delete(&TaskSchema{}, "id = ?", id).Error
}

// GetByID checks if task exists and returns basic info
func (r *taskRepo) GetByID(id uuid.UUID) (*domain.Task, error) {
	var s TaskSchema

	// Use Preload to fetch Assigner data to verify if they are the task creator
	if err := r.db.Preload("Assigner").Preload("Column").Preload("Members").First(&s, "id = ?", id).Error; err != nil {
		return nil, err
	}

	task := &domain.Task{
		TaskID:      s.ID,
		ProjectID:   s.ProjectID,
		Subject:     s.Subject,
		Description: s.Description,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
		Assigner: &domain.User{
			ID:   s.Assigner.ID,
			Name: s.Assigner.Name,
		},
		Column: &domain.Column{
			ColumnID: s.Column.ID,
			Name:     s.Column.Name,
			Color:    s.Column.Color,
		},
	}
	// Map Many-to-Many Members
	for _, m := range s.Members {
		task.Members = append(task.Members, domain.User{
			ID:   m.ID,
			Name: m.Name,
		})
	}
	// Mapping data from TaskSchema back to Domain.Task
	return task, nil
}

// GetByProjectID fetches all tasks in a project with Assigner, Column, and Members
func (r *taskRepo) GetByProjectID(projectID uuid.UUID) ([]domain.Task, error) {
	var schemas []TaskSchema
	err := r.db.Preload("Assigner").Preload("Column").Preload("Members").Where("project_id = ?", projectID).Find(&schemas).Error
	if err != nil {
		return nil, err
	}

	var tasks []domain.Task
	for _, s := range schemas {
		task := domain.Task{
			TaskID:      s.ID,
			ProjectID:   s.ProjectID,
			Subject:     s.Subject,
			Description: s.Description,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
			Assigner: &domain.User{
				ID:    s.Assigner.ID,
				Name:  s.Assigner.Name,
				Email: s.Assigner.Email,
			},
			Column: &domain.Column{
				ColumnID: s.Column.ID,
				Name:     s.Column.Name,
				Color:    s.Column.Color,
			},
		}

		// Map Many-to-Many Members
		for _, m := range s.Members {
			task.Members = append(task.Members, domain.User{
				ID:   m.ID,
				Name: m.Name,
			})
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *taskRepo) IsMember(taskID uuid.UUID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Table("task_members").
		Where("task_schema_id = ? AND user_schema_id = ?", taskID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// AddMembers adds multiple users to a task (Many-to-Many)
func (r *taskRepo) AddMembers(taskID uuid.UUID, userIDs []uuid.UUID) error {
	for _, userID := range userIDs {
		if err := r.db.Exec(
			"INSERT INTO task_members (task_schema_id, user_schema_id) VALUES ($1, $2)",
			taskID, userID,
		).Error; err != nil {
			return err
		}
	}
	return nil
}

// RemoveMember removes a user from a task
func (r *taskRepo) RemoveMember(taskID uuid.UUID, userID uuid.UUID) error {
	return r.db.Exec(
		"DELETE FROM task_members WHERE task_schema_id = $1 AND user_schema_id = $2",
		taskID, userID,
	).Error
}
