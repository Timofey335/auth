package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
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

	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn, createdAtColumn).
		Values(user.Name, user.Email, user.Password, user.Role, time.Now()).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	var userId int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&userId)
	if err != nil {
		return 0, err
	}

	log.Println(color.BlueString("create user: %v, with ctx: %v", user, ctx))

	return userId, nil
}

// GetUser - get information of the user by id
func (r *repo) GetUser(ctx context.Context, userId int64) (*model.User, error) {
	var user modelRepo.User

	builderSelect := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": userId})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return converter.ToUserFromRepo(&user), nil
}

// UpdateUser - update information of the user by id
func (r *repo) UpdateUser(ctx context.Context, user *model.User) (*emptypb.Empty, error) {
	var name, password string
	var role int64

	builderSelect := sq.Select(nameColumn, passwordColumn, roleColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": user.ID})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&name, &password, &role)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

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

	builderUpdate := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(nameColumn, name).
		Set(passwordColumn, password).
		Set(roleColumn, role).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{"id": user.ID})

	query, args, err = builderUpdate.ToSql()
	if err != nil {
		return nil, err
	}

	res, err := r.db.Exec(ctx, query, args...)
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

	builderDelete := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": userId}).
		Suffix("RETURNING id")

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		log.Println(color.HiMagentaString("error while deleting the user: %v, with ctx: %v", err, ctx))
		return nil, err
	}

	log.Println(color.HiMagentaString("deleted the user: id %v, with ctx: %v", id, ctx))

	return &emptypb.Empty{}, nil
}
