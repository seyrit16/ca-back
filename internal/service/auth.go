package service

import (
	"ca/internal/model"
	"ca/internal/repository/postgres"
	"ca/internal/transport/dto"
	"ca/pkg/jwtutil"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo    postgres.AuthRepository
	jwtUtil *jwtutil.JWTUtil
}

func NewAuthService(repo postgres.AuthRepository, jwtUtil *jwtutil.JWTUtil) *AuthService {
	return &AuthService{
		repo:    repo,
		jwtUtil: jwtUtil,
	}
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest, userAgent, ip string) (string, string, error) {
	user, err := s.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		log.Println("GetUserByUsername error:", err)
		return "", "", errors.New("invalid credentials")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		log.Println("bcrypt error:", err)
		return "", "", errors.New("invalid credentials")
	}

	return s.generateAndSaveTokens(ctx, user, userAgent, ip)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken, userAgent, ip string) (string, string, error) {
	// Валидируем refresh токен
	claims, err := s.jwtUtil.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	// Проверяем наличие в БД
	tokenHash := hashToken(refreshToken)
	session, err := s.repo.GetRefreshSession(ctx, tokenHash)
	if err != nil {
		return "", "", errors.New("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		_ = s.repo.DeleteRefreshSession(ctx, tokenHash)
		return "", "", errors.New("refresh token expired")
	}

	// Удаляем старую сессию
	_ = s.repo.DeleteRefreshSession(ctx, tokenHash)

	user, err := s.repo.GetUserByUsername(ctx, claims.Username)
	if err != nil {
		return "", "", err
	}

	return s.generateAndSaveTokens(ctx, user, userAgent, ip)
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	tokenHash := hashToken(refreshToken)
	return s.repo.DeleteRefreshSession(ctx, tokenHash)
}

func (s *AuthService) generateAndSaveTokens(ctx context.Context, user *model.User, userAgent, ip string) (string, string, error) {
	accessToken, err := s.jwtUtil.GenerateAccessToken(user)
	if err != nil {

		log.Println("GenerateAccessToken error:", err)
		return "", "", err
	}

	refreshToken, err := s.jwtUtil.GenerateRefreshToken(user)
	if err != nil {

		log.Println("GenerateRefreshToken error:", err)
		return "", "", err
	}

	// Сохраняем хэш refresh токена
	err = s.repo.SaveRefreshSession(ctx, &model.RefreshSession{
		UserID:    user.Id,
		Token:     hashToken(refreshToken),
		ExpiresAt: time.Now().Add(s.jwtUtil.RefreshTTL),
		UserAgent: userAgent,
		IP:        ip,
	})
	if err != nil {

		log.Println("SaveRefreshSession error:", err)
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
