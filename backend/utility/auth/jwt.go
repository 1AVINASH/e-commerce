package auth

import (
	"context"
	"fmt"
	"time"

	redisclient "gotemplate/infra/redis"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("supersecretkey")

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// HashPassword hashes the plain password using bcrypt
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// CheckPasswordHash compares hashed password with plain text password
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

func GenerateJWT(email string) (string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	redisclient.RedisClient.Set(context.Background(), email, token, time.Hour)
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}
	username := claims.Username
	jwtInCache := redisclient.RedisClient.Get(context.Background(), username)
	if jwtInCache==nil {
		return nil, fmt.Errorf("Please login again")
	}
	return claims, nil
}
