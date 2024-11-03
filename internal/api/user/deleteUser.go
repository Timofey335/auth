package user

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/auth/internal/logger"
	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

// DeleteUser - удаляет пользователя
func (i *Implementation) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	_, err := i.userService.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	logger.Info("Deleted user", zap.Any("Id", req.Id))

	return &emptypb.Empty{}, nil
}
