package model

import "errors"

var (
	ErrWalletAlreadyExists                 = errors.New("this wallet address has already been added to your portfolio")
	ErrUserNotFound                        = errors.New("user not found")
	ErrMissingOrInvalidAuthorizationHeader = errors.New("missing or invalid authorization header")
	ErrInvalidOrExpiredToken               = errors.New("invalid or expired token")
	ErrTokenRevoked                        = errors.New("token revoked")
	ErrSessionNotFound                     = errors.New("session not found")
	ErrMissingXTenantIDHeader              = errors.New("missing X-Tenant-ID header")
	ErrInvalidTenantIDFormat               = errors.New("invalid tenant ID format")
	ErrInvalidInputError                   = errors.New("invalid input")
	ErrMissingRequiredFields               = errors.New("missing required fields")
	ErrInvalidEmailFormat                  = errors.New("invalid email format")
	ErrPasswordTooShort                    = errors.New("password too short (min 8 chars)")
	ErrUnsupportedProvider                 = errors.New("unsupported provider")

	ErrInvalidSignature = errors.New("invalid webhook signature")
)
