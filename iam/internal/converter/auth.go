package converter

import (
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Alexey-step/rocket-factory/iam/internal/model"
	commonV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/common/v1"
)

func SessionToProto(session model.Session) *commonV1.Session {
	return &commonV1.Session{
		Uuid:      session.UUID,
		CreatedAt: timestamppb.New(session.CreatedAt),
		UpdatedAt: timestamppb.New(lo.FromPtr(session.UpdatedAt)),
		ExpiresAt: timestamppb.New(session.ExpiresAt),
	}
}
