package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/Alexey-step/rocket-factory/iam/internal/model"
	repoConverter "github.com/Alexey-step/rocket-factory/iam/internal/repository/converter"
	repoModel "github.com/Alexey-step/rocket-factory/iam/internal/repository/model"
)

func (r *repository) GetUserByLogin(ctx context.Context, login, password string) (model.User, error) {
	query, args, err := sq.Select(
		"uuid",
		"login",
		"email",
		"notification_methods",
		"created_at",
		"updated_at",
		"password_hash").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where("login = ?", login).
		ToSql()
	if err != nil {
		return model.User{}, err
	}

	var user repoModel.User
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.UUID,
		&user.Login,
		&user.Email,
		&user.NotificationMethods,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Password,
	)
	if err != nil {
		return model.User{}, err
	}

	return repoConverter.UserToModel(user), nil
}
