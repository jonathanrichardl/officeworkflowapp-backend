package requirements

import (
	"clean/internal/entity"
	"database/sql"
)

//Reader interface
type Reader interface {
	Get(id int) (*entity.Requirements, error)
	Search(query string) ([]*entity.Requirements, error)
	List() ([]*entity.Requirements, error)
	CustomQuery(query string) (*sql.Rows, error)
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
	GetRequirementsbyOrderId(id string) ([]*entity.Requirements, error)
	SearchRequirements(query string) ([]*entity.Requirements, error)
	ListRequirements() ([]*entity.Requirements, error)
	CreateRequirement(request string, expectedOutcome string, orderID entity.ID) (int, error)
	UpdateRequirement(e *entity.Requirements) error
	DeleteRequirement(id entity.ID) error
}
