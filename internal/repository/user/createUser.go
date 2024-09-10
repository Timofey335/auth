package user

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/Timofey335/auth/internal/client/db"
	"github.com/Timofey335/auth/internal/model"
)

// CreateUser - создает нового пользователя
func (r *repo) CreateUser(ctx context.Context, user *model.UserModel) (int64, error) {
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
