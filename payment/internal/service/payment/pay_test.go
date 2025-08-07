package payment

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPayOrder(t *testing.T) {
	ctx := context.Background()
	var (
		orderUUID     = gofakeit.UUID()
		userUUID      = gofakeit.UUID()
		paymentMethod = gofakeit.RandomString([]string{"CARD", "SBP", "CREDIT_CARD"})
	)

	paymentService := NewService()

	transactionUUID, err := paymentService.PayOrder(ctx, orderUUID, userUUID, paymentMethod)
	assert.NoError(t, err)
	assert.NotEmpty(t, transactionUUID)
	parsed, err := uuid.Parse(transactionUUID)
	assert.NoError(t, err)
	assert.NotEmpty(t, parsed)
}
