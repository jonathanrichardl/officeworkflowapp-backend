package orders

import (
	"clean/internal/entity"
	"time"
)

//Reader interface
type Reader interface {
	Get(id entity.ID) (*entity.Orders, error)
	Search(query string) ([]*entity.Orders, error)
	List() ([]*entity.Orders, error)
}

//Writer book writer
type Writer interface {
	Create(e *entity.Orders) (entity.ID, error)
	Update(e *entity.Orders) error
	Delete(id entity.ID) error
}

//Repository interface
type Repository interface {
	Reader
	Writer
}

//UseCase interface
type UseCase interface {
	GetOrder(id entity.ID) (*entity.Orders, error)
	SearchOrders(query string) ([]*entity.Orders, error)
	ListOrders() ([]*entity.Orders, error)
	NewOrder(title string, description string, deadline time.Time) (entity.ID, error)
	UpdateOrder(o *entity.Orders) error
	DeleteOrder(id entity.ID) error
}
