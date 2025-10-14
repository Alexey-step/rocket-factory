package user

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/Alexey-step/rocket-factory/iam/internal/model"
	repoConverter "github.com/Alexey-step/rocket-factory/iam/internal/repository/converter"
	repoModel "github.com/Alexey-step/rocket-factory/iam/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, info model.UserInfo) (userUUID string, err error) {
	user := repoModel.User{
		UUID:                uuid.NewString(),
		Login:               info.Login,
		Email:               info.Email,
		NotificationMethods: repoConverter.NotificationMethodsToRepo(ctx, info.NotificationMethods),
		CreatedAt:           time.Now(),
		Password:            info.Password,
	}

	query, args, err := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("uuid", "login", "email", "notification_methods", "created_at", "password_hash").
		Values(user.UUID, user.Login, user.Email, user.NotificationMethods, user.CreatedAt, user.Password).
		Suffix("RETURNING uuid").
		ToSql()
	if err != nil {
		return "", err
	}

	var userUUIDResult string
	err = r.db.QueryRow(ctx, query, args...).Scan(&userUUIDResult)
	if err != nil {
		return "", err
	}

	return userUUIDResult, nil
}
