package cache

type CacheManager interface {
	SessionStore
	ProviderCacheStore
}

type unifiedCacheManager struct {
	SessionStore
	ProviderCacheStore
}

func NewCacheManager(sessionStore SessionStore,
	providerCacheStore ProviderCacheStore) CacheManager {
	return &unifiedCacheManager{
		SessionStore:       sessionStore,
		ProviderCacheStore: providerCacheStore,
	}
}
