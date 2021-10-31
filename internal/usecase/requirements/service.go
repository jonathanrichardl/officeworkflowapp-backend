package requirements

import (
	"errors"
	"strings"

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

func (s *Service) GetRequirementsbyOrderId(orderID string) ([]*entity.Requirements, error) {
	return s.repo.GetByOrderID(orderID)
}
func (s *Service) CreateRequirement(request string, expectedOutcome string, orderID string, userID *string) (int, error) {
	e := entity.NewRequirement(request, expectedOutcome, orderID)
	return s.repo.Create(e)
}

func (s *Service) GetRequirementbyID(id int) (*entity.Requirements, error) {
	return s.repo.Get(id)
}

func (s *Service) SearchRequirements(query string) ([]*entity.Requirements, error) {
	return s.repo.Search(strings.ToLower(query))
}

func (s *Service) ListRequirements() ([]*entity.Requirements, error) {
	return s.repo.List()
}

func (s *Service) DeleteRequirement(id int) error {
	u, err := s.GetRequirementbyID(id)
	if u == nil {
		return errors.New("not found")
	}
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *Service) UpdateRequirement(e *entity.Requirements) error {
	return s.repo.Update(e)
}
