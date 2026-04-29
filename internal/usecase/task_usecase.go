package usecase

import (
	"errors"

	"github.com/fablelie/trello-clone-backend/internal/domain"
	"github.com/google/uuid"
)

type taskUsecase struct {
	taskRepo    domain.TaskRepository
	projectRepo domain.ProjectRepository
	userRepo    domain.UserRepository
}

func NewTaskUsecase(taskRepo domain.TaskRepository, projectRepo domain.ProjectRepository, userRepo domain.UserRepository) domain.TaskUsecase {
	return &taskUsecase{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}

// checkPermission is a helper to verify if the actor is an Admin or the Task creator
func (u *taskUsecase) checkPermission(actorID, projectID uuid.UUID, assignerID uuid.UUID) (bool, error) {
	// Check task assigner
	if actorID == assignerID {
		return true, nil
	}

	member, err := u.projectRepo.GetMember(projectID, actorID)
	if err != nil {
		return false, errors.New("permission denied: you are not a member of this project")
	}

	// Or role is Admin
	if member.Role == "Admin" {
		return true, nil
	}

	return false, errors.New("permission denied: only Admin or Task creator can perform this action")
}

func (u *taskUsecase) CreateTask(actorID uuid.UUID, task *domain.Task) error {
	// Only project member can create task
	_, err := u.projectRepo.GetMember(task.ProjectID, actorID)
	if err != nil {
		return errors.New("permission denied: you must be a member to create a task")
	}

	if task.Column == nil {
		return errors.New("column_id is required")
	}

	task.TaskID = uuid.New()
	task.Assigner = &domain.User{ID: actorID}
	return u.taskRepo.Create(task)
}

func (u *taskUsecase) UpdateTask(actorID uuid.UUID, task *domain.Task) error {
	// Get existing task to check permission from project id and assigner id
	existingTask, err := u.taskRepo.GetByID(task.TaskID)
	if err != nil {
		return errors.New("task not found")
	}

	if _, err := u.checkPermission(actorID, existingTask.ProjectID, existingTask.Assigner.ID); err != nil {
		return err
	}

	return u.taskRepo.Update(task)
}

// DeleteTask this task from project
func (u *taskUsecase) DeleteTask(actorID uuid.UUID, taskID uuid.UUID) error {
	existingTask, err := u.taskRepo.GetByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}

	if _, err := u.checkPermission(actorID, existingTask.ProjectID, existingTask.Assigner.ID); err != nil {
		return err
	}

	return u.taskRepo.Delete(taskID)
}

// MoveTask task to another column(update task status)
func (u *taskUsecase) MoveTask(actorID uuid.UUID, taskID uuid.UUID, targetColumnID uuid.UUID) error {
	existingTask, err := u.taskRepo.GetByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}

	if _, err := u.checkPermission(actorID, existingTask.ProjectID, existingTask.Assigner.ID); err != nil {
		return err
	}

	task := &domain.Task{
		TaskID: taskID,
		Column: &domain.Column{ColumnID: targetColumnID},
	}
	return u.taskRepo.Update(task)
}

func (u *taskUsecase) CheckUserInTask(taskID uuid.UUID, userID uuid.UUID) (bool, error) {
	exists, err := u.taskRepo.IsMember(taskID, userID)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// AssignMembers to this task.
func (u *taskUsecase) AssignMembers(actorID uuid.UUID, taskID uuid.UUID, emails []string) error {
	existingTask, err := u.taskRepo.GetByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}

	if _, err := u.checkPermission(actorID, existingTask.ProjectID, existingTask.Assigner.ID); err != nil {
		return err
	}

	var userIDs []uuid.UUID
	for _, email := range emails {
		user, err := u.userRepo.GetByEmail(email)
		if err != nil {
			return errors.New("user not found: " + email)
		}

		_, err = u.projectRepo.GetMember(existingTask.ProjectID, user.ID)
		if err != nil {
			return errors.New(email + " are not a member of this project")
		}

		userIDs = append(userIDs, user.ID)
	}

	return u.taskRepo.AddMembers(taskID, userIDs)
}

// RemoveMember from this task.
func (u *taskUsecase) RemoveMember(actorID uuid.UUID, taskID uuid.UUID, email string) error {
	existingTask, err := u.taskRepo.GetByID(taskID)
	if err != nil {
		return errors.New("task not found")
	}

	if _, err := u.checkPermission(actorID, existingTask.ProjectID, existingTask.Assigner.ID); err != nil {
		return err
	}

	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return errors.New("user not found: " + email)
	}

	b, err := u.CheckUserInTask(taskID, user.ID)
	if err != nil || !b {
		return errors.New(email + " is not associated with this task")
	}

	return u.taskRepo.RemoveMember(taskID, user.ID)
}

// GetTasksByProject get all tasks in project (only project member)
func (u *taskUsecase) GetTasksByProject(actorID uuid.UUID, projectID uuid.UUID) ([]domain.Task, error) {
	// Check permission
	_, err := u.projectRepo.GetMember(projectID, actorID)
	if err != nil {
		return nil, errors.New("permission denied: you are not a member of this project")
	}

	// get all tasks from repo
	return u.taskRepo.GetByProjectID(projectID)
}
