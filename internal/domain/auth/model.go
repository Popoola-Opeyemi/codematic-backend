package auth

import "codematic/internal/shared/model"

type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	Provider   string `json:"provider"`
	SocialID   string `json:"social_id,omitempty"`
	WalletAddr string `json:"wallet_address,omitempty"`
}

type SignupRequest struct {
	Email       string             `json:"email"`
	Password    string             `json:"password"`
	Provider    model.AuthProvider `json:"provider"`
	SocialID    string             `json:"social_id,omitempty"`
	DisplayName string             `json:"display_name,omitempty"`
	AvatarURL   string             `json:"avatar_url,omitempty"`
	WalletAddr  string             `json:"wallet_address,omitempty"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
