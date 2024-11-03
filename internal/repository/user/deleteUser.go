package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/Timofey335/platform_common/pkg/db"
	"google.golang.org/protobuf/types/known/emptypb"
)

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
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
