package tasks

import (
	"order-validation-v2/internal/entity"
)

type Reader interface {
	Get(userID string) ([]*entity.Task, error)
	List() ([]*entity.Task, error)
}

type Writer interface {
	Create(r *entity.Task) (string, error)
	Update(r *entity.Task) error
	Delete(id string) error
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	ListAllTasks() ([]*entity.Task, error)
	GetTasksofUser(userID string) ([]*entity.Task, error)
	UpdateTask(t *entity.Task) error
	DeleteTask(id string) error
	CreateTask(requirementID int, userID string) (string, error)
}
