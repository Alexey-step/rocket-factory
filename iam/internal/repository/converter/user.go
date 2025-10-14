package converter

import (
	"github.com/Alexey-step/rocket-factory/iam/internal/model"
	repoModel "github.com/Alexey-step/rocket-factory/iam/internal/repository/model"
)

func UserToModel(user repoModel.User) model.User {
	info := model.UserInfo{
		Login:               user.Login,
		Email:               user.Email,
		NotificationMethods: notificationMethodsToModel(user.NotificationMethods),
	}

	return model.User{
		UUID:      user.UUID,
		Info:      info,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func notificationMethodToModel(notificationMethod repoModel.NotificationMethod) model.NotificationMethod {
	return model.NotificationMethod{
		ProviderName: notificationMethod.ProviderName,
		Target:       notificationMethod.Target,
	}
}

func notificationMethodsToModel(methods []repoModel.NotificationMethod) []model.NotificationMethod {
	notifyMethods := make([]model.NotificationMethod, len(methods))
	for i, m := range methods {
		notifyMethods[i] = notificationMethodToModel(m)
	}
	return notifyMethods
}

func UserToRepo(user model.User) repoModel.User {
	notifyMethods := make([]repoModel.NotificationMethod, len(user.Info.NotificationMethods))
	for i, m := range user.Info.NotificationMethods {
		notifyMethods[i] = notificationMethodToRepo(m)
	}
	return repoModel.User{
		UUID:                user.UUID,
		Login:               user.Info.Login,
		Email:               user.Info.Email,
		NotificationMethods: notifyMethods,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}
}

func notificationMethodToRepo(notificationMethod model.NotificationMethod) repoModel.NotificationMethod {
	return repoModel.NotificationMethod{
		ProviderName: notificationMethod.ProviderName,
		Target:       notificationMethod.Target,
	}
}

func NotificationMethodsToRepo(methods []model.NotificationMethod) []repoModel.NotificationMethod {
	notifyMethods := make([]repoModel.NotificationMethod, len(methods))
	for i, m := range methods {
		notifyMethods[i] = notificationMethodToRepo(m)
	}
	return notifyMethods
}
