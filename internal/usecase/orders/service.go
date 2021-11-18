package orders

import (
	"errors"
	"strings"
	"time"

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

func (s *Service) NewOrder(title string, description string, deadline time.Time) (string, error) {
	o := entity.NewOrder(title, description, deadline)
	return s.repo.Create(o)
}

func (s *Service) GetOrder(id string) (*entity.Orders, error) {
	o, err := s.repo.Get(id)

	if err != nil {
		return nil, err
	}

	if o == nil {
		return nil, errors.New("not found")
	}

	return o, nil
}

func (s *Service) SearchOrders(query string) ([]*entity.Orders, error) {
	orders, err := s.repo.Search(strings.ToLower(query))
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, errors.New("not found")
	}
	return orders, nil
}

func (s *Service) ListOrders() ([]*entity.Orders, error) {
	orders, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, errors.New("not found")
	}
	return orders, nil
}

func (s *Service) DeleteOrder(id string) error {
	_, err := s.GetOrder(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *Service) UpdateOrder(o *entity.Orders) error {
	return s.repo.Update(o)
}
