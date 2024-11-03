package user_saver

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/IBM/sarama"
	"go.uber.org/zap"

	"github.com/Timofey335/auth/internal/logger"
	"github.com/Timofey335/auth/internal/model"
)

// UserSaveHandler - сохраняет нового пользователя
func (s *service) UserSaveHandler(ctx context.Context, msg *sarama.ConsumerMessage) error {
	user := &model.UserModel{}
	err := json.Unmarshal(msg.Value, user)
	if err != nil {
		return err
	}

	if user.Password != user.PasswordConfirm {
		err := errors.New("password doesn't match")

		return err
	}

	id, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	logger.Info("Created user", zap.Any("Id", id), zap.String("Name", user.Name), zap.String("Email", user.Email))

	return nil
}
