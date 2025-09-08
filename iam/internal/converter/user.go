package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Alexey-step/rocket-factory/iam/internal/model"
	commonV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/common/v1"
)

func UserToProto(user model.User) *commonV1.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt != nil {
		updatedAt = timestamppb.New(*user.UpdatedAt)
	}

	return &commonV1.User{
		Uuid: user.UUID,
		Info: &commonV1.UserInfo{
			Login:               user.Info.Login,
			Email:               user.Info.Email,
			NotificationMethods: notificationMethodToProto(user.Info.NotificationMethods),
		},
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func notificationMethodToProto(methods []model.NotificationMethod) []*commonV1.NotificationMethod {
	result := make([]*commonV1.NotificationMethod, len(methods))
	for _, method := range methods {
		result = append(result, &commonV1.NotificationMethod{
			ProviderName: method.ProviderName,
			Target:       method.Target,
		})
	}

	return result
}

func UserInfoToModel(user *commonV1.UserInfo) model.UserInfo {
	return model.UserInfo{
		Login:               user.GetLogin(),
		Email:               user.GetEmail(),
		NotificationMethods: notificationMethodToModel(user.GetNotificationMethods()),
	}
}

func notificationMethodToModel(methods []*commonV1.NotificationMethod) []model.NotificationMethod {
	result := make([]model.NotificationMethod, 0, len(methods))
	for _, method := range methods {
		result = append(result, model.NotificationMethod{
			ProviderName: method.GetProviderName(),
			Target:       method.GetTarget(),
		})
	}
	return result
}
