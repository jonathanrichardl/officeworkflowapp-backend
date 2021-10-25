package user

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"order-validation-v2/internal/entity"
	"strings"
)

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) GetUserbyID(id string) (*entity.User, error) {
	u, err := s.repo.GetbyID(id)

	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, errors.New("not found")
	}

	return u, nil
}

func (s *Service) CreateUser(username string, email string, password string) (string, error) {
	u := entity.NewUser(email, username, password)
	return s.repo.Create(u)

}

func (s *Service) GetUserbyUsername(username string) (*entity.User, error) {
	u, err := s.repo.GetbyUsername(username)

	if err != nil {
		return nil, err
	}

	if u == nil {
		return nil, errors.New("not found")
	}

	return u, nil
}

func (s *Service) SearchUser(query string) ([]*entity.User, error) {
	users, err := s.repo.Search(strings.ToLower(query))
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("not found")
	}
	return users, nil
}

func (s *Service) ListUsers() ([]*entity.User, error) {
	users, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.New("not found")
	}
	return users, nil
}

func (s *Service) DeleteUser(id string) error {
	_, err := s.GetUserbyID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *Service) UpdateUser(u *entity.User) error {
	return s.repo.Update(u)
}

func (s *Service) Login(username string, password string) (string, bool, error) {
	u, err := s.repo.GetbyUsername(username)
	if err != nil {
		return username, false, err
	}

	if u == nil {
		return username, false, errors.New("Username/Password wrong")
	}
	sum := sha256.Sum256([]byte(password))
	incomingPassword := fmt.Sprintf("\\x+%x", sum)

	fmt.Println(incomingPassword)
	fmt.Println(u.Password)
	if u.Password == incomingPassword {
		return u.ID, true, nil

	}
	return username, false, errors.New("Username/Password wrong")

}
