package converter

import (
	"time"

	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
	commonV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/common/v1"
)

func SessionToModel(in *commonV1.Session) model.Session {
	var updatedAt *time.Time
	if in.UpdatedAt != nil {
		t := in.UpdatedAt.AsTime()
		updatedAt = &t
	}

	return model.Session{
		UUID:      in.Uuid,
		CreatedAt: in.CreatedAt.AsTime(),
		UpdatedAt: updatedAt,
		ExpiresAt: in.ExpiresAt.AsTime(),
	}
}
