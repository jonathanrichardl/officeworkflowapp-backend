package submissions

import (
	"order-validation-v2/internal/entity"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}
func (s *Service) NewSubmission(message string, images []string, TaskID string) (string, error) {
	submission := entity.NewSubmission(message, images, TaskID)
	return s.repo.Create(submission)
}
func (s *Service) EditSubmission(e *entity.Submission) error {
	return s.repo.Update(e)
}

func (s *Service) DeleteSubmission(id string) error {
	return s.repo.Delete(id)

}

func (s *Service) GetSubmission(id string) (*entity.Submission, error) {
	return s.repo.Get(id)
}
