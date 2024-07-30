package server

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/Timofey335/auth/pkg/authservice_v1"
)

// UpdateUser - update information of the user by id
func (s *Server) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
