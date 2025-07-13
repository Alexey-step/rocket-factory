package order

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"

	"github.com/Alexey-step/rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestCreateOrderSuccess() {
	userUUID := gofakeit.UUID()
	orderUUID := gofakeit.UUID()
	partUUIDs := []string{gofakeit.UUID(), gofakeit.UUID()}

	var (
		name          = gofakeit.Name()
		description   = gofakeit.Paragraph(3, 5, 5, " ")
		price         = gofakeit.Price(100, 1000)
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

	part := model.Part{
		UUID:          orderUUID,
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

	info := model.OrderCreationInfo{
		OrderUUID:  orderUUID,
		TotalPrice: price,
	}

	filter := model.PartsFilter{
		Uuids: partUUIDs,
	}

	expectedParts := []model.Part{part}

	s.inventoryClient.On("ListParts", s.ctx, filter).Return(expectedParts, nil).Once()
	s.orderRepository.On("CreateOrder", s.ctx, userUUID, expectedParts).Return(info, nil).Once()
	resp, err := s.service.CreateOrder(s.ctx, userUUID, partUUIDs)

	s.NoError(err)
	s.Equal(info, resp)
}
