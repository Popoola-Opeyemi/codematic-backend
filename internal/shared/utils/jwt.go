package utils

import (
	"codematic-backend/internal/shared/model"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	JWTSecret          []byte
	RefreshTokenSecret []byte
}

func NewJWTManager(jwtSecret, refreshTokenSecret string) *JWTManager {
	return &JWTManager{
		JWTSecret:          []byte(jwtSecret),
		RefreshTokenSecret: []byte(refreshTokenSecret),
	}
}

func (j *JWTManager) GenerateJWT(userID, tokenID string) (string, error) {
	claims := model.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ID:        tokenID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.JWTSecret)
}

func (j *JWTManager) GenerateRefreshToken(userID string, tokenID string) (string, error) {
	claims := model.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ID:        tokenID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.RefreshTokenSecret)
}

func (j *JWTManager) ParseJWT(tokenStr string) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

func (j *JWTManager) ParseRefreshToken(tokenStr string) (*model.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.RefreshTokenSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}

// VerifyRefreshToken validates a refresh token and returns its claims
func (j *JWTManager) VerifyRefreshToken(refreshToken string) (*model.Claims, error) {
	claims, err := j.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Optional: you can also manually check expiry, although jwt.ParseWithClaims already does it internally
	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now().UTC()) {
		return nil, errors.New("refresh token has expired")
	}

	return claims, nil
}

// utils/jwt.go
func (j *JWTManager) ExtractTokenFromHeader(authHeader string) string {
	const prefix = "Bearer "
	if len(authHeader) > len(prefix) && authHeader[:len(prefix)] == prefix {
		return authHeader[len(prefix):]
	}
	return ""
}
