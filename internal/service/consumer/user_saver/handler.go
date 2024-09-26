package user_saver

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/IBM/sarama"
	"github.com/fatih/color"

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
		log.Println(color.HiMagentaString("error while creating the new user: %v, with ctx: %v", err, ctx))

		return err
	}

	id, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	log.Println(color.BlueString("create user: id-%d %v, with ctx: %v", id, user, ctx))

	return nil
}
