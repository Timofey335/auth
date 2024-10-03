package user

import (
	"context"
	"errors"
	"log"

	"github.com/fatih/color"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/auth/internal/model"
)

// UpdateUser - обновляет данные пользователя
func (s *serv) UpdateUser(ctx context.Context, user *model.UserUpdateModel) (*emptypb.Empty, error) {
	if *user.Password != *user.PasswordConfirm {
		err := errors.New("password doesn't match")
		log.Println(color.HiMagentaString("error while updating the new user: %v, with ctx: %v", err, ctx))

		return nil, err
	}

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error

		_, errTx = s.userRepository.UpdateUser(ctx, user)
		if errTx != nil {
			return errTx
		}

		errTx = s.cache.DeleteUser(ctx, user.ID)
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	log.Println(color.BlueString("updated the user %v, with ctx: %v", user.ID, ctx))

	return &emptypb.Empty{}, nil
}
