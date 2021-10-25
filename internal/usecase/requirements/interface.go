package requirements

import (
	"order-validation-v2/internal/entity"
)

//Reader interface
type Reader interface {
	Get(id int) (*entity.Requirements, error)
	Search(query string) ([]*entity.Requirements, error)
	List() ([]*entity.Requirements, error)
	GetByOrderID(orderID string) ([]*entity.Requirements, error)
	GetByUserID(userID string) ([]*entity.Requirements, error)
}

//Writer user writer
type Writer interface {
	Create(r *entity.Requirements) (int, error)
	Update(r *entity.Requirements) error
	Delete(id int) error
}

//Repository interface
type Repository interface {
	Reader
	Writer
}

//UseCase interface
type UseCase interface {
	GetRequirementbyID(id int) (*entity.Requirements, error)
	GetRequirementsbyOrderId(orderID string) ([]*entity.Requirements, error)
	GetRequirementsbyUserId(userID string) ([]*entity.Requirements, error)
	SearchRequirements(query string) ([]*entity.Requirements, error)
	ListRequirements() ([]*entity.Requirements, error)
	CreateRequirement(request string, expectedOutcome string, orderID string, userID *string) (int, error)
	UpdateRequirement(e *entity.Requirements) error
	DeleteRequirement(id int) error
}
