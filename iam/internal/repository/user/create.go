package user

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Alexey-step/rocket-factory/iam/internal/model"
	repoConverter "github.com/Alexey-step/rocket-factory/iam/internal/repository/converter"
	repoModel "github.com/Alexey-step/rocket-factory/iam/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, info model.UserInfo, password string) (userUUID string, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := repoModel.User{
		UUID:      uuid.NewString(),
		Info:      repoConverter.UserInfoToRepo(info),
		CreatedAt: time.Now(),
		Password:  hashedPassword,
	}

	query, args, err := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("uuid", "info", "created_at", "password_hash").
		Values(user.UUID, user.Info, user.CreatedAt, hashedPassword).
		Suffix("RETURNING uuid").
		ToSql()
	if err != nil {
		return "", err
	}

	var uuid string
	err = r.db.QueryRow(ctx, query, args...).Scan(&uuid)
	if err != nil {
		return "", err
	}

	return uuid, nil
}
