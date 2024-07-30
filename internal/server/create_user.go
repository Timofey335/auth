package server

import (
	"context"
	"log"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"

	desc "github.com/Timofey335/auth/pkg/authservice_v1"
)

// CreateUser - create a new user
func (s *Server) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	log.Println(color.BlueString("Create user: %v, with ctx: %v", req, ctx))

	return &desc.CreateUserResponse{
		Id: gofakeit.Int64(),
	}, nil
}
