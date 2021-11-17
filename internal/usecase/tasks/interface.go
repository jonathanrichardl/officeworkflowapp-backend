package tasks

import (
	"order-validation-v2/internal/entity"
	"time"
)

type Reader interface {
	Get(id string) (*entity.Task, error)
	GetbyUserID(userID string) ([]*entity.TaskWithDetails, error)
	List() ([]*entity.TaskWithDetails, error)
}

type Writer interface {
	Create(t *entity.Task) (string, error)
	Update(t *entity.Task) error
	Delete(id string) error
	RemovePrerequisite(prerequisiteID string) ([]*entity.Task, error)
}

type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	Get(id string) (*entity.Task, error)
	ListAllTasks() ([]*entity.TaskWithDetails, error)
	GetTasksofUser(userID string) ([]*entity.TaskWithDetails, error)
	UpdateTask(t *entity.Task) error
	DeleteTask(id string) error
	CreateTask(requirementID int, userID string, Note string, prerequisiteTaskID []string, Deadline time.Time) (string, error)
	RemovePrerequisite(prerequisiteTaskID string) ([]*entity.Task, error)
	SaveTask(t *entity.Task) (string, error)
}
