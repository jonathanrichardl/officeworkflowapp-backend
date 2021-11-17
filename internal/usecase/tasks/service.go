package tasks

import (
	"order-validation-v2/internal/entity"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) ListAllTasks() ([]*entity.TaskWithDetails, error) {
	return s.repo.List()
}

func (s *Service) GetTasksofUser(userID string) ([]*entity.TaskWithDetails, error) {
	return s.repo.GetbyUserID(userID)

}

func (s *Service) Get(id string) (*entity.Task, error) {
	return s.repo.Get(id)
}

func (s *Service) UpdateTask(t *entity.Task) error {
	return s.repo.Update(t)
}

func (s *Service) DeleteTask(id string) error {
	return s.repo.Delete(id)
}

func (s *Service) CreateTask(requirementID int, userID string, Note string, prerequisiteTaskID []string, Deadline time.Time) (string, error) {
	task := entity.NewTask(requirementID, userID, Note, prerequisiteTaskID, Deadline)
	return s.repo.Create(task)
}

func (s *Service) SaveTask(task *entity.Task) (string, error) {
	return s.repo.Create(task)
}

func (s *Service) RemovePrerequisite(prerequisiteID string) ([]*entity.Task, error) {
	return s.repo.RemovePrerequisite(prerequisiteID)
}
