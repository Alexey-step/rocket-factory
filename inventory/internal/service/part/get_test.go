package part

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"

	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
	"github.com/Alexey-step/rocket-factory/inventory/internal/repository/mocks"
)

func TestGetPartRepoSuccess(t *testing.T) {
	ctx := context.Background()
	uuid := gofakeit.UUID()
	part := getRepoMockPart(uuid)

	inventoryRepository := mocks.NewInventoryRepository(t)
	inventoryService := NewService(inventoryRepository)

	inventoryRepository.On("GetPart", ctx, uuid).Return(part, nil)

	res, err := inventoryService.GetPart(ctx, uuid)
	assert.NoError(t, err)
	assert.Equal(t, part, res)
}

func TestGetPartRepoError(t *testing.T) {
	ctx := context.Background()
	var (
		repoErr = gofakeit.Error()
		uuid    = gofakeit.UUID()
	)

	inventoryRepository := mocks.NewInventoryRepository(t)
	inventoryService := NewService(inventoryRepository)

	inventoryRepository.On("GetPart", ctx, uuid).Return(model.Part{}, repoErr)

	res, err := inventoryService.GetPart(ctx, uuid)
	assert.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Empty(t, res)
}

func getRepoMockPart(uuid string) model.Part {
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
		metadata = map[string]model.Metadata{
			"int64Value":  {Int64Value: lo.ToPtr(gofakeit.Int64())},
			"stringValue": {StringValue: lo.ToPtr(gofakeit.Word())},
			"doubleValue": {DoubleValue: lo.ToPtr(gofakeit.Float64())},
			"boolValue":   {BoolValue: lo.ToPtr(gofakeit.Bool())},
		}
		createdAt = time.Now()
	)

	tags := make([]string, gofakeit.Number(1, 5))
	for i := range tags {
		tags[i] = gofakeit.Word()
	}

	part := model.Part{
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

	return part
}
