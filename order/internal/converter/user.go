package converter

import (
	"time"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	commonV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/common/v1"
)

func UserToModel(in *commonV1.User) model.User {
	var updatedAt *time.Time
	if in.UpdatedAt != nil {
		t := in.UpdatedAt.AsTime()
		updatedAt = &t
	}

	return model.User{
		UUID:      in.Uuid,
		Info:      userInfoToModel(in.Info),
		CreatedAt: in.CreatedAt.AsTime(),
		UpdatedAt: updatedAt,
	}
}

func userInfoToModel(in *commonV1.UserInfo) model.UserInfo {
	if in == nil {
		return model.UserInfo{}
	}
	return model.UserInfo{
		Login:               in.Login,
		Email:               in.Email,
		NotificationMethods: notificationMethodsToModel(in.NotificationMethods),
	}
}

func notificationMethodsToModel(in []*commonV1.NotificationMethod) []model.NotificationMethod {
	if in == nil {
		return nil
	}
	out := make([]model.NotificationMethod, len(in))
	for i, nm := range in {
		if nm != nil {
			out[i] = model.NotificationMethod{
				ProviderName: nm.ProviderName,
				Target:       nm.Target,
			}
		}
	}
	return out
}
