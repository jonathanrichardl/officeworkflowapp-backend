package user

import (
	"database/sql"
	"order-validation-v2/internal/entity"
)

//Reader interface
type Reader interface {
	GetbyID(ID string) (*entity.User, error)
	GetbyUsername(username string) (*entity.User, error)
	Search(query string) ([]*entity.User, error)
	List() ([]*entity.User, error)
	CustomQuery(query string) (*sql.Rows, error)
}

//Writer user writer
type Writer interface {
	Create(r *entity.User) (string, error)
	Update(r *entity.User) error
	Delete(ID string) error
}

//Repository interface
type Repository interface {
	Reader
	Writer
}

type UseCase interface {
	GetUserbyUsername(username string) (*entity.User, error)
	SearchUser(query string) ([]*entity.User, error)
	ListUsers() ([]*entity.User, error)
	CreateUser(username string, email string, password string, role string) (string, error)
	UpdateUser(u *entity.User) error
	DeleteUser(username string) error
	Login(username string, password string) (string, string, bool, error)
}
