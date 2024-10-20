package access

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	descAccess "github.com/Timofey335/auth/pkg/access_v1"
)

// Check - проверяет разрешение на доступ к эндпоинту
func (a *AccessImplementation) Check(ctx context.Context, req *descAccess.CheckRequest) (*emptypb.Empty, error) {
	_, err := a.accessService.Check(ctx, req.EndpointAddress)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
