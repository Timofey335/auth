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
	user := User{
		Name:             req.Name,
		Email:            req.Email,
		Password:         req.Password,
		Password_confirm: req.PasswordConfirm,
	}

	if err := user.userValidation(); err != nil {
		log.Println(color.HiMagentaString("Error while creating a new user '%v', email '%v'. %v", user.Name, user.Email, err))

		return nil, err
	} else {
		log.Println(color.BlueString("Create user: %v, with ctx: %v", req, ctx))

		return &desc.CreateUserResponse{
			Id: gofakeit.Int64(),
		}, nil
	}
}
