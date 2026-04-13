package domain

import (
	"time"

	"github.com/google/uuid"
)

// Project is main entity.
type Project struct {
	ProjectID   uuid.UUID       `json:"project_id"`
	ProjectName string          `json:"project_name"`
	Members     []ProjectMember `json:"members"`
	Columns     []Column        `json:"columns"`
	CreatedAt   time.Time       `json:"created_at"`
}

// ProjectMember is middle table to stroage members in project.
type ProjectMember struct {
	ProjectID uuid.UUID `json:"project_id"`
	UserID    uuid.UUID `josn:"user_id"`
	Role      string    `json:"role"`
	Position  string    `json:"position"`
}

// Column is status of Task in Project like (Todo, Doing, Done)
type Column struct {
	ColumnID  uuid.UUID `json:"column_id"`
	ProjectID uuid.UUID `json:"project_id"`
	Name      string    `json:"name"`
	Order     int       `json:"order"`
	Color     string    `json:"color"`
}

type ProjectRepository interface {
	Create(project *Project) error
	GetByID(id uuid.UUID) (*Project, error)
	GetByUserID(userID uuid.UUID) ([]Project, error)

	// For manage Column in Project.
	CreateColumn(column *Column) error
	UpdateColumn(column *Column) error
	DeleteColumn(id uuid.UUID) error

	// For manage Members.
	AddMember(member *ProjectMember) error
	GetMember(projectID uuid.UUID, userID uuid.UUID) (*ProjectMember, error)
	UpdateMember(member *ProjectMember) error
	RemoveMember(projectID, userID uuid.UUID) error
}

type ProjectUsecase interface {
	CreateProject(name string, ownerID uuid.UUID) error

	// Requie actorID to check role.
	UpdateColumn(actorID uuid.UUID, column *Column) error
	AddMember(actorID uuid.UUID, member *ProjectMember) error
	RemoveMember(actorID uuid.UUID, projectID uuid.UUID, userID uuid.UUID) error
}
