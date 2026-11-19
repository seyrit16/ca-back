package dto

import (
	"github.com/google/uuid"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // секунды
}

type UserCreateRequest struct {
	Nickname string `json:"nickname"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Color    string `json:"color"`
}

type UserUpdateProfileDataRequest struct {
	Id       uuid.UUID `json:"id"`
	Nickname *string   `json:"nickname"`
	Username *string   `json:"username"`
	Email    *string   `json:"email"`
	Color    *string   `json:"color"`
}

// Draft
type DraftPickRequest struct {
	PlayerId   string `json:"playerId"`
	TemplateId int16  `json:"templateId"`
}
