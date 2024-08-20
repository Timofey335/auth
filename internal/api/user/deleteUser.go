package user

import (
	"context"

	desc "github.com/Timofey335/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	_, err := i.userService.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
