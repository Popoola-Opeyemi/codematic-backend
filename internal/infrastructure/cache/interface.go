package cache

type CacheManager interface {
	SessionStore
}

type unifiedCacheManager struct {
	SessionStore
}

func NewCacheManager(sessionStore SessionStore) CacheManager {
	return &unifiedCacheManager{
		SessionStore: sessionStore,
	}
}
