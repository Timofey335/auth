package users

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Timofey335/auth/internal/model"
	"github.com/Timofey335/auth/internal/repository/user/converter"
	"github.com/fatih/color"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	repository "github.com/Timofey335/auth/internal/repository"
	modelRepo "github.com/Timofey335/auth/internal/repository/user/model"
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

func NewRepository(db *pgxpool.Pool) repository.UsersRepository {
	return &repo{db: db}
}

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

func (r *repo) GetUser(ctx context.Context, userId int64) (*model.User, error) {
	var user modelRepo.User

	err := r.db.QueryRow(ctx, `SELECT id, name, email, role, created_at, updated_at
		FROM users WHERE id = $1`, userId).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return converter.ToUserFromRepo(&user), nil
}
