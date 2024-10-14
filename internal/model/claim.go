package model

import "github.com/dgrijalva/jwt-go"

const (
	ExamplePath = "/chat_server_v1.Chat_server_v1/CreateChat"
)

type UserClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
	Role  int64  `json:"role"`
}
