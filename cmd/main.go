package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/Timofey335/auth/pkg/authservice_v1"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedAuthserviceV1Server
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf(color.RedString("failed listen: %v", err))
	}

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthserviceV1Server(s, &server{})
	log.Println(color.BlueString("server listening at %v", lis.Addr()))
	if err := s.Serve(lis); err == nil {
		log.Fatalf(color.RedString("failed to serve: %v", err))
	}
}

// CreateUser - create a new user
func (s *server) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	log.Println(color.BlueString("Create user: %v, with ctx: %v", req, ctx))

	return &desc.CreateUserResponse{
		Id: gofakeit.Int64(),
	}, nil
}

func getUserRole() desc.Role {
	roles := []desc.Role{desc.Role_UNKNOWN, desc.Role_USER, desc.Role_ADMIN}
	rand.NewSource(time.Now().UnixNano())
	randomRole := rand.Intn(len(roles))

	return roles[randomRole]
}

// GetUser - get information of the user by id
func (s *server) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
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

// DeleteUser - delete a user by id
func (s *server) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	log.Println(color.HiMagentaString("Delete user: id %d, with ctx: %v", req.Id, ctx))

	return &emptypb.Empty{}, nil
}

// UpdateUser - update information of the user by id
func (s *server) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
