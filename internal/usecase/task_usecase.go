package usecase

import (
	"errors"

	"github.com/fablelie/trello-clone-backend/internal/domain"
	"github.com/google/uuid"
)

type taskUsecase struct {
	taskRepo    domain.TaskRepository
	projectRepo domain.ProjectRepository
}

func NewTaskUsecase(taskRepo domain.TaskRepository, projectRepo domain.ProjectRepository) domain.TaskUsecase {
	return &taskUsecase{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
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

func (u *taskUsecase) AssignMember(actorID uuid.UUID, taskID uuid.UUID, targetUserID uuid.UUID) {
	existingTask, err := u.taskRepo.GetByID(taskID)
	if err != nil {
		return
	}

	if _, err := u.checkPermission(actorID, existingTask.ProjectID, existingTask.Assigner.ID); err == nil {
		u.taskRepo.AddMember(taskID, targetUserID)
	}
}
