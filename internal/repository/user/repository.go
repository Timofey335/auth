package users

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/auth/internal/model"
	repository "github.com/Timofey335/auth/internal/repository"
	"github.com/Timofey335/auth/internal/repository/user/converter"
	modelRepo "github.com/Timofey335/auth/internal/repository/user/model"
	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	tableName = "users"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	passwordColumn  = "password"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{db: db}
}

// CreateUser - create a new user
func (r *repo) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	err := validation.Validate(user.Name, validation.Required, validation.Length(2, 50))
	if err != nil {
		log.Println(color.HiMagentaString("error while creating the new user: %v, with ctx: %v", err, ctx))

		return 0, err
	}

	err = validation.Validate(user.Email, validation.Required, is.Email)
	if err != nil {
		log.Println(color.HiMagentaString("error while creating the new user: %v, with ctx: %v", err, ctx))

		return 0, err
	}

	if user.Password != user.PasswordConfirm {
		err := errors.New("password doesn't match")
		log.Println(color.HiMagentaString("error while creating the new user: %v, with ctx: %v", err, ctx))

		return 0, err
	}

	err = validation.Validate(user.Password, validation.Required, validation.Length(8, 50))
	if err != nil {
		log.Println(color.HiMagentaString("error while creating the new user: %v, with ctx: %v", err, ctx))

		return 0, err
	}

	var userId int64
	err = r.db.QueryRow(ctx, `INSERT INTO users (name, email, password, role, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
		user.Name, user.Email, user.Password, user.Role, time.Now()).Scan(&userId)
	if err != nil {
		return 0, err
	}

	log.Println(color.BlueString("create user: %v, with ctx: %v", user, ctx))

	return userId, nil
}

// GetUser - get information of the user by id
func (r *repo) GetUser(ctx context.Context, userId int64) (*model.User, error) {
	var user modelRepo.User

	err := r.db.QueryRow(ctx, `SELECT id, name, email, role, created_at, updated_at
		FROM users WHERE id = $1`, userId).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return converter.ToUserFromRepo(&user), nil
}

// UpdateUser - update information of the user by id
func (r *repo) UpdateUser(ctx context.Context, user *model.User) (*emptypb.Empty, error) {
	var name, password string
	var role int64

	if user.Name != "" {
		name = user.Name
		err := validation.Validate(name, validation.Required, validation.Length(2, 50))
		if err != nil {
			log.Println(color.HiMagentaString("error while updating the user with id '%v'; %v", user.ID, err))

			return nil, err
		}
	}

	if user.Password != "" {
		password = user.Password
		if password != user.PasswordConfirm {
			err := errors.New("password doesn't match")
			log.Println(color.HiMagentaString("error while updating password the user with id '%v; %v'", user.ID, err))

			return nil, err
		}

		err := validation.Validate(&password, validation.Required, validation.Length(8, 50))
		if err != nil {
			log.Println(color.HiMagentaString("error while updating password the user with id '%v'; %v", user.ID, err))

			return nil, err
		}

	}

	if user.Role != 0 {
		role = user.Role
	}

	res, err := r.db.Exec(ctx, `UPDATE users SET name = CASE WHEN $1 = true THEN $2 ELSE name END, password = CASE WHEN $3 = true THEN $4 ELSE password END,
	role = CASE WHEN $5 = true THEN $6 ELSE role END, updated_at = $7 WHERE id = $8;`,
		user.Name != "", name, user.Password != "", password, user.Role != 0, role, time.Now(), user.ID)
	if err != nil {
		return nil, err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return nil, fmt.Errorf("updating failed")
	}

	log.Println(color.BlueString("updated the user %v, with ctx: %v", user.ID, ctx))

	return &emptypb.Empty{}, nil
}

// DeleteUser - delete a user by id
func (r *repo) DeleteUser(ctx context.Context, userId int64) (*emptypb.Empty, error) {
	var id int64
	err := r.db.QueryRow(ctx, `DELETE FROM users WHERE id = $1 RETURNING id`, userId).Scan(&id)
	if err != nil {
		log.Println(color.HiMagentaString("error while deleting the user: %v, with ctx: %v", err, ctx))
		return nil, err
	}

	log.Println(color.HiMagentaString("deleted the user: id %v, with ctx: %v", id, ctx))

	return &emptypb.Empty{}, nil
}
