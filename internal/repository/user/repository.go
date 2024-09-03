package user

import (
	"context"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/color"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/auth/internal/client/db"
	"github.com/Timofey335/auth/internal/model"
	repository "github.com/Timofey335/auth/internal/repository"
	"github.com/Timofey335/auth/internal/repository/user/converter"
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
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

// CreateUser - создает нового пользователя
func (r *repo) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, passwordColumn, roleColumn, createdAtColumn).
		Values(user.Name, user.Email, user.Password, user.Role, time.Now()).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	q := db.Query{
		Name:     "user_repository.CreateUser",
		QueryRaw: query,
	}

	var userId int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

// GetUser - получает данные пользователя
func (r *repo) GetUser(ctx context.Context, userId int64) (*model.User, error) {
	builderSelect := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userId})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.GetUser",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
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
		Where(sq.Eq{idColumn: user.ID})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	qSel := db.Query{
		Name:     "user_repository.SelectUser",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, qSel, args...).Scan(&name, &password, &role)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if user.Name != "" {
		name = user.Name
	}

	if user.Password != "" {
		password = user.Password
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
		Where(sq.Eq{idColumn: user.ID})

	query, args, err = builderUpdate.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.UpdateUser",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
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

// DeleteUser - удаляет пользователя
func (r *repo) DeleteUser(ctx context.Context, userId int64) (*emptypb.Empty, error) {
	var id int64

	builderDelete := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: userId}).
		Suffix("RETURNING id")

	query, args, err := builderDelete.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.DeleteUser",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		log.Println(color.HiMagentaString("error while deleting the user: %v, with ctx: %v", err, ctx))
		return nil, err
	}

	log.Println(color.HiMagentaString("deleted the user: id %v, with ctx: %v", id, ctx))

	return &emptypb.Empty{}, nil
}
