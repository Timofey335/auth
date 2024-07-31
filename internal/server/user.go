package server

import (
	"fmt"
	"time"

	"github.com/go-ozzo/ozzo-validation/is"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	desc "github.com/Timofey335/auth/pkg/authservice_v1"
)

type User struct {
	ID               int64
	Name             string
	Email            string
	Password         string
	Password_confirm string
	Role             desc.Role
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (u User) userValidation() error {
	if u.Password != u.Password_confirm {
		return fmt.Errorf("password doesn't match")
	} else {
		return validation.ValidateStruct(&u,
			validation.Field(&u.Name, validation.Required, validation.Length(2, 50)),
			validation.Field(&u.Email, validation.Required, is.Email),
			validation.Field(&u.Password, validation.Required, validation.Length(8, 50)),
		)
	}
}
