package session

import (
	"context"
	"errors"

	redigo "github.com/gomodule/redigo/redis"

	"github.com/Alexey-step/rocket-factory/iam/internal/model"
	repoConverter "github.com/Alexey-step/rocket-factory/iam/internal/repository/converter"
	repoModel "github.com/Alexey-step/rocket-factory/iam/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, sessionUUID string) (model.Session, model.User, error) {
	cacheKey := r.getCacheKey(sessionUUID)

	values, err := r.cache.HGetAll(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return model.Session{}, model.User{}, model.ErrSessionNotFound
		}
		return model.Session{}, model.User{}, err
	}

	if len(values) == 0 {
		return model.Session{}, model.User{}, model.ErrSessionNotFound
	}

	var sessionRedisView repoModel.SessionRedisView
	err = redigo.ScanStruct(values, &sessionRedisView)
	if err != nil {
		return model.Session{}, model.User{}, err
	}

	session, user := repoConverter.SessionAndUserFromRedisView(sessionRedisView)

	return session, user, nil
}
