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

func (s *ServiceSuite) TestListPartsRepoError() {
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

	s.inventoryRepository.On("ListParts", s.ctx, filter).Return([]model.Part{}, repoErr)

	res, err := s.service.ListParts(s.ctx, filter)
	s.Error(err)
	s.ErrorIs(err, repoErr)
	s.Empty(res)
}
