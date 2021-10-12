package requirements

import (
	"errors"
	"fmt"
	"strings"

	"clean/internal/entity"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) GetRequirementsbyOrderId(OrderID string) ([]*entity.Requirements, error) {
	rows, err := s.repo.CustomQuery(fmt.Sprintf(
		"SELECT id, request,outcome, status FROM requirements WHERE order_id = %s ",
		OrderID))
	if err != nil {
		return nil, err
	}
	var requirements []*entity.Requirements
	for rows.Next() {
		var r entity.Requirements
		err = rows.Scan(&r.Id, &r.Request, &r.ExpectedOutcome, &r.Status)
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, &r)
	}
	return requirements, nil

}
func (s *Service) CreateRequirement(request string, expectedOutcome string, orderID entity.ID) (int, error) {
	e := entity.NewRequirement(request, expectedOutcome, orderID)
	return s.repo.Create(e)
}

func (s *Service) GetRequirementbyID(id int) (*entity.Requirements, error) {
	return s.repo.Get(id)
}

func (s *Service) SearchRequirements(query string) ([]*entity.Requirements, error) {
	return s.repo.Search(strings.ToLower(query))
}

//ListUsers List users
func (s *Service) ListRequirements() ([]*entity.Requirements, error) {
	return s.repo.List()
}

//DeleteUser Delete an user
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

//UpdateUser Update an user
func (s *Service) UpdateRequirement(e *entity.Requirements) error {
	return s.repo.Update(e)
}
