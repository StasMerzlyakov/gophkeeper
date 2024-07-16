package domain

import (
	"github.com/golang-jwt/jwt/v4"
)

type Token string

type Claims struct {
	jwt.RegisteredClaims
	UserID UserID
}
