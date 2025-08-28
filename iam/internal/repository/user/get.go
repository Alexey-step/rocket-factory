package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/Alexey-step/rocket-factory/iam/internal/model"
	repoConverter "github.com/Alexey-step/rocket-factory/iam/internal/repository/converter"
	repoModel "github.com/Alexey-step/rocket-factory/iam/internal/repository/model"
)

func (r *repository) GetUser(ctx context.Context, userUUID string) (model.User, error) {
	query, args, err := sq.Select(
		"uuid",
		"info",
		"created_at",
		"updated_at").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"uuid": userUUID}).
		ToSql()
	if err != nil {
		return model.User{}, err
	}

	var user repoModel.User
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&user.UUID,
		&user.Info,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return model.User{}, err
	}

	return repoConverter.UserToModel(user), nil
}
