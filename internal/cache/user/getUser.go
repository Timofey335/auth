package user

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/gomodule/redigo/redis"

	"github.com/Timofey335/auth/internal/cache/user/converter"
	cacheModel "github.com/Timofey335/auth/internal/cache/user/model"
	"github.com/Timofey335/auth/internal/model"
)

// GetUser - получает информацию о пользователе из кэша
func (c *cacheImplementation) GetUser(ctx context.Context, id int64) (*model.UserModel, error) {
	idStr := strconv.Itoa(int(id))
	userCache, err := c.cacheClient.HGetAll(ctx, idStr)
	if err != nil {
		return nil, err
	}

	if len(userCache) == 0 {
		return nil, errors.New("user not found in the cache")
	}

	var user cacheModel.UserCacheModel
	if err = redis.ScanStruct(userCache, &user); err != nil {
		log.Fatalf("failed to scan fields: %v", err)
	}

	return converter.ToUserModelFromUserCache(&user), nil
}
