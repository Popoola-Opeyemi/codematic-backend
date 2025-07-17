package model

type AuthProvider string

type Providers struct {
}

const (
	ProviderEmail   AuthProvider = "email"
	ProviderGoogle  AuthProvider = "google"
	ProviderApple   AuthProvider = "apple"
	ProviderDiscord AuthProvider = "discord"
	ProviderWallet  AuthProvider = "wallet"
)
