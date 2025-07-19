package model

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID   string `json:"sub"`
	Email    string `json:"email"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type JWTData struct {
	UserID   string
	Email    string
	TenantID string
	TokenID  string
	Role     string
}
