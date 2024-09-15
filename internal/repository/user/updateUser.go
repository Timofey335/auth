package user

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/auth/internal/client/db"
	"github.com/Timofey335/auth/internal/model"
)

// UpdateUser - update information of the user by id
func (r *repo) UpdateUser(ctx context.Context, user *model.UserUpdateModel) (*emptypb.Empty, error) {
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

	if *user.Name != name && *user.Name != "" {
		name = *user.Name
	}

	if *user.Password != password && *user.Password != "" {
		password = *user.Password
	}

	if *user.Role != 0 {
		role = *user.Role
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

	return &emptypb.Empty{}, nil
}
