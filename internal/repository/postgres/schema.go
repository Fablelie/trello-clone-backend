package postgres

import (
	"time"

	"github.com/google/uuid"
)

// UserSchema represents the "users" table
type UserSchema struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time
}

func (UserSchema) TableName() string { return "users" }

// ProjectSchema represents the "projects" table
type ProjectSchema struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name      string    `gorm:"not null"`
	OwnerID   uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt time.Time

	// Relationships
	Columns []ColumnSchema        `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	Members []ProjectMemberSchema `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
}

func (ProjectSchema) TableName() string { return "projects" }

// ColumnSchema represents the "columns" table (formerly known as "statuses")
type ColumnSchema struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID uuid.UUID `gorm:"type:uuid;not null;index"`
	Name      string    `gorm:"not null"`
	Position  int       `gorm:"not null"`
	Color     string    `gorm:"size:7"` // ex. #ff0000
}

func (ColumnSchema) TableName() string { return "columns" }

// ProjectMemberSchema is an association table between Users and Projects with extra fields
type ProjectMemberSchema struct {
	ProjectID uuid.UUID `gorm:"type:uuid;primaryKey;constraint:OnDelete:CASCADE;references:id"`
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey;constraint:OnDelete:CASCADE;references:id"`
	Role      string    `gorm:"not null"`
	Position  string    `gorm:"not null"`
}

func (ProjectMemberSchema) TableName() string { return "project_members" }

// TaskSchema represents the "tasks" table
type TaskSchema struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID  uuid.UUID `gorm:"type:uuid;not null;index"`
	ColumnID   uuid.UUID `gorm:"type:uuid;not null"`
	AssignerID uuid.UUID `gorm:"type:uuid;not null"`

	Assigner UserSchema   `gorm:"foreignKey:AssignerID"`
	Column   ColumnSchema `gorm:"foreignKey:ColumnID"`
	// Many-to-many relationship using a junction table "task_members"
	Members []UserSchema `gorm:"many2many:task_members;constraint:OnDelete:CASCADE"`

	Subject     string `gorm:"not null"`
	Description string `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (TaskSchema) TableName() string { return "tasks" }
