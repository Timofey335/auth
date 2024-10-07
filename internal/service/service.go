package service

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/auth/internal/model"
)

// UserService - интерфейс сервисного слоя
type UserService interface {
	CreateUser(ctx context.Context, user *model.UserModel) (int64, error)
	DeleteUser(ctx context.Context, userId int64) (*emptypb.Empty, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
	GetRefreshToken(ctx context.Context, token string) (string, error)
	GetUser(ctx context.Context, userId int64) (*model.UserModel, error)
	Login(ctx context.Context, userLoginData *model.UserLoginModel) (string, error)
	UpdateUser(ctx context.Context, user *model.UserUpdateModel) (*emptypb.Empty, error)
}

// ConsumerService - интерфейс для consumer
type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}
