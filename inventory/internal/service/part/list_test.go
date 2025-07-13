package part

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"

	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
)

func (s *ServiceSuite) TestListPartsRepoSuccess() {
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
		metadata = model.Metadata{
			Int64Value: lo.ToPtr(gofakeit.Int64()),
		}
		createdAt = time.Now()
	)

	partsUUIDs := make([]string, gofakeit.Number(1, 5))
	for i := range partsUUIDs {
		partsUUIDs[i] = gofakeit.UUID()
	}

	partsNames := make([]string, gofakeit.Number(1, 5))
	for i := range partsNames {
		partsNames[i] = gofakeit.Name()
	}

	partsCategories := make([]model.Category, gofakeit.Number(1, 5))
	for i := range partsCategories {
		partsNames[i] = gofakeit.RandomString([]string{"UNKNOWN", "ENGINE", "FUEL", "PORTHOLE", "WING"})
	}

	manufactureCauntries := make([]string, gofakeit.Number(1, 5))
	for i := range partsCategories {
		partsNames[i] = gofakeit.Country()
	}

	tags := make([]string, gofakeit.Number(1, 5))
	for i := range tags {
		tags[i] = gofakeit.Word()
	}

	filter := model.PartsFilter{
		Uuids:                 partsUUIDs,
		Names:                 partsNames,
		Categories:            partsCategories,
		ManufacturerCountries: manufactureCauntries,
		Tags:                  tags,
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

	expectedParts := []model.Part{part}

	s.inventoryRepository.On("ListParts", s.ctx, filter).Return(expectedParts, nil)

	res, err := s.service.ListParts(s.ctx, filter)
	s.NoError(err)
	s.Equal(expectedParts, res)
}

func (s *ServiceSuite) TestGListPartsRepoError() {
	repoErr := gofakeit.Error()

	partsUUIDs := make([]string, gofakeit.Number(1, 5))
	for i := range partsUUIDs {
		partsUUIDs[i] = gofakeit.UUID()
	}

	partsNames := make([]string, gofakeit.Number(1, 5))
	for i := range partsNames {
		partsNames[i] = gofakeit.Name()
	}

	partsCategories := make([]model.Category, gofakeit.Number(1, 5))
	for i := range partsCategories {
		partsNames[i] = gofakeit.RandomString([]string{"UNKNOWN", "ENGINE", "FUEL", "PORTHOLE", "WING"})
	}

	manufactureCauntries := make([]string, gofakeit.Number(1, 5))
	for i := range partsCategories {
		partsNames[i] = gofakeit.Country()
	}

	tags := make([]string, gofakeit.Number(1, 5))
	for i := range tags {
		tags[i] = gofakeit.Word()
	}

	filter := model.PartsFilter{
		Uuids:                 partsUUIDs,
		Names:                 partsNames,
		Categories:            partsCategories,
		ManufacturerCountries: manufactureCauntries,
		Tags:                  tags,
	}

	s.inventoryRepository.On("ListParts", s.ctx, filter).Return([]model.Part{}, repoErr)

	res, err := s.service.ListParts(s.ctx, filter)
	s.Error(err)
	s.ErrorIs(err, repoErr)
	s.Empty(res)
}
