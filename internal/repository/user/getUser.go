package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/Timofey335/platform_common/pkg/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Timofey335/auth/internal/model"
	"github.com/Timofey335/auth/internal/repository/user/converter"
	modelRepo "github.com/Timofey335/auth/internal/repository/user/model"
)

// GetUser - получает данные пользователя
func (r *repo) GetUser(ctx context.Context, userId int64) (*model.UserModel, error) {
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

	var user modelRepo.UserRepoModel
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return converter.ToUserFromRepo(&user), nil
}
