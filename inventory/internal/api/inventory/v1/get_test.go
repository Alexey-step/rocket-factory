package v1

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Alexey-step/rocket-factory/inventory/internal/converter"
	"github.com/Alexey-step/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/Alexey-step/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (s *ServiceSuite) TestGetInventorySuccess() {
	uuid := gofakeit.UUID()

	part := model.Part{
		UUID:          uuid,
		Name:          gofakeit.Name(),
		Description:   "Primary propulsion unit",
		Price:         gofakeit.Float64Range(100, 10_000),
		StockQuantity: int64(gofakeit.Number(1, 100)),
		Category:      "ENGINE",
		Dimensions: model.Dimensions{
			Width:  gofakeit.Float64Range(0.1, 10.0),
			Height: gofakeit.Float64Range(0.1, 10.0),
			Length: gofakeit.Float64Range(0.1, 10.0),
			Weight: gofakeit.Float64Range(0.1, 10.0),
		},
		Manufacturer: model.Manufacturer{
			Name:    gofakeit.Name(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},
		Tags: []string{gofakeit.EmojiTag(), gofakeit.EmojiTag()},
		Metadata: model.Metadata{
			StringValue: lo.ToPtr(gofakeit.Word()),
			Int64Value:  lo.ToPtr(gofakeit.Int64()),
			DoubleValue: lo.ToPtr(gofakeit.Float64()),
			BoolValue:   lo.ToPtr(gofakeit.Bool()),
		},
		CreatedAt: timestamppb.Now().AsTime(),
	}

	s.service.On("GetPart", s.ctx, uuid).Return(part, nil)

	resp, err := s.api.GetPart(s.ctx, &inventoryV1.GetPartRequest{
		Uuid: uuid,
	})

	expectedPart := &inventoryV1.GetPartResponse{
		Part: converter.PartToProto(part),
	}

	s.NoError(err)
	s.Equal(resp, expectedPart)
}

func (s *ServiceSuite) TestGetInventoryFail() {
	uuid := gofakeit.UUID()

	expectedErr := model.ErrPartNotFound
	s.service.On("GetPart", s.ctx, uuid).Return(model.Part{}, expectedErr)

	resp, err := s.api.GetPart(s.ctx, &inventoryV1.GetPartRequest{
		Uuid: uuid,
	})

	s.Error(err)
	s.Equal(codes.NotFound, status.Code(err))
	s.Contains(err.Error(), "part with UUID")
	s.Empty(resp)
}
