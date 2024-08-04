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

type User struct {
	ID               int64
	Name             string
	Email            string
	Password         string
	Password_confirm string
	Role             desc.Role
	CreatedAt        time.Time
	UpdatedAt        sql.NullTime
}

func (u User) userValidation() error {
	if u.Password != u.Password_confirm {
		return fmt.Errorf("password doesn't match")
	} else {
		return validation.ValidateStruct(&u,
			validation.Field(&u.Name, validation.Required, validation.Length(2, 50)),
			validation.Field(&u.Email, validation.Required, is.Email),
			validation.Field(&u.Password, validation.Required, validation.Length(8, 50)),
		)
	}
}

// CreateUser - create a new user
func (s *server) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.CreateUserResponse, error) {
	user := User{
		Name:             req.Name,
		Email:            req.Email,
		Password:         req.Password,
		Password_confirm: req.PasswordConfirm,
		Role:             req.Role,
	}

	if err := user.userValidation(); err != nil {
		log.Println(color.HiMagentaString("error while creating a new user '%v', email '%v'. %v", user.Name, user.Email, err))

		return nil, err
	} else {
		var userId int64
		err = s.pool.QueryRow(ctx, `INSERT INTO 
									users (name, email, password, role, created_at) 
									VALUES ($1, $2, $3, $4, $5)
									RETURNING id;`, &user.Name, &user.Email, &user.Password, &user.Role, time.Now()).Scan(&userId)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		log.Println(color.BlueString("create user: %v, with ctx: %v", req, ctx))

		return &desc.CreateUserResponse{
			Id: userId,
		}, nil

	}
}

// GetUser - get information of the user by id
func (s *server) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	var user User
	err := s.pool.QueryRow(ctx, `SELECT id, name, email, role, created_at, updated_at 
								FROM users WHERE id = $1`, req.Id).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Println(color.BlueString("Get user by id: %d", req.GetId()))

	var updatedAtTime *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAtTime = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.GetUserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: timestamppb.New(user.CreatedAt),
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
			log.Println(color.HiMagentaString("error while updating the user with id '%v'. %v", req.Id, err))

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

	res, err := s.pool.Exec(ctx, `UPDATE users SET
								name = CASE WHEN $1 = true THEN $2 ELSE name END,
								password = CASE WHEN $3 = true THEN $4 ELSE password END,
								role = CASE WHEN $5 = true THEN $6 ELSE role END,
								updated_at = $7
								WHERE id = $8;`, req.Name != nil, name, req.Password != nil, password, req.Role != nil, role, time.Now(), req.Id)
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
	log.Println(color.HiMagentaString("Delete user: id %d, with ctx: %v", req.Id, ctx))

	return &emptypb.Empty{}, nil
}
