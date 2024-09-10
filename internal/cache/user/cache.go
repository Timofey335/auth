package user

import (
	"github.com/Timofey335/auth/internal/cache"
	cacheClient "github.com/Timofey335/auth/internal/client/cache"
)

type cacheImplementation struct {
	cacheClient cacheClient.RedisClient
}

func NewCache(cacheClient cacheClient.RedisClient) cache.UserCache {
	return &cacheImplementation{
		cacheClient: cacheClient,
	}
}
