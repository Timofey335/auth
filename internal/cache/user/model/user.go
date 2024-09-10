package model

type UserCacheModel struct {
	ID              int64  `redis:"id"`
	Name            string `redis:"name"`
	Email           string `redis:"email"`
	Password        string `redis:"password"`
	PasswordConfirm string `redis:"password_confirm"`
	Role            int64  `redis:"role"`
	CreatedAt       int64  `redis:"created_at"`
	UpdatedAt       *int64 `redis:"updated_at"`
}
