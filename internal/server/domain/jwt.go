package domain

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Token string

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}
