package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/auth/internal/model"
)

// UpdateUser - обновляет данные пользователя
func (s *serv) UpdateUser(ctx context.Context, user *model.User) (*emptypb.Empty, error) {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		_, errTx = s.userRepository.UpdateUser(ctx, user)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.userRepository.GetUser(ctx, user.ID)
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
