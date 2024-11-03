package user

import (
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/auth/internal/converter"
	"github.com/Timofey335/auth/internal/logger"
	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

// UpdateUser - обновляет данные пользователя
func (i *Implementation) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	_, err := i.userService.UpdateUser(ctx, converter.ToUserFromDescUpd(req))
	if err != nil {
		return nil, err
	}

	logger.Info("Updated user", zap.Any("Id", req.Id))

	return &emptypb.Empty{}, nil
}
