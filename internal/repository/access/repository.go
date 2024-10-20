package access

import (
	"github.com/Timofey335/platform_common/pkg/db"

	repository "github.com/Timofey335/auth/internal/repository"
)

const (
	tableName = "roles"

	idColumn       = "id"
	endpointColumn = "endpoint"
	roleColumn     = "role"
)

type repo struct {
	db db.Client
}

// NewRepository - создает новый объект repo
func NewRepository(db db.Client) repository.AccessRepository {
	return &repo{db: db}
}
