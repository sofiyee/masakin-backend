package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID int, name, role string) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	claims := JwtClaims{
		UserID: userID,
		Name:   name,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}


func ParseToken(tokenStr string) (*JwtClaims, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.ParseWithClaims(
		tokenStr,
		&JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
