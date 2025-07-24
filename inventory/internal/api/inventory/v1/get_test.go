package v1

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Alexey-step/rocket-factory/inventory/internal/converter"
	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
	"github.com/Alexey-step/rocket-factory/inventory/internal/service/mocks"
	inventoryV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/inventory/v1"
)

func TestGetInventorySuccess(t *testing.T) {
	ctx := context.Background()
	uuid := "123e4567-e89b-12d3-a456-426614174000"

	part := getMockedModelPart(uuid)

	service := mocks.NewInventoryService(t)
	inventoryApi := NewAPI(service)

	service.On("GetPart", ctx, uuid).Return(part, nil)

	resp, err := inventoryApi.GetPart(ctx, &inventoryV1.GetPartRequest{
		Uuid: uuid,
	})

	expectedPart := &inventoryV1.GetPartResponse{
		Part: converter.PartToProto(part),
	}

	assert.NoError(t, err)
	assert.Equal(t, resp, expectedPart)
}

func TestGetInventoryFail(t *testing.T) {
	ctx := context.Background()
	uuid := gofakeit.UUID()

	expectedErr := model.ErrPartNotFound

	service := mocks.NewInventoryService(t)
	inventoryApi := NewAPI(service)

	service.On("GetPart", ctx, uuid).Return(model.Part{}, expectedErr)

	resp, err := inventoryApi.GetPart(ctx, &inventoryV1.GetPartRequest{
		Uuid: uuid,
	})

	assert.Error(t, err)
	assert.Equal(t, codes.NotFound, status.Code(err))
	assert.Empty(t, resp)
}

func getMockedModelPart(uuid string) model.Part {
	return model.Part{
		UUID:          uuid,
		Name:          "Falcon Engine",
		Description:   "Primary propulsion unit",
		Price:         5000.0,
		StockQuantity: 10,
		Category:      model.CategoryEngine,
		Dimensions: model.Dimensions{
			Width:  2.5,
			Height: 1.2,
			Length: 3.0,
			Weight: 150.0,
		},
		Manufacturer: model.Manufacturer{
			Name:    "SpaceX",
			Country: "USA",
			Website: "https://spacex.com",
		},
		Tags: []string{"rocket", "engine"},
		Metadata: map[string]model.Metadata{
			"stringValue": {StringValue: lo.ToPtr("serial-001")},
			"int64Value":  {Int64Value: lo.ToPtr(int64(42))},
			"doubleValue": {DoubleValue: lo.ToPtr(3.14)},
			"boolValue":   {BoolValue: lo.ToPtr(true)},
		},
		CreatedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
	}
}
