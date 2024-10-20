package model

import "github.com/dgrijalva/jwt-go"

// UserClaims - модель для работы c jwt токеном
type UserClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
	Role  int64  `json:"role"`
}
