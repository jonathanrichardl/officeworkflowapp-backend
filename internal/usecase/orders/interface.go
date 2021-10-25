package orders

import (
	"clean/internal/entity"
	"time"
)

//Reader interface
type Reader interface {
	Get(id string) (*entity.Orders, error)
	Search(query string) ([]*entity.Orders, error)
	List() ([]*entity.Orders, error)
}

//Writer book writer
type Writer interface {
	Create(e *entity.Orders) (string, error)
	Update(e *entity.Orders) error
	Delete(id string) error
}

//Repository interface
type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	GetOrder(id string) (*entity.Orders, error)
	SearchOrders(query string) ([]*entity.Orders, error)
	ListOrders() ([]*entity.Orders, error)
	NewOrder(title string, description string, deadline time.Time) (string, error)
	UpdateOrder(o *entity.Orders) error
	DeleteOrder(id string) error
}
