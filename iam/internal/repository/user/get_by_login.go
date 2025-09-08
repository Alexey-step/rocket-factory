package user

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/crypto/bcrypt"

	"github.com/Alexey-step/rocket-factory/iam/internal/model"
	repoConverter "github.com/Alexey-step/rocket-factory/iam/internal/repository/converter"
	repoModel "github.com/Alexey-step/rocket-factory/iam/internal/repository/model"
)

func (r *repository) GetUserByLogin(ctx context.Context, login, password string) (model.User, error) {
	query, args, err := sq.Select(
		"uuid",
		"info",
		"created_at",
		"updated_at",
		"password_hash").
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where("info->>'login' = ?", login).
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
		&user.Password,
	)
	if err != nil {
		return model.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return model.User{}, err
	}

	return repoConverter.UserToModel(user), nil
}
