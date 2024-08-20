package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/fatih/color"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	userApi "github.com/Timofey335/auth/internal/api/user"
	user "github.com/Timofey335/auth/internal/repository/user"
	userService "github.com/Timofey335/auth/internal/service/user"
	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

const (
	grpcPort = 50051
	dbDSN    = "host=localhost port=54321 dbname=users user=user password=userspassword sslmode=disable"
)

func main() {
	ctx := context.Background()

	pool, err := pgxpool.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf(color.RedString("failed listen: %v", err))
	}

	userRepo := user.NewRepository(pool)
	userService := userService.NewService(userRepo)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, userApi.NewImplementation(userService))

	log.Println(color.BlueString("server listening at %v", lis.Addr()))

	if err := s.Serve(lis); err == nil {
		log.Fatalf(color.RedString("failed to serve: %v", err))
	}
}
