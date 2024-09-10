package user

import (
	"github.com/Timofey335/auth/internal/client/db"
	repository "github.com/Timofey335/auth/internal/repository"
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

// NewRepository - создает новый объект repo
func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}
