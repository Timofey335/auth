package user

import (
	"context"

	"go.uber.org/zap"

	"github.com/Timofey335/auth/internal/converter"
	"github.com/Timofey335/auth/internal/logger"
	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

// CreateUser - создает нового пользователя
func (i *Implementation) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	id, err := i.userService.CreateUser(ctx, converter.ToUserFromDesc(req))
	if err != nil {
		return nil, err
	}

	logger.Info("Created user", zap.Any("Id", id), zap.String("Name", req.Name), zap.String("Email", req.Email))

	return &desc.CreateUserResponse{
		Id: id,
	}, nil
}
