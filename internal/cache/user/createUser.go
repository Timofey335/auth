package user

import (
	"context"
	"strconv"

	"github.com/Timofey335/auth/internal/cache/user/converter"
	"github.com/Timofey335/auth/internal/model"
)

// CreateUser - создает пользователя в кэше
func (c *cacheImplementation) CreateUser(ctx context.Context, user *model.UserModel) error {
	userCache := converter.ToUserCacheFromUserModel(user)
	id := strconv.FormatInt(userCache.ID, 10)
	if err := c.cacheClient.HashSet(ctx, id, userCache); err != nil {
		return err
	}

	// Устанавливает TTL для кэшированных данных пользователя 3600 секунд
	if err := c.cacheClient.Expire(ctx, id, 3600); err != nil {
		return err
	}

	return nil
}
