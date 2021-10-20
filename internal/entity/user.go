package entity

type User struct {
	ID       string
	Username string
	Email    string
	Password string
}

func NewUser(email string, username string, password string) *User {
	u := User{
		ID:       NewUUID().String(),
		Username: username,
		Email:    email,
		Password: password,
	}
	return &u
}
