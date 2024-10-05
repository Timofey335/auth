package user

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/Timofey335/platform_common/pkg/db"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	modelRepo "github.com/Timofey335/auth/internal/repository/user/model"
)

// GetUserData - получает данные о пользователя из базы
func (r *repo) GetUserData(ctx context.Context, userEmail string) (*modelRepo.UserRepoModel, error) {
	builderSelect := sq.Select(passwordColumn, roleColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{emailColumn: userEmail})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.GetUserData",
		QueryRaw: query,
	}

	var user modelRepo.UserRepoModel
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &modelRepo.UserRepoModel{
		ID:       user.ID,
		Password: user.Password,
		Role:     user.Role,
	}, nil
}
