package grpc

import (
	"context"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
	authV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/auth/v1"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) (parts []model.Part, err error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, userUUID, orderUUID, paymentMethod string) (transactionUUID string, err error)
}

type IamClient interface {
	Login(ctx context.Context, req *authV1.LoginRequest) (*authV1.LoginResponse, error)
	Whoami(ctx context.Context, req *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error)
}
