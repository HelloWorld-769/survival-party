package model

import (
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	Id string
	jwt.RegisteredClaims
}
