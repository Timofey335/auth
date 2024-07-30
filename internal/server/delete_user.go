package server

import (
	"context"
	"log"

	"github.com/fatih/color"
	"google.golang.org/protobuf/types/known/emptypb"

	desc "github.com/Timofey335/auth/pkg/authservice_v1"
)

// DeleteUser - delete a user by id
func (s *Server) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	log.Println(color.HiMagentaString("Delete user: id %d, with ctx: %v", req.Id, ctx))

	return &emptypb.Empty{}, nil
}
