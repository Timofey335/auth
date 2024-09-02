package user

import (
	"context"

	"github.com/Timofey335/auth/internal/model"
)

// CreateUser - создает нового пользователя
func (s *serv) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	var id int64
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.userRepository.CreateUser(ctx, user)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.userRepository.GetUser(ctx, id)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
