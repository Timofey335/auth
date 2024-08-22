package user

import (
	"context"
	"log"

	"github.com/fatih/color"

	"github.com/Timofey335/auth/internal/converter"
	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

func (i *Implementation) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	id, err := i.userService.CreateUser(ctx, converter.ToUserFromDesc(req))
	if err != nil {
		return nil, err
	}

	log.Println(color.BlueString("create user: %v, with ctx: %v", req, ctx))

	return &desc.CreateUserResponse{
		Id: id,
	}, nil
}
