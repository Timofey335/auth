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

	repository "github.com/Timofey335/auth/internal/repository"
	user "github.com/Timofey335/auth/internal/repository/user"
	desc "github.com/Timofey335/auth/pkg/auth_v1"
)

const (
	grpcPort = 50051
	dbDSN    = "host=localhost port=54321 dbname=users user=user password=userspassword sslmode=disable"
)

type server struct {
	desc.UnimplementedAuthV1Server
	usersRepository repository.UsersRepository
}

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

	usersRepo := user.NewRepository(pool)

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthV1Server(s, &server{usersRepository: usersRepo})

	log.Println(color.BlueString("server listening at %v", lis.Addr()))

	if err := s.Serve(lis); err == nil {
		log.Fatalf(color.RedString("failed to serve: %v", err))
	}
}

func (s *server) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	id, err := s.usersRepository.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	log.Println(color.BlueString("create user: %v, with ctx: %v", req, ctx))

	return &desc.CreateUserResponse{
		Id: id,
	}, nil
}

func (s *server) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	userObj, err := s.usersRepository.GetUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	log.Println(color.BlueString("Get user by id: %d", userObj.Id))

	return &desc.GetUserResponse{
		Id:        userObj.Id,
		Name:      userObj.Name,
		Email:     userObj.Email,
		Role:      userObj.Role,
		CreatedAt: userObj.CreatedAt,
		UpdatedAt: userObj.UpdatedAt,
	}, nil

}

// UpdateUser - update information of the user by id
// func (s *server) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
// 	var name, password string
// 	var role desc.Role

// 	if req.Name != nil {
// 		name = req.Name.Value
// 		err := validation.Validate(name, validation.Required, validation.Length(2, 50))
// 		if err != nil {
// 			log.Println(color.HiMagentaString("error while updating the user with id '%v'; %v", req.Id, err))

// 			return nil, err
// 		}
// 	}

// 	if req.Password != nil {
// 		password = req.Password.Value
// 		if password != req.PasswordConfirm.Value {
// 			err := errors.New("password doesn't match")
// 			log.Println(color.HiMagentaString("error while updating password the user with id '%v; %v'", req.Id, err))

// 			return nil, err
// 		}

// 		err := validation.Validate(&password, validation.Required, validation.Length(8, 50))
// 		if err != nil {
// 			log.Println(color.HiMagentaString("error while updating password the user with id '%v'; %v", req.Id, err))

// 			return nil, err
// 		}

// 	}

// 	if req.Role != nil {
// 		role = *req.Role
// 	}

// 	res, err := s.pool.Exec(ctx, `UPDATE users SET name = CASE WHEN $1 = true THEN $2 ELSE name END, password = CASE WHEN $3 = true THEN $4 ELSE password END,
// 	role = CASE WHEN $5 = true THEN $6 ELSE role END, updated_at = $7 WHERE id = $8;`,
// 		req.Name != nil, name, req.Password != nil, password, req.Role != nil, role, time.Now(), req.Id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	rowsAffected := res.RowsAffected()
// 	if rowsAffected == 0 {
// 		return nil, fmt.Errorf("updating failed")
// 	}

// 	log.Println(color.BlueString("updated the user %v, with ctx: %v", req, ctx))

// 	return &emptypb.Empty{}, nil
// }

// DeleteUser - delete a user by id
// func (s *server) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
// 	var userId int64
// 	err := s.pool.QueryRow(ctx, `DELETE FROM users WHERE id = $1 RETURNING id`, req.Id).Scan(&userId)
// 	if err != nil {
// 		log.Println(color.HiMagentaString("error while deleting the user: %v, with ctx: %v", err, ctx))
// 		return nil, err
// 	}

// 	log.Println(color.HiMagentaString("deleted the user: id %v, with ctx: %v", userId, ctx))

// 	return &emptypb.Empty{}, nil
// }
