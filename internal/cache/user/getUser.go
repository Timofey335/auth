package user

import (
	"context"

	"github.com/Timofey335/auth/internal/model"
)

func (c *cacheImplementation) GetUser(ctx context.Context, id int64) (*model.UserModel, error) {
	return nil, nil
}
