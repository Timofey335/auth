package model

import (
	"time"
)

// User - модель User
type UserModel struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
	PasswordConfirm string    `json:"password_confirm"`
	Role            int64     `json:"role"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// UserUpdateModel - модель для метода Update
type UserUpdateModel struct {
	ID              int64
	Name            *string
	Password        *string
	PasswordConfirm *string
	Role            *int64
}

// UserData ...
type UserData struct {
	Username string `json:"username"`
	Role     int64  `json:"role"`
}

type UserLoginModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
