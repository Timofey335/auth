package model

import "github.com/dgrijalva/jwt-go"

const (
	ExamplePath = "/note_v1.NoteV1/Get"
)

type UserClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
	Role  int64  `json:"role"`
}
