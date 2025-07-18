package cache

import "context"

type CacheManager interface {
	SessionStore
	GetTokenIDForUser(ctx context.Context, userID string) (string, error)
}

type unifiedCacheManager struct {
	SessionStore
}

func NewCacheManager(sessionStore SessionStore) CacheManager {
	return &unifiedCacheManager{
		SessionStore: sessionStore,
	}
}

func (u *unifiedCacheManager) GetTokenIDForUser(ctx context.Context, userID string) (string, error) {
	return u.SessionStore.GetTokenIDForUser(ctx, userID)
}
