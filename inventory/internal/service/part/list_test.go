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

func TestListPartsRepoSuccess(t *testing.T) {
	ctx := context.Background()
	filters := getMockFilters()
	part := getMockRepoPart()

	expectedParts := []model.Part{part}

	inventoryRepository := mocks.NewInventoryRepository(t)
	inventoryService := NewService(inventoryRepository)

	inventoryRepository.On("ListParts", ctx, filters).Return(expectedParts, nil)

	res, err := inventoryService.ListParts(ctx, filters)
	assert.NoError(t, err)
	assert.Equal(t, expectedParts, res)
}

func TestListPartsRepoError(t *testing.T) {
	ctx := context.Background()
	repoErr := gofakeit.Error()

	partsUUIDs := []string{gofakeit.UUID(), gofakeit.UUID()}
	partsNames := []string{gofakeit.Name(), gofakeit.Name()}
	partsCategories := []model.Category{"UNKNOWN", "ENGINE", "FUEL", "PORTHOLE", "WING"}
	manufactureCountries := []string{gofakeit.Country(), gofakeit.Country()}
	tags := []string{gofakeit.Word(), gofakeit.Word()}

	filter := model.PartsFilter{
		Uuids:                 partsUUIDs,
		Names:                 partsNames,
		Categories:            partsCategories,
		ManufacturerCountries: manufactureCountries,
		Tags:                  tags,
	}

	inventoryRepository := mocks.NewInventoryRepository(t)
	inventoryService := NewService(inventoryRepository)

	inventoryRepository.On("ListParts", ctx, filter).Return([]model.Part{}, repoErr)

	res, err := inventoryService.ListParts(ctx, filter)
	assert.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Empty(t, res)
}

func getMockFilters() model.PartsFilter {
	partsUUIDs := []string{gofakeit.UUID(), gofakeit.UUID()}
	partsNames := []string{gofakeit.Name(), gofakeit.Name()}
	partsCategories := []model.Category{"UNKNOWN", "ENGINE", "FUEL", "PORTHOLE", "WING"}
	manufactureCountries := []string{gofakeit.Country(), gofakeit.Country()}
	tags := []string{gofakeit.Word(), gofakeit.Word()}

	return model.PartsFilter{
		Uuids:                 partsUUIDs,
		Names:                 partsNames,
		Categories:            partsCategories,
		ManufacturerCountries: manufactureCountries,
		Tags:                  tags,
	}
}

func getMockRepoPart() model.Part {
	var (
		uuid          = gofakeit.UUID()
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

	tags := []string{gofakeit.Word(), gofakeit.Word()}

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
