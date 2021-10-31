package tasks

import "order-validation-v2/internal/entity"

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) ListAllTasks() ([]*entity.Task, error) {
	return s.repo.List()
}

func (s *Service) GetTasksofUser(userID string) ([]*entity.Task, error) {
	return s.repo.Get(userID)

}

func (s *Service) UpdateTask(t *entity.Task) error {
	return s.repo.Update(t)
}

func (s *Service) DeleteTask(id string) error {
	return s.repo.Delete(id)
}

func (s *Service) CreateTask(requirementID int, userID string) (string, error) {
	task := entity.NewTask(requirementID, userID)
	return s.repo.Create(task)
}
