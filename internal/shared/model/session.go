package model

import "time"

type UserSessionInfo struct {
	UserID    string    `json:"user_id"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`
	TokenID   string    `json:"token_id"`
	LoginTime time.Time `json:"login_time"`
	LastSeen  time.Time `json:"last_seen"`
	IsActive  bool      `json:"is_active"`
}
