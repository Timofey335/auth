package model

import (
	"database/sql"
	"time"
)

// User - модель User
type User struct {
	ID              int64
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            int64
	CreatedAt       time.Time
	UpdatedAt       sql.NullTime
}
