package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/fatih/color"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/Timofey335/auth/pkg/authservice_v1"
)

const (
	grpcPort = 50051
	dbDSN    = "host=localhost port=54321 dbname=users user=user password=userspassword sslmode=disable"
)

type server struct {
	desc.UnimplementedAuthserviceV1Server
	pool *pgxpool.Pool
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

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterAuthserviceV1Server(s, &server{pool: pool})
	log.Println(color.BlueString("server listening at %v", lis.Addr()))
	if err := s.Serve(lis); err == nil {
		log.Fatalf(color.RedString("failed to serve: %v", err))
	}
}

// CreateUser - create a new user
func (s *server) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	err := validation.Validate(req.Name, validation.Required, validation.Length(2, 50))
	if err != nil {
		log.Println(color.HiMagentaString("error while creating the new user: %v, with ctx: %v", err, ctx))

		return nil, err
	}

	err = validation.Validate(req.Email, validation.Required, is.Email)
	if err != nil {
		log.Println(color.HiMagentaString("error while creating the new user: %v, with ctx: %v", err, ctx))

		return nil, err
	}

	if req.Password != req.PasswordConfirm {
		err := errors.New("password doesn't match")
		log.Println(color.HiMagentaString("error while creating the new user: %v, with ctx: %v", err, ctx))

		return nil, err
	}

	err = validation.Validate(req.Password, validation.Required, validation.Length(8, 50))
	if err != nil {
		log.Println(color.HiMagentaString("error while creating the new user: %v, with ctx: %v", err, ctx))

		return nil, err
	}

	var userId int64
	err = s.pool.QueryRow(ctx, `INSERT INTO users (name, email, password, role, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
		req.Name, req.Email, req.Password, req.Role, time.Now()).Scan(&userId)
	if err != nil {
		return nil, err
	}

	log.Println(color.BlueString("create user: %v, with ctx: %v", req, ctx))

	return &desc.CreateUserResponse{
		Id: userId,
	}, nil
}

// GetUser - get information of the user by id
func (s *server) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	var id int64
	var name, email string
	var role desc.Role
	var createdAt time.Time
	var updatedAt sql.NullTime

	err := s.pool.QueryRow(ctx, `SELECT id, name, email, role, created_at, updated_at 
	FROM users WHERE id = $1`, req.Id).Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Println(color.BlueString("Get user by id: %d", req.GetId()))

	var updatedAtTime *timestamppb.Timestamp
	if updatedAt.Valid {
		updatedAtTime = timestamppb.New(updatedAt.Time)
	}

	return &desc.GetUserResponse{
		Id:        id,
		Name:      name,
		Email:     email,
		Role:      role,
		CreatedAt: timestamppb.New(createdAt),
		UpdatedAt: updatedAtTime,
	}, nil
}

// UpdateUser - update information of the user by id
func (s *server) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*emptypb.Empty, error) {
	var name, password string
	var role desc.Role

	if req.Name != nil {
		name = req.Name.Value
		err := validation.Validate(name, validation.Required, validation.Length(2, 50))
		if err != nil {
			log.Println(color.HiMagentaString("error while updating the user with id '%v'; %v", req.Id, err))

			return nil, err
		}
	}

	if req.Password != nil {
		password = req.Password.Value
		if password != req.PasswordConfirm.Value {
			err := errors.New("password doesn't match")
			log.Println(color.HiMagentaString("error while updating password the user with id '%v; %v'", req.Id, err))

			return nil, err
		}

		err := validation.Validate(&password, validation.Required, validation.Length(8, 50))
		if err != nil {
			log.Println(color.HiMagentaString("error while updating password the user with id '%v'; %v", req.Id, err))

			return nil, err
		}

	}

	if req.Role != nil {
		role = *req.Role
	}

	res, err := s.pool.Exec(ctx, `UPDATE users SET name = CASE WHEN $1 = true THEN $2 ELSE name END, password = CASE WHEN $3 = true THEN $4 ELSE password END,
	role = CASE WHEN $5 = true THEN $6 ELSE role END, updated_at = $7 WHERE id = $8;`,
		req.Name != nil, name, req.Password != nil, password, req.Role != nil, role, time.Now(), req.Id)
	if err != nil {
		return nil, err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return nil, fmt.Errorf("updating failed")
	}

	log.Println(color.BlueString("updated the user %v, with ctx: %v", req, ctx))

	return &emptypb.Empty{}, nil
}

// DeleteUser - delete a user by id
func (s *server) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	var userId int64
	err := s.pool.QueryRow(ctx, `DELETE FROM users WHERE id = $1 RETURNING id`, req.Id).Scan(&userId)
	if err != nil {
		log.Println(color.HiMagentaString("error while deleting the user: %v, with ctx: %v", err, ctx))
		return nil, err
	}

	log.Println(color.HiMagentaString("deleted the user: id %v, with ctx: %v", userId, ctx))

	return &emptypb.Empty{}, nil
}
