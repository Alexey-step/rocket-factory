package session

import (
	"context"
	"time"

	"github.com/Alexey-step/rocket-factory/iam/internal/model"
	repoConverter "github.com/Alexey-step/rocket-factory/iam/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, session model.Session, user model.User, ttl time.Duration) error {
	cacheKey := r.getCacheKey(session.UUID)

	redisView := repoConverter.SessionAndUserToRedisView(session, user)

	err := r.cache.HashSet(ctx, cacheKey, redisView)
	if err != nil {
		return err
	}

	return r.cache.Expire(ctx, cacheKey, ttl)
}
