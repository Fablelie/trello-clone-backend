package repository

import (
	"github.com/fablelie/trello-clone-backend/internal/domain"
	my_postgres "github.com/fablelie/trello-clone-backend/internal/repository/postgres"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of the user repository
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *domain.User) error {
	// Mapping: Convert Domain entity to DB Schema
	schema := my_postgres.UserSchema{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}

	// Persist to database
	if err := r.db.Create(&schema).Error; err != nil {
		return err
	}

	// Update Domain entity with DB-generated fields (ID, CreatedAt)
	user.ID = schema.ID
	user.CreatedAt = schema.CreatedAt
	return nil
}

func (r *userRepo) GetByEmail(email string) (*domain.User, error) {
	var schema my_postgres.UserSchema
	if err := r.db.Where("email = ?", email).First(&schema).Error; err != nil {
		return nil, err
	}

	// Mapping: Convert Schema back to Domain entity
	return &domain.User{
		ID:        schema.ID,
		Name:      schema.Name,
		Email:     schema.Email,
		Password:  schema.Password,
		CreatedAt: schema.CreatedAt,
	}, nil
}

func (r *userRepo) GetByID(id uuid.UUID) (*domain.User, error) {
	var schema my_postgres.UserSchema
	if err := r.db.First(&schema, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &domain.User{
		ID:        schema.ID,
		Name:      schema.Name,
		Email:     schema.Email,
		Password:  schema.Password,
		CreatedAt: schema.CreatedAt,
	}, nil
}
