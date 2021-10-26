package entity

type Role uint8

const (
	Admin Role = iota
	Worker
)

type User struct {
	ID       string
	Username string
	Email    string
	Password string
	UserRole string
}

func NewUser(email string, username string, password string, Role string) *User {
	u := User{
		ID:       NewUUID().String(),
		Username: username,
		Email:    email,
		Password: password,
		UserRole: Role,
	}
	return &u
}
