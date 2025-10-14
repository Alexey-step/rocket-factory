package converter

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"github.com/Alexey-step/rocket-factory/iam/internal/model"
	repoModel "github.com/Alexey-step/rocket-factory/iam/internal/repository/model"
	"github.com/Alexey-step/rocket-factory/platform/pkg/logger"
)

func UserToModel(ctx context.Context, user repoModel.User) model.User {
	info := model.UserInfo{
		Login:               user.Login,
		Email:               user.Email,
		NotificationMethods: NotificationMethodsToModel(ctx, user.NotificationMethods),
		Password:            user.Password,
	}

	return model.User{
		UUID:      user.UUID,
		Info:      info,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func NotificationMethodsToModel(ctx context.Context, methods []byte) []model.NotificationMethod {
	var notificationMethods []model.NotificationMethod
	err := json.Unmarshal(methods, &notificationMethods)
	if err != nil {
		logger.Error(ctx, "failed to unmarshal notification methods", zap.Error(err))
	}
	return notificationMethods
}

func UserToRepo(ctx context.Context, user model.User) repoModel.User {
	return repoModel.User{
		UUID:                user.UUID,
		Login:               user.Info.Login,
		Email:               user.Info.Email,
		NotificationMethods: NotificationMethodsToRepo(ctx, user.Info.NotificationMethods),
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}
}

func NotificationMethodsToRepo(ctx context.Context, methods []model.NotificationMethod) []byte {
	notificationMethods, err := json.Marshal(methods)
	if err != nil {
		logger.Error(ctx, "failed to marshal notification methods", zap.Error(err))
	}
	return notificationMethods
}
