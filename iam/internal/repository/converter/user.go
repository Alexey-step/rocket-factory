package converter

import (
	"github.com/Alexey-step/rocket-factory/iam/internal/model"
	repoModel "github.com/Alexey-step/rocket-factory/iam/internal/repository/model"
)

func UserToModel(user repoModel.User) model.User {
	return model.User{
		UUID:      user.UUID,
		Info:      userInfoToModel(user.Info),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func userInfoToModel(userInfo repoModel.UserInfo) model.UserInfo {
	notifyMethods := make([]model.NotificationMethod, len(userInfo.NotificationMethods))
	for i, m := range userInfo.NotificationMethods {
		notifyMethods[i] = notificationMethodToModel(m)
	}
	return model.UserInfo{
		Login:               userInfo.Login,
		Email:               userInfo.Email,
		NotificationMethods: notifyMethods,
	}
}

func notificationMethodToModel(notificationMethod repoModel.NotificationMethod) model.NotificationMethod {
	return model.NotificationMethod{
		ProviderName: notificationMethod.ProviderName,
		Target:       notificationMethod.Target,
	}
}

func UserToRepo(user model.User) repoModel.User {
	return repoModel.User{
		UUID:      user.UUID,
		Info:      UserInfoToRepo(user.Info),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserInfoToRepo(userInfo model.UserInfo) repoModel.UserInfo {
	notifyMethods := make([]repoModel.NotificationMethod, len(userInfo.NotificationMethods))
	for i, m := range userInfo.NotificationMethods {
		notifyMethods[i] = notificationMethodToRepo(m)
	}
	return repoModel.UserInfo{
		Login:               userInfo.Login,
		Email:               userInfo.Email,
		NotificationMethods: notifyMethods,
	}
}

func notificationMethodToRepo(notificationMethod model.NotificationMethod) repoModel.NotificationMethod {
	return repoModel.NotificationMethod{
		ProviderName: notificationMethod.ProviderName,
		Target:       notificationMethod.Target,
	}
}
