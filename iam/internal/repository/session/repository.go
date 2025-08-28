package session

import (
	"fmt"

	iamRepo "github.com/Alexey-step/rocket-factory/iam/internal/repository"
	"github.com/Alexey-step/rocket-factory/platform/pkg/cache"
)

var _ iamRepo.SessionRepository = (*repository)(nil)

const (
	cacheKeyPrefix = "iam:session:"
)

type repository struct {
	cache cache.RedisClient
}

func NewRepository(cache cache.RedisClient) *repository {
	return &repository{
		cache: cache,
	}
}

func (r *repository) GetCacheKey(uuid string) string {
	return fmt.Sprintf("%s%s", cacheKeyPrefix, uuid)
}
