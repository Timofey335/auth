package user

import (
	"context"

	"go.uber.org/zap"

	"github.com/Timofey335/auth/internal/converter"
	"github.com/Timofey335/auth/internal/logger"
	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

// GetUser - получает данные о пользователе
func (i *Implementation) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	userObj, err := i.userService.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	logger.Info("Get user", zap.Any("Id", req.Id), zap.String("Name", userObj.Name), zap.String("Email", userObj.Email))

	userObjConvert := converter.ToUserFromService(userObj)

	return &desc.GetUserResponse{
		Id:        userObjConvert.Id,
		Name:      userObj.Name,
		Email:     userObj.Email,
		Role:      userObjConvert.Role,
		CreatedAt: userObjConvert.CreatedAt,
		UpdatedAt: userObjConvert.UpdatedAt,
	}, nil
}
