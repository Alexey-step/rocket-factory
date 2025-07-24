package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Alexey-step/rocket-factory/payment/internal/model"
	paymentV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	transactionUUID, err := a.service.PayOrder(ctx, req.GetOrderUuid(), req.GetOrderUuid(), req.GetPaymentMethod().String())
	if err != nil {
		if errors.Is(err, model.ErrPaymentInternalError) {
			return nil, status.Errorf(codes.Internal, "Payment service error: %v", err)
		}
		return nil, err
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
