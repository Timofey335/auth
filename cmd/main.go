package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	desc "github.com/Timofey335/auth/pkg/authservice_v1"
	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const grpcPort = 50051

type server struct {
	desc.UnimplementedAuthserviceV1Server
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Println(color.BlueString("Create user: %v, with ctx: %v", req, ctx))
	return &desc.CreateResponse{
		Id: gofakeit.Int64(),
	}, nil
}

func getUserRole() desc.Role {
	roles := []desc.Role{desc.Role_UNKNOW, desc.Role_USER, desc.Role_ADMIN}
	rand.NewSource(time.Now().UnixNano())
	randomRole := rand.Intn(len(roles))
	return roles[randomRole]
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	log.Println(color.BlueString("Note id: %d", req.GetId()))
	return &desc.GetResponse{
		Id:        gofakeit.Int64(),
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      getUserRole(),
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	log.Println(color.HiMagentaString("Delete user: id %d, with ctx: %v", req.Id, ctx))
	return &emptypb.Empty{}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
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
