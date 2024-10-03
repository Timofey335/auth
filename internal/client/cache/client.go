package cache

import (
	"context"
)

// RedisClient - интерфейс для использования методов redis
type RedisClient interface {
	// HashSet - запись целой структуры (ключ-значение) по ключу
	HashSet(ctx context.Context, key string, values interface{}) error

	// Set - записать значение по ключу
	Set(ctx context.Context, key string, value interface{}) error

	// HGetAll - получение всех полей и значений из хэш-структуры данных
	HGetAll(ctx context.Context, key string) ([]interface{}, error)

	// GET - получение значения по ключу в бд
	Get(ctx context.Context, key string) (interface{}, error)

	// DeleteHashSet - удаляет Hash set в redis по ключу
	DeleteHashSet(ctx context.Context, key string) error

	// Expire - устанавлиает время жизни для ключа
	Expire(ctx context.Context, key string) error

	// Ping - проверка соединения
	Ping(ctx context.Context) error
}
