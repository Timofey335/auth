package access

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Timofey335/platform_common/pkg/db"
)

// GetRole - возвращает значение роли в соответствии с эндпоинтом
func (r *repo) GetRole(ctx context.Context, endpoint string) (int64, error) {
	builderSelect := sq.Select(roleColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{endpointColumn: endpoint})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "access_repository.GetRole",
		QueryRaw: query,
	}

	var role int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&role)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return role, nil
}
