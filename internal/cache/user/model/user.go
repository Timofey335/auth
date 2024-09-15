package model

// UserCacheModel - модель UserCacheModel
type UserCacheModel struct {
	ID        int64  `redis:"id"`
	Name      string `redis:"name"`
	Email     string `redis:"email"`
	Role      int64  `redis:"role"`
	CreatedAt int64  `redis:"created_at"`
	UpdatedAt *int64 `redis:"updated_at"`
}
