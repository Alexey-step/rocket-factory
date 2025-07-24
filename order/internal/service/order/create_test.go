package order

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"

	clientMocks "github.com/Alexey-step/rocket-factory/order/internal/client/grpc/mocks"
	"github.com/Alexey-step/rocket-factory/order/internal/model"
	"github.com/Alexey-step/rocket-factory/order/internal/repository/mocks"
)

func TestCreateOrderSuccess(t *testing.T) {
	ctx := context.Background()
	userUUID := gofakeit.UUID()
	orderUUID := gofakeit.UUID()
	partUUIDs := []string{gofakeit.UUID()}
	price := gofakeit.Price(100, 1000)

	part := getMockedPart(orderUUID, price)

	info := model.OrderCreationInfo{
		OrderUUID:  orderUUID,
		TotalPrice: price,
	}

	filter := model.PartsFilter{
		Uuids: partUUIDs,
	}

	listParts := []model.Part{part}

	orderRepository := mocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)

	orderService := NewService(
		orderRepository,
		inventoryClient,
		paymentClient,
	)

	inventoryClient.On("ListParts", ctx, filter).Return(listParts, nil).Once()
	orderRepository.On("CreateOrder", ctx, userUUID, listParts).Return(info, nil).Once()
	resp, err := orderService.CreateOrder(ctx, userUUID, partUUIDs)

	assert.NoError(t, err)
	assert.Equal(t, info, resp)
}

func TestCreateOrderListPartsFail(t *testing.T) {
	ctx := context.Background()
	userUUID := gofakeit.UUID()
	partUUIDs := []string{gofakeit.UUID()}

	filter := model.PartsFilter{
		Uuids: partUUIDs,
	}

	expectedListPartsError := gofakeit.Error()

	orderRepository := mocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)

	orderService := NewService(
		orderRepository,
		inventoryClient,
		paymentClient,
	)

	inventoryClient.On("ListParts", ctx, filter).Return(nil, expectedListPartsError).Once()
	resp, err := orderService.CreateOrder(ctx, userUUID, partUUIDs)

	assert.Error(t, err)
	assert.Empty(t, resp)
	assert.Equal(t, err, expectedListPartsError)
}

func TestCreateOrderBadRequest(t *testing.T) {
	ctx := context.Background()
	userUUID := gofakeit.UUID()
	orderUUID := gofakeit.UUID()
	partUUIDs := []string{gofakeit.UUID()}
	price := gofakeit.Price(100, 1000)

	part := getMockedPart(orderUUID, price)

	part2 := getMockedPart(orderUUID, price)

	filter := model.PartsFilter{
		Uuids: partUUIDs,
	}

	listParts := []model.Part{part, part2}
	expectedErr := model.ErrOrderConflict

	orderRepository := mocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)

	orderService := NewService(
		orderRepository,
		inventoryClient,
		paymentClient,
	)

	inventoryClient.On("ListParts", ctx, filter).Return(listParts, nil).Once()
	resp, err := orderService.CreateOrder(ctx, userUUID, partUUIDs)

	assert.Error(t, err)
	assert.Empty(t, resp)
	assert.Equal(t, err, expectedErr)
}

func TestCreateOrderRepoErr(t *testing.T) {
	ctx := context.Background()
	userUUID := gofakeit.UUID()
	orderUUID := gofakeit.UUID()
	partUUIDs := []string{gofakeit.UUID()}
	price := gofakeit.Price(100, 1000)

	part := getMockedPart(orderUUID, price)

	filter := model.PartsFilter{
		Uuids: partUUIDs,
	}

	listParts := []model.Part{part}
	expectedErr := gofakeit.Error()

	orderRepository := mocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)

	orderService := NewService(
		orderRepository,
		inventoryClient,
		paymentClient,
	)

	inventoryClient.On("ListParts", ctx, filter).Return(listParts, nil).Once()
	orderRepository.On("CreateOrder", ctx, userUUID, listParts).Return(model.OrderCreationInfo{}, expectedErr).Once()
	resp, err := orderService.CreateOrder(ctx, userUUID, partUUIDs)

	assert.Error(t, err)
	assert.Empty(t, resp)
	assert.Equal(t, err, expectedErr)
}

func getMockedPart(uuid string, price float64) model.Part {
	var (
		name          = gofakeit.Name()
		description   = gofakeit.Paragraph(3, 5, 5, " ")
		stockQuantity = gofakeit.Int64()
		category      = gofakeit.RandomString([]string{"UNKNOWN", "ENGINE", "FUEL", "PORTHOLE", "WING"})
		dimensions    = model.Dimensions{
			Height: gofakeit.Float64Range(1.0, 10.0),
			Width:  gofakeit.Float64Range(1.0, 10.0),
			Length: gofakeit.Float64Range(1.0, 10.0),
			Weight: gofakeit.Float64Range(0.1, 5.0),
		}
		manufacturer = model.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		}
		metadata = model.Metadata{
			Int64Value: lo.ToPtr(gofakeit.Int64()),
		}
		createdAt = time.Now()
	)

	tags := make([]string, gofakeit.Number(1, 5))
	for i := range tags {
		tags[i] = gofakeit.Word()
	}

	return model.Part{
		UUID:          uuid,
		Name:          name,
		Description:   description,
		Price:         price,
		StockQuantity: stockQuantity,
		Category:      model.Category(category),
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          tags,
		Metadata:      metadata,
		CreatedAt:     createdAt,
	}
}
