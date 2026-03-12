package domain

import "github.com/golang-jwt/jwt/v5"

type PasswordService interface {
	HashPassword(password string) (string, error)
	ComparePassword(hash string, password string) error
}

type JWTService interface {
	GenerateAccessToken(userID string, role string) (string, error)
	GenerateRefreshToken(userID string, role string) (string, error)
	ValidateToken(token string) (jwt.MapClaims, error)
}