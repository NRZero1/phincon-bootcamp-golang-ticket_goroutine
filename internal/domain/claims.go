package domain

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Email string `json:"Email"`
	Name string `json:"Name"`
	jwt.RegisteredClaims
}