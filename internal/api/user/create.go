package user

import (
	"context"
	"log"

	"github.com/Timofey335/auth/internal/converter"
	"github.com/fatih/color"

	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

func (s *Implementation) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	id, err := s.userService.CreateUser(ctx, converter.ToUserFromDesc(req))
	if err != nil {
		return nil, err
	}

	log.Println(color.BlueString("create user: %v, with ctx: %v", req, ctx))

	return &desc.CreateUserResponse{
		Id: id,
	}, nil
}
