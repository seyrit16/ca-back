package jwtutil

import (
	"ca/internal/model"
	"ca/internal/model/enum"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenClaims struct {
	UserId   uuid.UUID `db:"user_id"`
	Username string    `db:"username"`
	Role     enum.Role `db:"role"`
	jwt.RegisteredClaims
}

type JWTUtil struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

func NewJwtUtil(accessSecret string, refreshSecret string, accessTTL, refreshTTL time.Duration) *JWTUtil {
	return &JWTUtil{
		AccessSecret:  accessSecret,
		RefreshSecret: refreshSecret,
		AccessTTL:     accessTTL,
		RefreshTTL:    refreshTTL,
	}
}

func (jwtUtil *JWTUtil) GenerateAccessToken(user *model.User) (string, error) {
	claims := TokenClaims{
		UserId:   user.Id,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtUtil.AccessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtUtil.AccessSecret))
}

func (jwtUtil *JWTUtil) GenerateRefreshToken(user *model.User) (string, error) {
	claims := TokenClaims{
		UserId:   user.Id,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtUtil.RefreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtUtil.RefreshSecret))
}

func (jwtUtil *JWTUtil) ParseAccessToken(tokenString string) (*TokenClaims, error) {
	return jwtUtil.parseToken(tokenString, jwtUtil.AccessSecret)
}

func (jwtUtil *JWTUtil) ParseRefreshToken(tokenStr string) (*TokenClaims, error) {
	return jwtUtil.parseToken(tokenStr, jwtUtil.RefreshSecret)
}

func (jwtUtil *JWTUtil) parseToken(tokenStr, secret string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
