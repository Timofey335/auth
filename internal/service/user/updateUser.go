package user

import (
	"context"
	"errors"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/auth/internal/model"
)

// UpdateUser - обновляет данные пользователя
func (s *serv) UpdateUser(ctx context.Context, user *model.UserUpdateModel) (*emptypb.Empty, error) {
	if *user.Password != *user.PasswordConfirm {
		err := errors.New("password doesn't match")

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

	return &emptypb.Empty{}, nil
}
