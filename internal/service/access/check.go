package access

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (*serv) Check(ctx context.Context, s string) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
