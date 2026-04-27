package postgres

import (
	"github.com/fablelie/trello-clone-backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type projectRepo struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) domain.ProjectRepository {
	return &projectRepo{db: db}
}

// Create a new project with default columns and an admin member
func (r *projectRepo) Create(project *domain.Project) error {
	schema := ProjectSchema{
		ID:      project.ProjectID,
		Name:    project.ProjectName,
		OwnerID: project.Members[0].UserID, // Get owner from the first member
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&schema).Error; err != nil {
			return err
		}

		// Create default columns: Todo, Doing, Done
		defaultColumns := []ColumnSchema{
			{ProjectID: schema.ID, Name: "Todo", Position: 1, Color: "#D1D5DB"},
			{ProjectID: schema.ID, Name: "Doing", Position: 2, Color: "#3B82F6"},
			{ProjectID: schema.ID, Name: "Done", Position: 3, Color: "#10B981"},
		}

		if err := tx.Create(&defaultColumns).Error; err != nil {
			return err
		}

		// Add owner as Admin member
		ownerMember := ProjectMemberSchema{
			ProjectID: schema.ID,
			UserID:    schema.OwnerID,
			Role:      "Admin",
			Position:  "Owner",
		}

		if err := tx.Create(&ownerMember).Error; err != nil {
			return err
		}

		project.ProjectID = schema.ID
		project.CreatedAt = schema.CreatedAt
		return nil
	})
}

// Get project details by ID including columns and members
func (r *projectRepo) GetByID(id uuid.UUID) (*domain.Project, error) {
	var schema ProjectSchema
	if err := r.db.Preload("Columns").Preload("Members").First(&schema, "id = ?", id).Error; err != nil {
		return nil, err
	}

	project := &domain.Project{
		ProjectID:   schema.ID,
		ProjectName: schema.Name,
		CreatedAt:   schema.CreatedAt,
	}

	for _, col := range schema.Columns {
		project.Columns = append(project.Columns, domain.Column{
			ColumnID:  col.ID,
			ProjectID: col.ProjectID,
			Name:      col.Name,
			Order:     col.Position,
			Color:     col.Color,
		})
	}

	for _, mem := range schema.Members {
		project.Members = append(project.Members, domain.ProjectMember{
			ProjectID: mem.ProjectID,
			UserID:    mem.UserID,
			Role:      mem.Role,
			Position:  mem.Position,
		})
	}

	return project, nil
}

// Get all projects for a specific user via membership
func (r *projectRepo) GetByUserID(userID uuid.UUID) ([]domain.Project, error) {
	var memberSchemas []ProjectMemberSchema
	if err := r.db.Where("user_id = ?", userID).Find(&memberSchemas).Error; err != nil {
		return nil, err
	}

	var projects []domain.Project
	for _, ms := range memberSchemas {
		p, err := r.GetByID(ms.ProjectID)
		if err == nil {
			projects = append(projects, *p)
		}
	}
	return projects, nil
}

// --- Column Management ---

func (r *projectRepo) CreateColumn(column *domain.Column) error {
	schema := ColumnSchema{
		ProjectID: column.ProjectID,
		Name:      column.Name,
		Position:  column.Order,
		Color:     column.Color,
	}
	if err := r.db.Create(&schema).Error; err != nil {
		return err
	}
	column.ColumnID = schema.ID
	return nil
}

func (r *projectRepo) UpdateColumn(column *domain.Column) error {
	return r.db.Model(&ColumnSchema{}).Where("id = ?", column.ColumnID).Updates(map[string]interface{}{
		"name":     column.Name,
		"position": column.Order,
		"color":    column.Color,
	}).Error
}

func (r *projectRepo) DeleteColumn(id uuid.UUID) error {
	return r.db.Delete(&ColumnSchema{}, "id = ?", id).Error
}

// --- Member Management ---

func (r *projectRepo) AddMembers(members []domain.ProjectMember) error {
	var schemas []ProjectMemberSchema

	for _, m := range members {
		schemas = append(schemas, ProjectMemberSchema{
			ProjectID: m.ProjectID,
			UserID:    m.UserID,
			Role:      m.Role,
			Position:  m.Position,
		})
	}
	return r.db.Create(&schemas).Error
}

func (r *projectRepo) GetMember(projectID uuid.UUID, userID uuid.UUID) (*domain.ProjectMember, error) {
	var schema ProjectMemberSchema
	if err := r.db.Where("project_id = ? AND user_id = ?", projectID, userID).First(&schema).Error; err != nil {
		return nil, err
	}
	return &domain.ProjectMember{
		ProjectID: schema.ProjectID,
		UserID:    schema.UserID,
		Role:      schema.Role,
		Position:  schema.Position,
	}, nil
}

func (r *projectRepo) UpdateMember(member *domain.ProjectMember) error {
	return r.db.Model(&ProjectMemberSchema{}).
		Where("project_id = ? AND user_id = ?", member.ProjectID, member.UserID).
		Updates(map[string]interface{}{
			"role":     member.Role,
			"position": member.Position,
		}).Error
}

func (r *projectRepo) RemoveMember(projectID, userID uuid.UUID) error {
	return r.db.Where("project_id = ? AND user_id = ?", projectID, userID).Delete(&ProjectMemberSchema{}).Error
}
