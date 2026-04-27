package usecase

import (
	"errors"

	"github.com/fablelie/trello-clone-backend/internal/domain"
	"github.com/google/uuid"
)

type projectUsecase struct {
	projectRepo domain.ProjectRepository
	userRepo    domain.UserRepository
}

// NewProjectUsecase creates a new instance of project business logic
func NewProjectUsecase(repo domain.ProjectRepository, userRepo domain.UserRepository) domain.ProjectUsecase {
	return &projectUsecase{
		projectRepo: repo,
		userRepo:    userRepo,
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
func (u *projectUsecase) AddMembers(actorID uuid.UUID, projectID uuid.UUID, members []domain.AddMemberRequest) error {
	// Check the permissions of the actor
	actorMember, err := u.projectRepo.GetMember(projectID, actorID)
	if err != nil || actorMember == nil {
		return errors.New("permission denied: you are not a member of this project")
	}

	if actorMember.Role != "Admin" {
		return errors.New("permission denied: only Admin can add members")
	}

	var newMembers []domain.ProjectMember

	for _, member := range members {
		user, err := u.userRepo.GetByEmail(member.Email)
		if err != nil {
			continue
		}

		newMembers = append(newMembers, domain.ProjectMember{
			ProjectID: projectID,
			UserID:    user.ID,
			Role:      member.Role,
			Position:  member.Position,
		})
	}

	if len(newMembers) == 0 {
		return errors.New("newMember length == 0")
	}

	// If Admin, proceed to add the member
	return u.projectRepo.AddMembers(newMembers)
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

// GetMyProjects get all projects the user is a member of
func (u *projectUsecase) GetMyProjects(actorID uuid.UUID) ([]domain.Project, error) {
	return u.projectRepo.GetByUserID(actorID)
}

// GetProjectByID get project board details(members only)
func (u *projectUsecase) GetProjectByID(actorID uuid.UUID, projectID uuid.UUID) (*domain.Project, error) {
	member, err := u.projectRepo.GetMember(projectID, actorID)
	if err != nil || member == nil {
		return nil, errors.New("permission denied: you are not a member of this project")
	}

	return u.projectRepo.GetByID(projectID)
}
