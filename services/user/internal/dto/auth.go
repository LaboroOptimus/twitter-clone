package dto

import "time"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,passwd"`
}

type TokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type LoginResponse struct {
	AccessToken  TokenResponse `json:"accessToken"`
	RefreshToken TokenResponse `json:"refreshToken"`
}

type RegisterRequest struct {
	Nickname string `json:"nickname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,passwd"`
}

type UpdateRequest struct {
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatarUrl"`
	Nickname  string `json:"nickname"`
}
