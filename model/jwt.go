package model

import (
	"os"

	"github.com/golang-jwt/jwt"
)

var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}
