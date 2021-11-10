package models

import "order-validation-v2/internal/entity"

type NewUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type RetrievedUser struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func BuildUserProfile(user *entity.User) RetrievedUser {
	return RetrievedUser{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.UserRole,
	}

}
