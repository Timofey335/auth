package user

import (
	"context"
	"errors"
	"log"

	"github.com/fatih/color"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"

	"github.com/Timofey335/auth/internal/model"
)

// CreateUser - создает нового пользователя
func (s *serv) CreateUser(ctx context.Context, user *model.UserModel) (int64, error) {
	err := validation.Validate(user.Name, validation.Required, validation.Length(2, 50))
	if err != nil {
		log.Println(color.HiMagentaString("error while creating the new user: %v, with ctx: %v", err, ctx))

		return 0, err
	}

	err = validation.Validate(user.Email, validation.Required, is.Email)
	if err != nil {
		log.Println(color.HiMagentaString("error while creating the new user: %v, with ctx: %v", err, ctx))

		return 0, err
	}

	if user.Password != user.PasswordConfirm {
		err := errors.New("password doesn't match")
		log.Println(color.HiMagentaString("error while creating the new user: %v, with ctx: %v", err, ctx))

		return 0, err
	}

	err = validation.Validate(user.Password, validation.Required, validation.Length(8, 50))
	if err != nil {
		log.Println(color.HiMagentaString("error while creating the new user: %v, with ctx: %v", err, ctx))

		return 0, err
	}

	var id int64
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.userRepository.CreateUser(ctx, user)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	log.Println(color.BlueString("create user: %v, with ctx: %v", user, ctx))

	return id, nil
}
