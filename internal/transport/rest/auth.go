package rest

import (
	"ca/internal/service"
	"ca/internal/transport/dto"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const refreshTokenCookie = "refresh_token"

type AuthHandler struct {
	authService *service.AuthService
	refreshTTL  time.Duration
}

func NewAuthHandler(authService *service.AuthService, refreshTTL time.Duration) *AuthHandler {
	return &AuthHandler{authService: authService, refreshTTL: refreshTTL}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.authService.Login(
		c.Request.Context(), req,
		c.Request.UserAgent(),
		c.ClientIP(),
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Refresh токен — в HttpOnly cookie (недоступен из JS)
	//h.setRefreshCookie(c, refreshToken)

	// Access токен — в JSON ответе (фронт кладёт в память)
	c.JSON(http.StatusOK, dto.AccessTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(15 * 60), // 15 минут в секундах
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newAccess, newRefresh, err := h.authService.Refresh(
		c.Request.Context(), req.RefreshToken,
		c.Request.UserAgent(),
		c.ClientIP(),
	)
	if err != nil {
		h.clearRefreshCookie(c)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.AccessTokenResponse{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
		ExpiresIn:    int64(15 * 60),
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshTokenCookie)
	if err == nil {
		_ = h.authService.Logout(c.Request.Context(), refreshToken)
	}

	h.clearRefreshCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// ──────────────────────────────────────────────
// helpers
// ──────────────────────────────────────────────

func (h *AuthHandler) setRefreshCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		refreshTokenCookie,
		token,
		int(h.refreshTTL.Seconds()),
		"/api/auth", // ограничиваем путь — cookie отправляется только на /api/auth/*
		"",          // domain (пусто = текущий домен)
		true,        // Secure — только HTTPS (в проде!)
		true,        // HttpOnly — недоступен из JS!
	)
}

func (h *AuthHandler) clearRefreshCookie(c *gin.Context) {
	c.SetCookie(refreshTokenCookie, "", -1, "/api/auth", "", true, true)
}
