package user

import (
	"context"
	"errors"
	"log"

	"github.com/fatih/color"
	"golang.org/x/crypto/bcrypt"

	"github.com/Timofey335/auth/internal/model"
)

// CreateUser - создает нового пользователя
func (s *serv) CreateUser(ctx context.Context, user *model.UserModel) (int64, error) {
	if user.Password != user.PasswordConfirm {
		err := errors.New("password doesn't match")
		log.Println(color.HiMagentaString("error while creating the new user: %v, with ctx: %v", err, ctx))

		return 0, err
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errors.New("error while hashing of the password")
	}

	user.Password = string(passHash)

	var id int64
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.userRepository.CreateUser(ctx, user)
		if errTx != nil {
			return errTx
		}

		errTx = s.cache.CreateUser(ctx, &model.UserModel{
			ID:        id,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
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
