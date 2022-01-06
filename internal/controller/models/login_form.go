package models

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

type ChangePasswordForm struct {
	UserID      string `json:"userid"`
	OldPassword string `json:"password"`
	NewPassword string `json:"new_password"`
}
