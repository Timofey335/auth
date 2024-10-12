package model

import "github.com/dgrijalva/jwt-go"

const (
	ExamplePath = "/auth_v1.Auth_v1/GetUser"
)

type UserClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
	Role  int64  `json:"role"`
}
