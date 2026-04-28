package domain

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	TaskID      uuid.UUID `json:"task_id"`
	ProjectID   uuid.UUID `json:"project_id"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Assigner *User   `json:"assigner,omitempty"`
	Members  []User  `json:"members"`
	Column   *Column `json:"column,omitempty"`
}

type TaskRepository interface {
	Create(task *Task) error
	Update(task *Task) error
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*Task, error)
	GetByProjectID(projectID uuid.UUID) ([]Task, error)

	// For manage relationship Many-to-Many of members.
	AddMember(taskID uuid.UUID, userID uuid.UUID) error
	RemoveMember(taskID uuid.UUID, userID uuid.UUID) error
}

type TaskUsecase interface {
	CreateTask(actorID uuid.UUID, task *Task) error
	UpdateTask(actorID uuid.UUID, task *Task) error
	DeleteTask(actorID uuid.UUID, taskID uuid.UUID) error

	// For update Task status.
	MoveTask(actorID uuid.UUID, taskID uuid.UUID, targetColumnID uuid.UUID) error

	AssignMember(actorID uuid.UUID, taskID uuid.UUID, targetUserID uuid.UUID)

	GetTasksByProject(actorID uuid.UUID, projectID uuid.UUID) ([]Task, error)
}
