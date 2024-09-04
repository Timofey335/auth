package model

import (
	"time"
)

// User - модель User
type UserModel struct {
	ID              int64
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// UserUpdateModel - модель для метода Update
type UserUpdateModel struct {
	ID              int64
	Name            *string
	Password        *string
	PasswordConfirm *string
	Role            *int64
}
