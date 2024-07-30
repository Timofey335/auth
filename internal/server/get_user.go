package server

import (
	"context"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"golang.org/x/exp/rand"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/Timofey335/auth/pkg/authservice_v1"
)

func getUserRole() desc.Role {
	roles := []desc.Role{desc.Role_UNKNOWN, desc.Role_USER, desc.Role_ADMIN}
	rand.NewSource(uint64(time.Now().UnixNano()))
	randomRole := rand.Intn(len(roles))

	return roles[randomRole]
}

// GetUser - get information of the user by id
func (s *Server) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	log.Println(color.BlueString("Note id: %d", req.GetId()))

	return &desc.GetUserResponse{
		Id:        gofakeit.Int64(),
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      getUserRole(),
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}
