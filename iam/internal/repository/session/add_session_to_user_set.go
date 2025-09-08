package session

import "context"

func (r *repository) AddSessionToUserSet(ctx context.Context, userUUID, sessionUUID string) error {
	cacheKey := r.getCacheKey(userUUID)
	return r.cache.SAdd(ctx, cacheKey, sessionUUID)
}
