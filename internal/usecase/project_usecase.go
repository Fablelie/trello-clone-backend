package usecase

import (
	"errors"

	"github.com/fablelie/trello-clone-backend/internal/domain"
	"github.com/google/uuid"
)

type projectUsecase struct {
	projectRepo domain.ProjectRepository
}

// NewProjectUsecase creates a new instance of project business logic
func NewProjectUsecase(repo domain.ProjectRepository) domain.ProjectUsecase {
	return &projectUsecase{
		projectRepo: repo,
	}
}

// CreateProject handles the creation of a project and initializes its owner
func (u *projectUsecase) CreateProject(name string, ownerID uuid.UUID) error {
	projectID := uuid.New()
	project := &domain.Project{
		ProjectID:   projectID,
		ProjectName: name,
		// Assign owner as the first member (Admin)
		Members: []domain.ProjectMember{
			{
				ProjectID: projectID,
				UserID:    ownerID,
				Role:      "Admin",
				Position:  "Owner",
			},
		},
	}

	return u.projectRepo.Create(project)
}

// AddMember checks if the actor is an Admin before adding a new member
func (u *projectUsecase) AddMember(actorID uuid.UUID, member *domain.ProjectMember) error {
	// Check the permissions of the actor
	actorMember, err := u.projectRepo.GetMember(member.ProjectID, actorID)
	if err != nil {
		return errors.New("permission denied: you are not a member of this project")
	}

	if actorMember.Role != "Admin" {
		return errors.New("permission denied: only Admin can add members")
	}

	// If Admin, proceed to add the member
	return u.projectRepo.AddMember(member)
}

// UpdateColumn checks if the actor is an Admin before updating a column
func (u *projectUsecase) UpdateColumn(actorID uuid.UUID, column *domain.Column) error {
	actorMember, err := u.projectRepo.GetMember(column.ProjectID, actorID)
	if err != nil {
		return errors.New("permission denied: you are not a member of this project")
	}

	if actorMember.Role != "Admin" {
		return errors.New("permission denied: only Admin can manage columns")
	}

	return u.projectRepo.UpdateColumn(column)
}

// RemoveMember checks if the actor is an Admin before removing a member
func (u *projectUsecase) RemoveMember(actorID uuid.UUID, projectID uuid.UUID, userID uuid.UUID) error {
	actorMember, err := u.projectRepo.GetMember(projectID, actorID)
	if err != nil {
		return errors.New("permission denied: you are not a member of this project")
	}

	if actorMember.Role != "Admin" {
		return errors.New("permission denied: only Admin can remove members")
	}

	return u.projectRepo.RemoveMember(projectID, userID)
}
