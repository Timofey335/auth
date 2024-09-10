package user

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/color"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/Timofey335/auth/internal/client/db"
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
		log.Println(color.HiMagentaString("error while deleting the user: %v, with ctx: %v", err, ctx))
		return nil, err
	}

	log.Println(color.HiMagentaString("deleted the user: id %v, with ctx: %v", id, ctx))

	return &emptypb.Empty{}, nil
}
