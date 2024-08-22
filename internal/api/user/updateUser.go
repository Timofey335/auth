package user

import (
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/auth/internal/converter"
	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

func (i *Implementation) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	_, err := i.userService.UpdateUser(ctx, converter.ToUserFromDescUpd(req))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
