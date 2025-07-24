package part

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"

	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
)

func (s *ServiceSuite) TestGetPartRepoSuccess() {
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

	s.inventoryRepository.On("GetPart", s.ctx, uuid).Return(part, nil)

	res, err := s.service.GetPart(s.ctx, uuid)
	s.NoError(err)
	s.Equal(part, res)
}

func (s *ServiceSuite) TestGetPartRepoError() {
	var (
		repoErr = gofakeit.Error()
		uuid    = gofakeit.UUID()
	)

	s.inventoryRepository.On("GetPart", s.ctx, uuid).Return(model.Part{}, repoErr)

	res, err := s.service.GetPart(s.ctx, uuid)
	s.Error(err)
	s.ErrorIs(err, repoErr)
	s.Empty(res)
}
